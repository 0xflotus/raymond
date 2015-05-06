package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/aymerick/raymond/ast"
	"github.com/aymerick/raymond/lexer"
)

// References:
//   - https://github.com/wycats/handlebars.js/blob/master/src/handlebars.yy
//   - https://github.com/golang/go/blob/master/src/text/template/parse/parse.go

// Grammar parser
type Parser struct {
	// Lexer
	lex *lexer.Lexer

	// Root node
	root ast.Node

	// Tokens parsed but not consumed yet
	tokens []*lexer.Token

	// All tokens have been retreieved from lexer
	lexOver bool
}

var (
	rOpenComment  = regexp.MustCompile(`^\{\{~?!-?-?`)
	rCloseComment = regexp.MustCompile(`-?-?~?\}\}$`)
)

// instanciate a new parser
func New(input string) *Parser {
	return &Parser{
		lex: lexer.Scan(input),
	}
}

// parse given input and returns the ast root node
func Parse(input string) (ast.Node, error) {
	return New(input).ParseProgram()
}

func errExpect(msg string, expect lexer.TokenKind, tok *lexer.Token) error {
	return errors.New(fmt.Sprintf("%s Expected %s but got: %s", msg, expect, tok))
}

// program : statement*
func (p *Parser) ParseProgram() (*ast.Program, error) {
	result := ast.NewProgram(p.lex.Pos())

	for p.isStatement() {
		node, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		result.AddStatement(node)
	}

	return result, p.err()
}

// statement : mustache | block | rawBlock | partial | content | COMMENT
func (p *Parser) parseStatement() (ast.Node, error) {
	var result ast.Node
	var err error

	tok := p.next()

	switch tok.Kind {
	case lexer.TokenOpen, lexer.TokenOpenUnescaped:
		// mustache
		result, err = p.parseMustache()
	case lexer.TokenOpenBlock, lexer.TokenOpenInverse:
		// block
		result, err = p.parseBlock()
	case lexer.TokenOpenRawBlock:
		// rawBlock
		result, err = p.parseRawBlock()
	case lexer.TokenOpenPartial:
		// partial
		result, err = p.parsePartial()
	case lexer.TokenContent:
		// content
		result, err = p.parseContent()
	case lexer.TokenComment:
		// COMMENT
		result, err = p.parseComment()
	default:
		return nil, errors.New(fmt.Sprintf("Failed to parse statement: %s", tok))
	}

	if err != nil {
		return nil, err
	}

	return result, p.err()
}

// Returns true if next token starts a statement
func (p *Parser) isStatement() bool {
	if !p.have(1) {
		return false
	}

	switch p.next().Kind {
	case lexer.TokenOpen, lexer.TokenOpenUnescaped, lexer.TokenOpenBlock,
		lexer.TokenOpenInverse, lexer.TokenOpenRawBlock, lexer.TokenOpenPartial,
		lexer.TokenContent, lexer.TokenComment:
		return true
	}

	return false
}

// content : CONTENT
func (p *Parser) parseContent() (ast.Node, error) {
	tok := p.shift()
	if tok.Kind != lexer.TokenContent {
		return nil, errExpect("Failed to parse content.", lexer.TokenContent, tok)
	}

	return ast.NewContentStatement(tok.Pos, tok.Val), nil
}

// COMMENT
func (p *Parser) parseComment() (ast.Node, error) {
	tok := p.shift()
	if tok.Kind != lexer.TokenComment {
		return nil, errExpect("Failed to parse comment.", lexer.TokenComment, tok)
	}

	value := rOpenComment.ReplaceAllString(tok.Val, "")
	value = rCloseComment.ReplaceAllString(value, "")

	return ast.NewCommentStatement(tok.Pos, value), nil
}

// Parses `param* hash?`
func (p *Parser) parseExpressionParamsHash() (params []ast.Node, hash ast.Node, err error) {
	// params*
	if p.isParam() {
		params, err = p.parseParams()
		if err != nil {
			return
		}
	}

	// hash?
	if p.isHashSegment() {
		hash, err = p.parseHash()
	}

	return
}

