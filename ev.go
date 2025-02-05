package main

import (
	"errors"
	"fmt"
	"log"
)

func getOperatorPrecidence(operator Token) byte {
	if operator.tokenType == TokenAssignment {
		return 1
	}

	if operator.tokenType == TokenConditionalOperator {
		return 2
	}

	switch operator.value.(byte) {
	case '+', '-':
		return 3
	case '*', '/', '%':
		return 4
	case '^':
		return 5
	}

	return 0
}

func convertExpressionToPostfix(infix []Token) ([]Token, error) {
	stack := TokenStack{}
	result := TokenStack{}

	stack.Push(Token{TokenParentheses, byte('(')})
	infix = append(infix, Token{TokenParentheses, byte(')')})

	for _, token := range infix {
		switch token.tokenType {
		case TokenInteger, TokenIdentifier:
			result.Push(token)
			continue
		case TokenParentheses:
			char := token.byteValue()

			if char == '(' {
				stack.Push(token)
				continue
			}

			if char == ')' {
				topElem, err := stack.Peek()
				if err != nil {
					log.Println("error: while encountering a closing parenthesis")
					return nil, err
				}

				for topElem.tokenType != TokenParentheses {
					if _, err = stack.Pop(); err != nil { // will never happen
						log.Panicln(err)
					}
					result.Push(topElem)

					topElem, err = stack.Peek()
					if err != nil {
						log.Println("error: while dealing with a closing parenthesis")
						return nil, err
					}
				}
				stack.Pop()
				continue
			}
		case TokenUnaryOperator:
			stack.Push(token)
		case TokenBinaryOperator, TokenAssignment, TokenConditionalOperator:
			topElem, err := stack.Peek()

			if err != nil {
				log.Println("error: while encountering new operator")
				return nil, err
			}

			if topElem.tokenType == TokenParentheses {
				stack.Push(token)
				continue
			}

			if getOperatorPrecidence(topElem) >= getOperatorPrecidence(token) {
				if _, err = stack.Pop(); err != nil { // will never happen
					log.Panicln(err)
				}
				result.Push(topElem)
			}

			stack.Push(token)
		case TokenWhiteSpace:
			continue
		default:
			return nil, fmt.Errorf("unknown token %s found while parsing infix", token.ToString())
		}
	}

	if stack.Len() > 0 {
		return nil, errors.New("too many operators/too few operands found while parsing infix")
	}

	return result.ToArr(), nil
}

func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
