package main

import (
	"errors"
	"fmt"
)

func getOperatorPrecidence(operator byte) byte {
	switch operator {
	case '=':
		return 1
	case '+', '-':
		return 2
	case '*', '/', '%':
		return 3
	case '^':
		return 4
	default:
		return 0
	}
}

func convertExpressionToPostfix(infix []Token) ([]Token, error) {
	stack := Stack{}
	result := Stack{}

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
				for stack.Peek().tokenType != TokenParentheses {
					result.Push(stack.Pop())
				}
				stack.Pop()
				continue
			}
		case TokenUnaryOperator:
			stack.Push(token)
		case TokenBinaryOperator, TokenAssignment:
			if stack.Peek().tokenType == TokenParentheses {
				stack.Push(token)
				continue
			}

			if getOperatorPrecidence(stack.Peek().byteValue()) >= getOperatorPrecidence(token.byteValue()) {
				result.Push(stack.Pop())
			}

			stack.Push(token)
		case TokenWhiteSpace:
			continue
		default:
			return nil, fmt.Errorf("unknown token (type: %d, value: \"%v\") found while parsing infix", token.tokenType, token.value)
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

func evaluateExpression(tokens []Token, variables map[string]int) error {
	stack := Stack{}

	for _, token := range tokens {
		switch token.tokenType {
		case TokenInteger, TokenIdentifier:
			stack.Push(token)
		case TokenBinaryOperator:
			rhs := stack.Pop()
			lhs := stack.Pop()

			var rhsValue, lhsValue int
			var ok bool

			if rhs.tokenType == TokenIdentifier {
				rhsValue, ok = variables[rhs.stringValue()]

				if !ok {
					return fmt.Errorf("unintialized variable %q", rhs.stringValue())
				}
			} else {
				rhsValue = rhs.intValue()
			}

			if lhs.tokenType == TokenIdentifier {
				lhsValue, ok = variables[lhs.stringValue()]

				if !ok {
					return fmt.Errorf("unintialized variable %q", lhs.stringValue())
				}
			} else {
				lhsValue = lhs.intValue()
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
			}

		case TokenAssignment:
			value := stack.Pop().intValue()
			name := stack.Pop().stringValue()
			variables[name] = value

		case TokenUnaryOperator:
			a := stack.Pop().intValue()
			switch token.byteValue() {
			case '-':
				stack.Push(Token{TokenInteger, -1 * a})
			}
		}
	}

	return nil
}