// Parses an expression `helperName param* hash?`
func (p *Parser) parseExpression() (helperName ast.Node, params []ast.Node, hash ast.Node, err error) {
	// helperName
	helperName, err = p.parseHelperName()
	if err != nil {
		return
	}

	// param* hash?
	params, hash, err = p.parseExpressionParamsHash()

	return
}

// rawBlock : openRawBlock content endRawBlock
// openRawBlock : OPEN_RAW_BLOCK helperName param* hash? CLOSE_RAW_BLOCK
// endRawBlock : OPEN_END_RAW_BLOCK helperName CLOSE_RAW_BLOCK
func (p *Parser) parseRawBlock() (ast.Node, error) {
	var err error
	errMsg := "Failed to parse raw block."

	// OPEN_RAW_BLOCK
	tok := p.shift()

	result := ast.NewBlockStatement(tok.Pos)

	// helperName param* hash?
	result.Path, result.Params, result.Hash, err = p.parseExpression()
	if err != nil {
		return nil, err
	}

	openName, ok := result.Path.(*ast.PathExpression)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s Unexpected name in open block: %s", errMsg, result.Path))
	}

	// CLOSE_RAW_BLOCK
	tok = p.shift()
	if tok.Kind != lexer.TokenCloseRawBlock {
		return nil, errExpect(errMsg, lexer.TokenCloseRawBlock, tok)
	}

	// content
	content, err := p.parseContent()
	if err != nil {
		return nil, err
	}

	program := ast.NewProgram(tok.Pos)
	program.AddStatement(content)

	result.Program = program

	// OPEN_END_RAW_BLOCK
	tok = p.shift()
	if tok.Kind != lexer.TokenOpenEndRawBlock {
		return nil, errExpect(errMsg, lexer.TokenOpenEndRawBlock, tok)
	}

	// helperName
	endId, err := p.parseHelperName()
	if err != nil {
		return nil, err
	}

	closeName, ok := endId.(*ast.PathExpression)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s Unexpected name in end block: %s", errMsg, endId))
	}

	if openName.Original != closeName.Original {
		return nil, errors.New(fmt.Sprintf("%s Open and end blocks names mismatch: %s != %s", openName.Original, closeName.Original))
	}

	// CLOSE_RAW_BLOCK
	tok = p.shift()
	if tok.Kind != lexer.TokenCloseRawBlock {
		return nil, errExpect(errMsg, lexer.TokenCloseRawBlock, tok)
	}

	return result, nil
}

// openBlock : OPEN_BLOCK helperName param* hash? blockParams? CLOSE
// openInverse : OPEN_INVERSE helperName param* hash? blockParams? CLOSE
func (p *Parser) parseOpenBlock() (result *ast.BlockStatement, blockParams []string, inverse bool, err error) {
	errMsg := "Failed to parse open block."

	// OPEN_BLOCK | OPEN_INVERSE
	tokOpen := p.shift()

	result = ast.NewBlockStatement(tokOpen.Pos)

	// helperName param* hash?
	result.Path, result.Params, result.Hash, err = p.parseExpression()
	if err != nil {
		return
	}

	// blockParams?
	if p.isBlockParams() {
		blockParams, err = p.parseBlockParams()
		if err != nil {
			return
		}
	}

	// CLOSE
	tok := p.shift()
	if tok.Kind != lexer.TokenClose {
		err = errors.New(fmt.Sprintf("%s Unexpected name: %s", errMsg, result.Path))
		return
	}

	return
}

