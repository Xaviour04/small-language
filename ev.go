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

func evaluateExpression(tokens []Token, variables map[string]Token) ([]Token, error) {
	stack := TokenStack{}

	for _, token := range tokens {
		switch token.tokenType {
		case TokenInteger, TokenIdentifier:
			stack.Push(token)

		case TokenAssignment:
			value, err := stack.Pop()
			if err != nil {
				return nil, errors.New("no value found to assign")
			}

			var ok bool

			if value.tokenType == TokenIdentifier {
				value, ok = variables[value.stringValue()]

				if !ok {
					return nil, fmt.Errorf("trying to read from variable %s before initializing", value.stringValue())
				}
			}

			name, err := stack.Pop()
			if err != nil {
				return nil, errors.New("no identifier found to assign")
			}
			variables[name.stringValue()] = value

		case TokenUnaryOperator:
			value, err := stack.Pop()
			if err != nil {
				return nil, errors.New("no value found for unirary operator")
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

			var rhsValue, lhsValue int

			if rhs.tokenType == TokenIdentifier {
				rhsToken, ok := variables[rhs.stringValue()]

				if !ok {
					return nil, fmt.Errorf("unintialized variable %q", rhs.stringValue())
				}

				if rhsToken.tokenType != TokenInteger {
					return nil, fmt.Errorf("invalid token %s cannot be used for operations", rhsToken.ToString())
				}

				rhsValue = rhsToken.intValue()
			} else if rhs.tokenType == TokenInteger {
				rhsValue = rhs.intValue()
			} else {
				return nil, fmt.Errorf("invalid Token{type: %d, value: %v} found while evaluating the rhs for binary operator %c", rhs.tokenType, rhs.value, token.byteValue())
			}

			if lhs.tokenType == TokenIdentifier {
				lhsToken, ok := variables[lhs.stringValue()]

				if !ok {
					return nil, fmt.Errorf("unintialized variable %q", lhs.stringValue())
				}

				if lhsToken.tokenType != TokenInteger {
					return nil, fmt.Errorf("invalid token %s cannot be used for operations", lhs.ToString())
				}

				lhsValue = lhsToken.intValue()
			} else if lhs.tokenType == TokenInteger {
				lhsValue = lhs.intValue()
			} else {
				return nil, fmt.Errorf("invalid Token{type: %d, value: %v} found while evaluating the lhs for binary operator %c", rhs.tokenType, rhs.value, token.byteValue())
			}

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

	return stack.ToArr(), nil
}
