package parser

import (
	"errors"
	"fmt"
	"regexp"

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

// program : statement*
func (p *Parser) ParseProgram() (ast.Node, error) {
	result := ast.NewProgramNode(p.lex.Pos())

	for !p.over() {
		node, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		result.Statements = append(result.Statements, node)
	}

	return result, p.err()
}

// statement : mustache | block | rawBlock | partial | content | COMMENT
func (p *Parser) parseStatement() (ast.Node, error) {
	var result ast.Node

	tok := p.next()

	switch tok.Kind {
	case lexer.TokenContent:
		result = p.parseContent()
	case lexer.TokenComment:
		result = p.parseComment()
	default:
		return nil, errors.New(fmt.Sprintf("Failed to parse statement: %s", tok))
	}

	return result, p.err()
}

// content : CONTENT
func (p *Parser) parseContent() ast.Node {
	tok := p.shift()

	return ast.NewContentNode(tok.Pos, tok.Val)
}

// COMMENT
func (p *Parser) parseComment() ast.Node {
	tok := p.shift()

	value := rOpenComment.ReplaceAllString(tok.Val, "")
	value = rCloseComment.ReplaceAllString(value, "")

	return ast.NewCommentNode(tok.Pos, value)
}

// rawBlock : openRawBlock content END_RAW_BLOCK
func (p *Parser) parseRawBlock() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// openRawBlock : OPEN_RAW_BLOCK helperName param* hash? CLOSE_RAW_BLOCK
func (p *Parser) parseOpenRawBlock() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// block : openBlock program inverseChain? closeBlock
//       | openInverse program inverseAndProgram? closeBlock
func (p *Parser) parseBlock() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// openBlock : OPEN_BLOCK helperName param* hash? blockParams? CLOSE
func (p *Parser) parseOpenBlock() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// openInverse : OPEN_INVERSE helperName param* hash? blockParams? CLOSE
func (p *Parser) parseOpenInverse() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// openInverseChain : OPEN_INVERSE_CHAIN helperName param* hash? blockParams? CLOSE
func (p *Parser) parseOpenInverseChain() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// inverseAndProgram : INVERSE program
func (p *Parser) parseInverseAndProgram() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// inverseChain : openInverseChain program inverseChain?
//              | inverseAndProgram
func (p *Parser) parseInverseChain() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// closeBlock : OPEN_ENDBLOCK helperName CLOSE
func (p *Parser) parseCloseBlock() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// mustache : OPEN helperName param* hash? CLOSE
//          | OPEN_UNESCAPED helperName param* hash? CLOSE_UNESCAPED
func (p *Parser) parseMustache() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// partial : OPEN_PARTIAL partialName param* hash? CLOSE
func (p *Parser) parsePartial() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// param : helperName
//       | sexpr
func (p *Parser) parseParam() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// sexpr : OPEN_SEXPR helperName param* hash? CLOSE_SEXPR
func (p *Parser) parseSexpr() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// hash : hashSegment+
func (p *Parser) parseHash() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// hashSegment : ID EQUALS param
func (p *Parser) parseHashSegment() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// blockParams : OPEN_BLOCK_PARAMS ID+ CLOSE_BLOCK_PARAMS
func (p *Parser) parseBlockParams() (ast.Node, error) {

	return nil, nil
}

// helperName : path | dataName | STRING | NUMBER | BOOLEAN | UNDEFINED | NULL
func (p *Parser) parseHelperName() (ast.Node, error) {

	return nil, nil
}

// partialName : helperName | sexpr
func (p *Parser) parsePartialName() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// dataName : DATA pathSegments
func (p *Parser) parseDataName() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// x path : pathSegments
func (p *Parser) parsePath() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// pathSegments : pathSegments SEP ID
//              | ID
func (p *Parser) parsePathSegments() (ast.Node, error) {
	// @todo !!!
	return nil, nil
}

// Ensure there is at least a token to parse
func (p *Parser) ensure() {
	if len(p.tokens) == 0 {
		// fetch next token
		tok := p.lex.NextToken()

		// queue it
		p.tokens = append(p.tokens, &tok)
	}
}

// Returns next token without removing it from tokens buffer
func (p *Parser) next() *lexer.Token {
	p.ensure()

	return p.tokens[0]
}

// Returns next token and remove it from the tokens buffer
func (p *Parser) shift() *lexer.Token {
	var result *lexer.Token

	p.ensure()

	result, p.tokens = p.tokens[0], p.tokens[1:]

	return result
}

// Returns true if parsing is over
func (p *Parser) over() bool {
	tok := p.next()
	return (tok.Kind == lexer.TokenEOF) || (tok.Kind == lexer.TokenError)
}

// Returns lexer error, or nil if no error
func (p *Parser) err() error {
	if token := p.next(); token.Kind == lexer.TokenError {
		return errors.New(fmt.Sprintf("Lexer error: %s", token.String()))
	} else {
		return nil
	}
}