// block : openBlock program inverseChain? closeBlock
//       | openInverse program inverseAndProgram? closeBlock
//
// openInverseChain : OPEN_INVERSE_CHAIN helperName param* hash? blockParams? CLOSE
// inverseAndProgram : INVERSE program
// inverseChain : openInverseChain program inverseChain?
//              | inverseAndProgram
//
// closeBlock : OPEN_ENDBLOCK helperName CLOSE
func (p *Parser) parseBlock() (ast.Node, error) {
	var err error
	errMsg := "Failed to parse block."

	// openBlock | openInverse
	result, blockParams, inverse, err := p.parseOpenBlock()
	if err != nil {
		return nil, err
	}

	// program
	program, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}
	program.BlockParams = blockParams

	result.Program = program

	if inverse {
		// inverseAndProgram?
		if p.isInverse() {
			// @todo
		}
	} else {
		// inverseChain?
		if p.isInverseChain() {
			// @todo
		}
	}

	// OPEN_ENDBLOCK
	tok := p.shift()
	if tok.Kind != lexer.TokenOpenEndBlock {
		return nil, errExpect(errMsg, lexer.TokenOpenEndBlock, tok)
	}

	// helperName
	endId, err := p.parseHelperName()
	if err != nil {
		return nil, err
	}

	closeName, ok := endId.(*ast.PathExpression)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s Unexpected name in end block: %s", errMsg, endId))
	}

	openName, ok := result.Path.(*ast.PathExpression)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s Unexpected name: %s", errMsg, result.Path))
	}

	if openName.Original != closeName.Original {
		return nil, errors.New(fmt.Sprintf("%s Open and end blocks names mismatch: %s != %s", openName.Original, closeName.Original))
	}

	// CLOSE
	tok = p.shift()
	if tok.Kind != lexer.TokenCloseRawBlock {
		return nil, errExpect(errMsg, lexer.TokenCloseRawBlock, tok)
	}

	return result, p.err()
}

// mustache : OPEN helperName param* hash? CLOSE
//          | OPEN_UNESCAPED helperName param* hash? CLOSE_UNESCAPED
func (p *Parser) parseMustache() (ast.Node, error) {
	var err error

	// OPEN | OPEN_UNESCAPED
	tok := p.shift()

	closeToken := lexer.TokenClose
	if tok.Kind == lexer.TokenOpenUnescaped {
		closeToken = lexer.TokenCloseUnescaped
	}

	result := ast.NewMustacheStatement(tok.Pos)

	// helperName param* hash?
	result.Path, result.Params, result.Hash, err = p.parseExpression()
	if err != nil {
		return nil, err
	}

	// CLOSE | CLOSE_UNESCAPED
	tok = p.shift()
	if tok.Kind != closeToken {
		return nil, errExpect("Failed to parse Mustache Statement.", closeToken, tok)
	}

	return result, p.err()
}

// partial : OPEN_PARTIAL partialName param* hash? CLOSE
func (p *Parser) parsePartial() (ast.Node, error) {
	var err error

	// OPEN_PARTIAL
	tok := p.shift()

	result := ast.NewPartialStatement(tok.Pos)

	// partialName
	result.Name, err = p.parsePartialName()
	if err != nil {
		return nil, err
	}

	// param* hash?
	result.Params, result.Hash, err = p.parseExpressionParamsHash()
	if err != nil {
		return nil, err
	}

	// CLOSE
	tok = p.shift()
	if tok.Kind != lexer.TokenClose {
		return nil, errors.New(fmt.Sprintf("Failed to parse Partial Statement. Expected TokenClose, but got: %s", tok))
	}

	return result, p.err()
}

// Parses `helperName | sexpr`
func (p *Parser) parseHelperNameOrSexpr() (ast.Node, error) {
	if p.isSexpr() {
		// sexpr
		return p.parseSexpr()
	} else {
		// helperName
		return p.parseHelperName()
	}
}

// param : helperName | sexpr
func (p *Parser) parseParam() (ast.Node, error) {
	return p.parseHelperNameOrSexpr()
}

// Returns true if next tokens represent a `param`
func (p *Parser) isParam() bool {
	return (p.isSexpr() || p.isHelperName()) && !p.isHashSegment()
}

// parses `param*`
func (p *Parser) parseParams() ([]ast.Node, error) {
	var result []ast.Node

	for p.isParam() {
		param, err := p.parseParam()
		if err != nil {
			return nil, err
		}

		result = append(result, param)
	}

	return result, p.err()
}

