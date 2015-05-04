package ast

import (
	"bytes"
	"fmt"
)

// References:
//   - https://github.com/wycats/handlebars.js/blob/master/lib/handlebars/compiler/ast.js
//   - https://github.com/golang/go/blob/master/src/text/template/parse/node.go

type NodeType int

type Node interface {
	Type() NodeType

	String() string

	// byte position of start of node in full original input string
	Position() Pos
}

type Pos int

func (p Pos) Position() Pos {
	return p
}

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeProgram NodeType = iota
	NodeContent
	NodeComment
)

// Changed to "%q" in tests for better error messages
var textFormat = "%s" // "%q"

//
// Program
//

type ProgramNode struct {
	NodeType
	Pos
	Statements []Node
}

func NewProgramNode(pos int) *ProgramNode {
	return &ProgramNode{
		NodeType: NodeProgram,
		Pos:      Pos(pos),
	}
}

func (node *ProgramNode) String() string {
	b := new(bytes.Buffer)

	for _, n := range node.Statements {
		fmt.Fprint(b, n)
	}

	return b.String()
}

//
// Content
//

type ContentNode struct {
	NodeType
	Pos
	Content string
}

func NewContentNode(pos int, text string) *ContentNode {
	return &ContentNode{
		NodeType: NodeContent,
		Pos:      Pos(pos),
		Content:  text,
	}
}

func (node *ContentNode) String() string {
	return fmt.Sprintf(textFormat, node.Content)
}

//
// Comment
//

type CommentNode struct {
	NodeType
	Pos
	Comment string
}

func NewCommentNode(pos int, text string) *CommentNode {
	return &CommentNode{
		NodeType: NodeComment,
		Pos:      Pos(pos),
		Comment:  text,
	}
}

func (node *CommentNode) String() string {
	return fmt.Sprintf(textFormat, node.Comment)
}
