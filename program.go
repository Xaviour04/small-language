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

func (p *Program) GetValue(identifier Token) (Token, error) {
	switch identifier.tokenType {
	case TokenIdentifier:
		value, ok := p.variables[identifier.stringValue()]

		if !ok {
			return Token{}, fmt.Errorf("variable %q cannot be read before initializing", identifier.stringValue())
		}

		return value, nil
	case TokenInteger, TokenBoolean:
		return identifier, nil
	default:
		return Token{}, fmt.Errorf("cannot get value of %s", identifier.ToString())
	}

}

func (p *Program) computeExpression(tokens []Token) ([]Token, error) {
	tokens, err := convertExpressionToPostfix(tokens)

	if err != nil {
		log.Println("error: converting expression to postfix")
		return nil, err
	}

	return p.evaluateExpression(tokens)
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
		p.intendation = p.calcLineIntend(p.lineNo)

		if err != nil {
			log.Println("error: while going up pointer stack")
			return err
		}

		return
	}

	if lineIntend > p.intendation {
		return errors.New("error: invalid intendation found")
	}

	switch tokens[lineIntend].tokenType {
	case TokenIdentifier:
		_, err = p.computeExpression(tokens[lineIntend:])
		if err != nil {
			return err
		}
		p.lineNo++
		return
	case TokenKeyword:
		switch tokens[lineIntend].stringValue() {
		case "print":
			tokens, err = p.evaluateExpression(tokens[lineIntend+1:])
			if err != nil {
				log.Println("error evaluating print expression")
				return err
			}
			str := ""
			for _, token := range tokens {
				switch token.tokenType {
				case TokenInteger:
					str += fmt.Sprintf("%d ", token.intValue())
				case TokenBoolean:
					if token.value.(bool) {
						str += "True "
					} else {
						str += "False "
					}
				}
			}
			fmt.Println(str)
			p.lineNo++
			return

		case "while":
			condition := tokens[lineIntend+1:]
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

func (p *Program) evaluateExpression(tokens []Token) ([]Token, error) {
	stack := TokenStack{}

	for _, token := range tokens {
		switch token.tokenType {
		case TokenInteger, TokenBoolean, TokenIdentifier:
			stack.Push(token)

		case TokenAssignment:
			value, err := stack.Pop()
			if err != nil {
				return nil, errors.New("no value found to assign")
			}

			name, err := stack.Pop()
			if err != nil {
				return nil, errors.New("no identifier found to assign")
			}
			p.variables[name.stringValue()], err = p.GetValue(value)
			if err != nil {
				return nil, err
			}

		case TokenUnaryOperator:
			value, err := stack.Pop()
			if err != nil {
				return nil, errors.New("no value found for unirary operator")
			}
			value, err = p.GetValue(value)
			if err != nil {
				return nil, err
			}
			switch token.byteValue() {
			case '-':
				stack.Push(Token{TokenInteger, -1 * value.intValue()})
			}

		case TokenBinaryOperator, TokenConditionalOperator:
			rhs, err := stack.Pop()
			if err != nil {
				log.Printf("error: no values available for binary operator %c", token.byteValue())
				return nil, err
			}
			lhs, err := stack.Pop()
			if err != nil {
				log.Printf("error: only one value found for binary operator %c", token.byteValue())
				return nil, err
			}

			lhs, err = p.GetValue(lhs)
			if err != nil {
				return nil, err
			}

			rhs, err = p.GetValue(rhs)
			if err != nil {
				return nil, err
			}

			if rhs.tokenType != TokenInteger {
				return nil, fmt.Errorf("invalid Token{type: %d, value: %v} found while evaluating the rhs for binary operator %c", rhs.tokenType, rhs.value, token.byteValue())
			}

			if lhs.tokenType != TokenInteger {
				return nil, fmt.Errorf("invalid Token{type: %d, value: %v} found while evaluating the lhs for binary operator %c", rhs.tokenType, rhs.value, token.byteValue())
			}

			rhsValue := rhs.intValue()
			lhsValue := lhs.intValue()

			switch token.byteValue() {
			case '+':
				stack.Push(Token{TokenInteger, lhsValue + rhsValue})
			case '-':
				stack.Push(Token{TokenInteger, lhsValue - rhsValue})
			case '*':
				stack.Push(Token{TokenInteger, lhsValue * rhsValue})
			case '/':
				stack.Push(Token{TokenInteger, lhsValue / rhsValue})
			case '%':
				stack.Push(Token{TokenInteger, lhsValue % rhsValue})
			case '^':
				stack.Push(Token{TokenInteger, pow(lhsValue, rhsValue)})
			case '<':
				stack.Push(Token{TokenBoolean, lhsValue < rhsValue})
			case '>':
				stack.Push(Token{TokenBoolean, lhsValue > rhsValue})
			case 'l':
				stack.Push(Token{TokenBoolean, lhsValue <= rhsValue})
			case 'g':
				stack.Push(Token{TokenBoolean, lhsValue >= rhsValue})
			case 'e':
				stack.Push(Token{TokenBoolean, lhsValue == rhsValue})
			default:
				return nil, fmt.Errorf("invalid operator %c", token.byteValue())
			}
		}
	}

	result := make([]Token, stack.Len())

	for i, token := range stack.ToArr() {
		var err error
		if token.tokenType == TokenIdentifier {
			result[i], err = p.GetValue(token)
			if err != nil {
				return nil, err
			}
			continue
		}

		if token.tokenType == TokenInteger || token.tokenType == TokenBoolean {
			result[i] = token
			continue
		}

		return nil, fmt.Errorf("invalid expression evaluation found: %s", token.ToString())
	}

	return result, nil
}
