package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

type Program struct {
	lineNo           int
	linePointerStack IntStack

	code      [][]Token
	variables map[string]Token

	intendation int
}

func (p *Program) calcLineIntend(lineNo int) int {
	if len(p.code[lineNo]) == 0 {
		return 0
	}

	firstToken := p.code[lineNo][0]
	if firstToken.tokenType == TokenWhiteSpace && (firstToken.byteValue() == '\n' || firstToken.byteValue() == '\r') {
		return 0
	}

	counter := 0
	for _, token := range p.code[lineNo] {
		if token.tokenType == TokenWhiteSpace && token.byteValue() == '\t' {
			counter++
		}
	}
	return counter
}

func (p *Program) IsRunning() bool {
	if p.lineNo < len(p.code) {
		return true
	}

	if p.linePointerStack.Len() > 0 {
		p.lineNo, _ = p.linePointerStack.Pop()
		p.intendation = p.calcLineIntend(p.lineNo)
		return true
	}

	return false
}

func (p *Program) computeExpression(tokens []Token) ([]Token, error) {
	tokens, err := convertExpressionToPostfix(tokens)

	if err != nil {
		log.Println("error: converting expression to postfix")
		return nil, err
	}

	return evaluateExpression(tokens, p.variables)
}

func (p *Program) RunLine() (err error) {
	tokens := p.code[p.lineNo]

	if len(tokens) == 0 {
		p.lineNo++
		return
	}

	lineIntend := p.calcLineIntend(p.lineNo)

	if lineIntend < p.intendation {
		p.lineNo, err = p.linePointerStack.Pop()

		if err != nil {
			log.Println("error: while going up pointer stack")
			return err
		}

		return
	}

	if lineIntend > p.intendation {
		return errors.New("error: invalid intendation found")
	}

	switch p.code[p.lineNo][lineIntend].tokenType {
	case TokenIdentifier:
		_, err = p.computeExpression(p.code[p.lineNo][lineIntend:])
		if err != nil {
			return err
		}
		p.lineNo++
		return
	case TokenKeyword:
		switch p.code[p.lineNo][lineIntend].stringValue() {
		case "while":
			condition := p.code[p.lineNo][lineIntend+1:]
			condition, err := p.computeExpression(condition)
			if err != nil {
				return err
			}

			if len(condition) > 1 {
				return fmt.Errorf("%d conditons found for while loop, only needs one", len(condition))
			}

			if len(condition) == 0 {
				return errors.New("no conditons found for while loop, needs one")
			}

			var isConditionTrue bool

			if condition[0].tokenType == TokenBoolean {
				isConditionTrue = condition[0].value.(bool)
			} else if condition[0].tokenType == TokenInteger {
				isConditionTrue = condition[0].intValue() != 0
			} else {
				return fmt.Errorf("%s is not a valid token for while condition", condition[0].ToString())
			}

			if isConditionTrue {
				p.linePointerStack.Push(p.lineNo)
				p.lineNo++
				p.intendation++
				return nil
			}

			p.lineNo++
			for p.lineNo < len(p.code) {
				if p.calcLineIntend(p.lineNo) == lineIntend {
					break
				}
				p.lineNo++
			}
			return nil
		}
	default:
		return errors.New("invalid start of line")
	}

	return errors.New("unexpected line ending found")
}

func tokenizeSourceCode(data []byte) (Program, error) {
	lines := bytes.Split(data, []byte("\n"))

	program := Program{
		lineNo:           0,
		linePointerStack: IntStack{},

		code:      make([][]Token, len(lines)),
		variables: make(map[string]Token),
	}

	for lineNo, line := range lines {
		tokens, err := TokenizeLine(line)

		if err != nil {
			log.Printf("error: tokonizing line %d", lineNo+1)
			return program, err
		}

		program.code[lineNo] = tokens
	}

	return program, nil
}