// sexpr : OPEN_SEXPR helperName param* hash? CLOSE_SEXPR
func (p *Parser) parseSexpr() (ast.Node, error) {
	var err error
	errMsg := "Failed to parse SubExpression."

	// OPEN_SEXPR
	tok := p.shift()
	if tok.Kind != lexer.TokenOpenSexpr {
		return nil, errors.New(fmt.Sprintf("%s Expected TokenOpenSexpr: %s", errMsg, tok))
	}

	result := ast.NewSubExpression(tok.Pos)

	// helperName param* hash?
	result.Path, result.Params, result.Hash, err = p.parseExpression()
	if err != nil {
		return nil, err
	}

	// CLOSE_SEXPR
	tok = p.shift()
	if tok.Kind != lexer.TokenCloseSexpr {
		return nil, errors.New(fmt.Sprintf("%s Expected TokenCloseSexpr: %s", errMsg, tok))
	}

	return result, p.err()
}

// hash : hashSegment+
func (p *Parser) parseHash() (ast.Node, error) {
	var pairs []ast.Node

	for p.isHashSegment() {
		pair, err := p.parseHashSegment()
		if err != nil {
			return nil, err
		}

		pairs = append(pairs, pair)
	}

	if len(pairs) == 0 {
		return nil, errors.New(fmt.Sprintf("Failed to parse Hash: %s", p.next()))
	}

	result := ast.NewHash(int(pairs[0].Position()))
	result.Pairs = pairs

	return result, p.err()
}

// returns true if next tokens represents a `hashSegment`
func (p *Parser) isHashSegment() bool {
	return p.have(2) && (p.next().Kind == lexer.TokenID) && (p.nextAt(1).Kind == lexer.TokenEquals)
}

// hashSegment : ID EQUALS param
func (p *Parser) parseHashSegment() (ast.Node, error) {
	errMsg := "Failed to parse Hash Segment."

	// ID
	tokId := p.shift()
	if tokId.Kind != lexer.TokenID {
		return nil, errors.New(fmt.Sprintf("%s Expected an ID: %s", errMsg, tokId))
	}

	// EQUALS
	tokEquals := p.shift()
	if tokEquals.Kind != lexer.TokenEquals {
		return nil, errors.New(fmt.Sprintf("%s Expected an EQUAL: %s", errMsg, tokEquals))
	}

	// param
	param, err := p.parseParam()
	if err != nil {
		return nil, err
	}

	result := ast.NewHashPair(tokId.Pos)
	result.Key = tokId.Val
	result.Val = param

	return result, p.err()
}

// blockParams : OPEN_BLOCK_PARAMS ID+ CLOSE_BLOCK_PARAMS
func (p *Parser) parseBlockParams() ([]string, error) {
	var result []string
	errMsg := "Failed to parse Block Params."

	// OPEN_BLOCK_PARAMS
	tok := p.shift()
	if tok.Kind != lexer.TokenOpenBlockParams {
		return nil, errors.New(fmt.Sprintf("%s Expected a TokenOpenBlockParams: %s", errMsg, tok))
	}

	// ID+
	for p.isID() {
		result = append(result, p.shift().Val)
	}

	if len(result) == 0 {
		return nil, errors.New(fmt.Sprintf("%s Missing ID: %s", errMsg, tok))
	}

	// CLOSE_BLOCK_PARAMS
	tok = p.shift()
	if tok.Kind != lexer.TokenCloseBlockParams {
		return nil, errors.New(fmt.Sprintf("%s Expected a TokenCloseBlockParams: %s", errMsg, tok))
	}

	return result, p.err()
}

// helperName : path | dataName | STRING | NUMBER | BOOLEAN | UNDEFINED | NULL
func (p *Parser) parseHelperName() (ast.Node, error) {
	var result ast.Node
	var err error

	tok := p.next()

	switch tok.Kind {
	case lexer.TokenBoolean:
		// BOOLEAN
		p.shift()
		result = ast.NewBooleanLiteral(tok.Pos, (tok.Val == "true"), tok.Val)
	case lexer.TokenNumber:
		// NUMBER
		p.shift()
		val, err := strconv.Atoi(tok.Val)
		if err != nil {
			return nil, err
		}
		result = ast.NewNumberLiteral(tok.Pos, val, tok.Val)
	case lexer.TokenString:
		// STRING
		p.shift()
		result = ast.NewStringLiteral(tok.Pos, tok.Val)
	case lexer.TokenData:
		// dataName
		result, err = p.parseDataName()
		if err != nil {
			return nil, err
		}
	default:
		// path
		result, err = p.parsePath(false)
		if err != nil {
			return nil, err
		}
	}

	return result, p.err()
}

