package parser

import (
	"fmt"

	"ksm/lexer"
)

type NodeType string

const (
	NodeVarDecl   NodeType = "VAR_DECL"
	NodePrint     NodeType = "PRINT"
	NodeIf        NodeType = "IF"
	NodeBlock     NodeType = "BLOCK"
	NodeOtherwise NodeType = "OTHERWISE"
)

type Node struct {
	Type     NodeType
	Literal  string
	Children []*Node
}

type Parser struct {
	lexer   *lexer.Lexer
	current lexer.Token
}

// NewParser creates and returns a new instance of Parser.
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}
	p.nextToken() // Initialize the first token
	return p
}

// Parse is the method that starts parsing input and returns the root node of the AST.
func (p *Parser) Parse() (*Node, error) {
	// Start parsing with a block node to handle multiple statements
	block := &Node{Type: NodeBlock}

	// Loop until reaching the end of the file (or input)
	for p.current.Type != lexer.TokenEOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, fmt.Errorf("error parsing statement: %v", err)
		}
		if stmt != nil {
			block.Children = append(block.Children, stmt)
		}
	}

	return block, nil
}

// Add other parsing methods here...

func (p *Parser) parseStatement() (*Node, error) {
	switch p.current.Type {
	case lexer.TokenKeyword:
		switch p.current.Literal {
		case "declare":
			return p.parseVarDecl()
		case "displayln":
			return p.parsePrint()
		case "if":
			return p.parseIf()
		case "otherwise":
			return p.parseOtherwise()
		default:
			return nil, fmt.Errorf("unexpected keyword: %s", p.current.Literal)
		}
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.current)
	}
}

func (p *Parser) parseVarDecl() (*Node, error) {
	p.nextToken() // move past 'declare'
	if p.current.Type != lexer.TokenIdent {
		return nil, fmt.Errorf("expected identifier, got %v", p.current)
	}
	varName := p.current.Literal
	p.nextToken() // move past variable name
	if p.current.Literal != "=" {
		return nil, fmt.Errorf("expected '=', got %v", p.current)
	}
	p.nextToken() // move past '='
	if p.current.Type != lexer.TokenNumber && p.current.Type != lexer.TokenIdent && p.current.Type != lexer.TokenString {
		return nil, fmt.Errorf("expected number, identifier, or string, got %v", p.current)
	}
	varValue := p.current.Literal
	p.nextToken() // move past value

	return &Node{
		Type:    NodeVarDecl,
		Literal: varName + " = " + varValue,
	}, nil
}

func (p *Parser) parsePrint() (*Node, error) {
	p.nextToken() // move past 'displayln'
	if p.current.Literal != "(" {
		return nil, fmt.Errorf("expected '(', got %v", p.current)
	}
	p.nextToken() // move past '('
	if p.current.Type != lexer.TokenIdent && p.current.Type != lexer.TokenNumber && p.current.Type != lexer.TokenString {
		return nil, fmt.Errorf("expected identifier, number, or string, got %v", p.current)
	}
	printValue := p.current.Literal
	p.nextToken() // move past value
	if p.current.Literal != ")" {
		return nil, fmt.Errorf("expected ')', got %v", p.current)
	}
	p.nextToken() // move past ')'

	return &Node{
		Type:    NodePrint,
		Literal: "print " + printValue,
	}, nil
}

func (p *Parser) parseIf() (*Node, error) {
	p.nextToken() // move past 'if'
	if p.current.Literal != "case" {
		return nil, fmt.Errorf("expected 'case', got %v", p.current)
	}
	p.nextToken() // move past 'case'

	condition, err := p.parseCondition()
	if err != nil {
		return nil, err
	}

	if p.current.Literal != "{" {
		return nil, fmt.Errorf("expected '{', got %v", p.current)
	}
	p.nextToken() // move past '{'

	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	if p.current.Literal != "}" {
		return nil, fmt.Errorf("expected '}', got %v", p.current)
	}
	p.nextToken() // move past '}'

	return &Node{
		Type:     NodeIf,
		Literal:  "if " + condition,
		Children: []*Node{block},
	}, nil
}

func (p *Parser) parseCondition() (string, error) {
	left := p.current.Literal
	p.nextToken() // move past left operand

	if p.current.Type != lexer.TokenOperator {
		return "", fmt.Errorf("expected operator, got %v", p.current)
	}
	operator := p.current.Literal
	p.nextToken() // move past operator

	right := p.current.Literal
	p.nextToken() // move past right operand

	return fmt.Sprintf("%s %s %s", left, operator, right), nil
}

func (p *Parser) parseOtherwise() (*Node, error) {
	p.nextToken() // move past 'otherwise'
	if p.current.Literal != "{" {
		return nil, fmt.Errorf("expected '{', got %v", p.current)
	}
	p.nextToken() // move past '{'
	block, err := p.parseBlock()
	if err != nil {
		return nil, err
	}
	if p.current.Literal != "}" {
		return nil, fmt.Errorf("expected '}', got %v", p.current)
	}
	p.nextToken() // move past '}'

	return &Node{
		Type:     NodeOtherwise,
		Literal:  "otherwise",
		Children: []*Node{block},
	}, nil
}

func (p *Parser) parseBlock() (*Node, error) {
	block := &Node{Type: NodeBlock}
	for p.current.Type != lexer.TokenEOF && p.current.Literal != "}" {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Children = append(block.Children, stmt)
		}
	}
	return block, nil
}

func (p *Parser) nextToken() {
	p.current = p.lexer.NextToken()
}
