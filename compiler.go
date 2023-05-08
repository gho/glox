package main

import (
	"fmt"
	"strconv"
)

type Compiler interface {
	compile(source string) (*Chunk, error)
}

type precedence byte

const (
	precNone       precedence = iota
	precAssignment            // =
	precOr                    // or
	precAnd                   // and
	precEquality              // == !=
	precComparison            // < <= >= >
	precTerm                  // + -
	precFactor                // * /
	precUnary                 // - !
	precCall                  // . ()
	precPrimary
)

type parseFn func(chunk *Chunk) error

type parseRule struct {
	prefix     parseFn
	infix      parseFn
	precedence precedence
}

type compiler struct {
	scanner    Scanner
	parseRules map[TokenType]parseRule
	current    Token
	previous   Token
}

func newCompiler() Compiler {
	c := &compiler{}
	c.parseRules = map[TokenType]parseRule{
		TokenEOF:        {nil, nil, precNone},
		TokenNil:        {c.literal, nil, precNone},
		TokenFalse:      {c.literal, nil, precNone},
		TokenTrue:       {c.literal, nil, precNone},
		TokenLeftParen:  {c.grouping, nil, precNone},
		TokenRightParen: {nil, nil, precNone},
		TokenPlus:       {c.unary, c.binary, precTerm},
		TokenMinus:      {c.unary, c.binary, precTerm},
		TokenStar:       {c.unary, c.binary, precFactor},
		TokenSlash:      {c.unary, c.binary, precFactor},
		TokenEqualEqual: {nil, c.binary, precEquality},
		TokenGreater:    {nil, c.binary, precComparison},
		TokenLess:       {nil, c.binary, precComparison},
		TokenBang:       {c.unary, nil, precNone},
		TokenNumber:     {c.number, nil, precNone},
	}
	return c
}

func (c *compiler) advance() {
	c.previous = c.current
	c.current = c.scanner.nextToken()
}

func (c *compiler) consume(typ TokenType) error {
	if c.current.typ != typ {
		return fmt.Errorf("expected %v, got %v", typ, c.current.typ)
	}
	c.advance()
	return nil
}

func (c *compiler) compile(source string) (*Chunk, error) {
	chunk := &Chunk{}
	c.scanner = newScanner(source)

	c.advance()

loop:
	for {
		t := c.current
		switch t.typ {
		case TokenError:
			return nil, fmt.Errorf("%d: %s", t.line, t.data)
		case TokenEOF:
			break loop
		default:
			if err := c.expression(chunk); err != nil {
				return nil, err
			}
		}
	}

	chunk.addOp(OpReturn)

	return chunk, nil
}

func (c *compiler) parse(chunk *Chunk, prec precedence) error {
	c.advance()

	rule, err := c.getParseRule(c.previous.typ)
	if err != nil {
		return err
	}

	prefix := rule.prefix
	if prefix == nil {
		return fmt.Errorf("syntax error")
	}

	if err = prefix(chunk); err != nil {
		return err
	}

	for {
		rule, err = c.getParseRule(c.current.typ)
		if err != nil {
			return err
		}

		if prec > rule.precedence {
			break
		}

		c.advance()
		infix := rule.infix
		if err = infix(chunk); err != nil {
			return err
		}
	}

	return nil
}

func (c *compiler) getParseRule(typ TokenType) (*parseRule, error) {
	rule, ok := c.parseRules[typ]
	if !ok {
		return nil, fmt.Errorf("unknown token type: %s", typ)
	}
	return &rule, nil
}

func (c *compiler) expression(chunk *Chunk) error {
	return c.parse(chunk, precAssignment)
}

var literalOps = map[TokenType]Op{
	TokenNil:   OpNil,
	TokenFalse: OpFalse,
	TokenTrue:  OpTrue,
}

func (c *compiler) literal(chunk *Chunk) error {
	typ := c.previous.typ

	op, ok := literalOps[typ]
	if !ok {
		return fmt.Errorf("unknown literal token: %s", typ)
	}
	chunk.addOp(op)
	return nil
}

func (c *compiler) number(chunk *Chunk) error {
	f, err := strconv.ParseFloat(c.previous.data, 64)
	if err != nil {
		return err
	}

	val := numberValue(f)

	index := chunk.addVal(val)
	if index > 255 {
		return fmt.Errorf("too many constants")
	}

	chunk.addOp(OpConstant)
	chunk.addByte(byte(index))

	return nil
}

func (c *compiler) grouping(chunk *Chunk) error {
	if err := c.expression(chunk); err != nil {
		return err
	}
	return c.consume(TokenRightParen)
}

var unaryOps = map[TokenType]Op{
	TokenMinus: OpNegate,
	TokenBang:  OpNot,
}

func (c *compiler) unary(chunk *Chunk) error {
	typ := c.previous.typ

	if err := c.parse(chunk, precUnary); err != nil {
		return err
	}

	op, ok := unaryOps[typ]
	if !ok {
		return fmt.Errorf("unknown unary op: %s", typ)
	}
	chunk.addOp(op)

	return nil
}

var binaryOps = map[TokenType]Op{
	TokenPlus:       OpAdd,
	TokenMinus:      OpSubtract,
	TokenStar:       OpMultiply,
	TokenSlash:      OpDivide,
	TokenEqualEqual: OpEqual,
	TokenGreater:    OpGreater,
	TokenLess:       OpLess,
}

func (c *compiler) binary(chunk *Chunk) error {
	typ := c.previous.typ

	rule, err := c.getParseRule(typ)
	if err != nil {
		return err
	}

	if err := c.parse(chunk, rule.precedence+1); err != nil {
		return err
	}

	op, ok := binaryOps[typ]
	if !ok {
		return fmt.Errorf("unknown binary op: %s", typ)
	}
	chunk.addOp(op)

	return nil
}
