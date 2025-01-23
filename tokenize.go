package main

import (
	"fmt"
)

type TokenType int

const (
	TokenWhiteSpace     TokenType = iota
	TokenInteger                  = iota
	TokenBinaryOperator           = iota
	TokenUnaryOperator            = iota
	TokenParentheses              = iota
	TokenAssignment               = iota
	TokenIdentifier               = iota
)

type Token struct {
	tokenType TokenType
	value     interface{}
}

func (token Token) byteValue() byte {
	return token.value.(byte)
}

func (token Token) stringValue() string {
	return string(token.value.([]byte))
}

func (token Token) intValue() int {
	return token.value.(int)
}

func Tokenize(line []byte) ([]Token, error) {
	tokens := Stack{}

	token := Token{TokenWhiteSpace, nil}

	for _, char := range line {
		isDigit := '0' <= char && char <= '9'
		isLower := 'a' <= char && char <= 'z'
		isUpper := 'A' <= char && char <= 'Z'
		isAlpha := isLower || isUpper
		isPlusOrMinus := char == '+' || char == '-'
		isBinaryOperator := char == '+' || char == '-' || char == '/' || char == '*' || char == '^' || char == '%'

		if char == ' ' {
			if token.tokenType != TokenWhiteSpace {
				tokens.Push(token)
			}
			token = Token{TokenWhiteSpace, nil}
			continue
		}

		if isDigit {
			if token.tokenType == TokenInteger {
				token.value = token.value.(int)*10 + int(char-'0')
				continue
			}

			if token.tokenType == TokenIdentifier {
				token.value = append(token.value.([]byte), char)
				continue
			}

			if token.tokenType != TokenWhiteSpace {
				tokens.Push(token)
			}
			token = Token{TokenInteger, int(char - '0')}
			continue
		}

		if isAlpha || char == '_' {
			if token.tokenType == TokenIdentifier {
				token.value = append(token.value.([]byte), char)
				continue
			}

			if token.tokenType != TokenWhiteSpace {
				tokens.Push(token)
			}
			token = Token{TokenIdentifier, []byte{char}}
			continue
		}

		if char == '(' || char == ')' {
			if token.tokenType != TokenWhiteSpace {
				tokens.Push(token)
			}
			token = Token{TokenParentheses, byte(char)}
			continue
		}

		isPrevTokenInteger := token.tokenType == TokenInteger
		isPrevTokenIdentifier := token.tokenType == TokenIdentifier
		isPrevTokenParenthesis := token.tokenType == TokenParentheses

		if token.tokenType == TokenWhiteSpace {
			if tokens.Len() >= 1 {
				isPrevTokenInteger = tokens.Peek().tokenType == TokenInteger
				isPrevTokenIdentifier = tokens.Peek().tokenType == TokenIdentifier
				isPrevTokenParenthesis = tokens.Peek().tokenType == TokenParentheses
			} else {
				isPrevTokenInteger = false
				isPrevTokenIdentifier = false
				isPrevTokenParenthesis = false
			}
		}

		if isPlusOrMinus && !(isPrevTokenInteger || isPrevTokenIdentifier || isPrevTokenParenthesis) {
			if token.tokenType != TokenWhiteSpace {
				tokens.Push(token)
			}
			token = Token{TokenUnaryOperator, char}
			continue
		}

		if isBinaryOperator {
			if token.tokenType != TokenWhiteSpace {
				tokens.Push(token)
			}
			token = Token{TokenBinaryOperator, char}
			continue
		}

		if char == '=' {
			if token.tokenType != TokenWhiteSpace {
				tokens.Push(token)
			}
			token = Token{TokenAssignment, char}
			continue
		}

		return tokens.ToArr(), fmt.Errorf("invalid character %q found while tokenizing", string([]byte{char}))
	}

	if token.tokenType != TokenWhiteSpace {
		tokens.Push(token)
	}

	return tokens.ToArr(), nil
}