// Returns true if next tokens represent a `helperName`
func (p *Parser) isHelperName() bool {
	switch p.next().Kind {
	case lexer.TokenBoolean, lexer.TokenNumber, lexer.TokenString, lexer.TokenData, lexer.TokenID:
		return true
	}

	return false
}

// partialName : helperName | sexpr
func (p *Parser) parsePartialName() (ast.Node, error) {
	return p.parseHelperNameOrSexpr()
}

// dataName : DATA pathSegments
func (p *Parser) parseDataName() (ast.Node, error) {
	tok := p.shift()
	if tok.Kind != lexer.TokenData {
		return nil, errors.New(fmt.Sprintf("Failed to parse data: %s", tok))
	}

	return p.parsePath(true)
}

// path : pathSegments
// pathSegments : pathSegments SEP ID
//              | ID
func (p *Parser) parsePath(data bool) (ast.Node, error) {
	var tok *lexer.Token

	// ID
	tok = p.shift()
	if tok.Kind != lexer.TokenID {
		return nil, errors.New(fmt.Sprintf("Failed to parse path, expecting ID: %s", tok))
	}

	result := ast.NewPathExpression(tok.Pos, data)
	result.Part(tok.Val)

	for p.isPathSep() {
		// SEP
		tok = p.shift()
		result.Sep(tok.Val)

		// ID
		tok = p.shift()
		if tok.Kind != lexer.TokenID {
			return nil, errors.New(fmt.Sprintf("Failed to parse path, expecting ID after separator: %s", tok))
		}
		result.Part(tok.Val)
	}

	return result, p.err()
}

// Ensure there is token to parse at given index
func (p *Parser) ensure(index int) {
	if p.lexOver {
		// nothing more to grab
		return
	}

	nb := index + 1

	for len(p.tokens) < nb {
		// fetch next token
		tok := p.lex.NextToken()

		// queue it
		p.tokens = append(p.tokens, &tok)

		if (tok.Kind == lexer.TokenEOF) || (tok.Kind == lexer.TokenError) {
			p.lexOver = true
			break
		}
	}
}

// Returns true is there are a list given number of tokens to consume left
func (p *Parser) have(nb int) bool {
	p.ensure(nb - 1)

	return len(p.tokens) >= nb
}

// Returns next token at given index, without consuming it
func (p *Parser) nextAt(index int) *lexer.Token {
	p.ensure(index)

	return p.tokens[index]
}

// Returns next token without consuming it
func (p *Parser) next() *lexer.Token {
	return p.nextAt(0)
}

// Returns next token and remove it from the tokens buffer
func (p *Parser) shift() *lexer.Token {
	var result *lexer.Token

	p.ensure(0)

	result, p.tokens = p.tokens[0], p.tokens[1:]

	return result
}

// Returns lexer error, or nil if no error
func (p *Parser) err() error {
	if token := p.next(); token.Kind == lexer.TokenError {
		return errors.New(fmt.Sprintf("Lexer error: %s", token.String()))
	} else {
		return nil
	}
}

// Returns true if next token is of given type
func (p *Parser) isToken(kind lexer.TokenKind) bool {
	return p.have(1) && p.next().Kind == kind
}

// returns true if next token starts a sexpr
func (p *Parser) isSexpr() bool {
	return p.isToken(lexer.TokenOpenSexpr)
}

// Returns true if next token is a path separator
func (p *Parser) isPathSep() bool {
	return p.isToken(lexer.TokenSep)
}

// Returns true if next token is an ID
func (p *Parser) isID() bool {
	return p.isToken(lexer.TokenID)
}

// Returns true if next token starts a block params
func (p *Parser) isBlockParams() bool {
	return p.isToken(lexer.TokenOpenBlockParams)
}

// Returns true if next token starts an INVERSE sequence
func (p *Parser) isInverse() bool {
	return p.isToken(lexer.TokenOpenInverse)
}

// Returns true if next token starts an INVERSE_CHAIN v
func (p *Parser) isInverseChain() bool {
	return p.isToken(lexer.TokenOpenInverseChain)
}
