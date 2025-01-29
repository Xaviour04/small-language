package main

import (
	"errors"
	"fmt"
)

type TokenType int

const (
	TokenWhiteSpace          TokenType = iota
	TokenInteger                       = iota
	TokenBinaryOperator                = iota
	TokenUnaryOperator                 = iota
	TokenConditionalOperator           = iota
	TokenParentheses                   = iota
	TokenAssignment                    = iota
	TokenIdentifier                    = iota
	TokenBoolean                       = iota
	TokenKeyword                       = iota
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

func (token Token) ToString() string {
	switch token.tokenType {
	case TokenWhiteSpace:
		if token.byteValue() == '\t' {
			return "<indentation>"
		}
		return "<blank>"
	case TokenInteger:
		return fmt.Sprintf("<int:%d>", token.value)
	case TokenBinaryOperator:
		return fmt.Sprintf("<bi-op:%c>", token.value)
	case TokenUnaryOperator:
		return fmt.Sprintf("<uni-op:%c>", token.value)
	case TokenConditionalOperator:
		return fmt.Sprintf("<cond-op:%c>", token.value)
	case TokenParentheses:
		return fmt.Sprintf("<paren:%c>", token.value)
	case TokenAssignment:
		return "<=>"
	case TokenIdentifier:
		return fmt.Sprintf("<var:%s>", token.stringValue())
	case TokenBoolean:
		return fmt.Sprintf("<bool:%t>", token.value)
	case TokenKeyword:
		return fmt.Sprintf("<keyword:%s>", token.stringValue())
	}
	return "<unknown>"
}

func TokenizeLine(line []byte) ([]Token, error) {
	tokens := TokenStack{}

	token := Token{TokenWhiteSpace, nil}

	i := 0
	intendation := line[0]
	for i < len(line) {
		if line[i] == ' ' && intendation == ' ' {
			if line[i+1] == ' ' && line[i+2] == ' ' && line[i+3] == ' ' {
				tokens.Push(Token{TokenWhiteSpace, byte('\t')})
				i += 4
				continue
			}
			return nil, errors.New("invalid number of spaces for intendation")
		}

		if line[i] == '\t' && intendation == '\t' {
			tokens.Push(Token{TokenWhiteSpace, byte('\t')})
			i++
			continue
		}

		break
	}

	for i < len(line) {
		char := line[i]
		i++
		isDigit := '0' <= char && char <= '9'
		isLower := 'a' <= char && char <= 'z'
		isUpper := 'A' <= char && char <= 'Z'
		isAlpha := isLower || isUpper
		isPlusOrMinus := char == '+' || char == '-'
		isBinaryOperator := char == '+' || char == '-' || char == '/' || char == '*' || char == '^' || char == '%'
		isConditionalOperator := char == '>' || char == '<' || char == '='

		if char == ' ' || char == '\n' || char == '\r' {
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
				topToken, _ := tokens.Peek()
				isPrevTokenInteger = topToken.tokenType == TokenInteger
				isPrevTokenIdentifier = topToken.tokenType == TokenIdentifier
				isPrevTokenParenthesis = topToken.tokenType == TokenParentheses
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
			if token.tokenType == TokenConditionalOperator {
				token.value = byte('g') // g -> >=
				if token.value.(byte) == '<' {
					token.value = byte('l') // l -> <=
				}
				continue
			}

			if token.tokenType != TokenWhiteSpace {
				tokens.Push(token)
			}

			if len(line) > i && line[i] == '=' {
				token = Token{TokenConditionalOperator, byte('e')} // e -> ==
				continue
			}

			token = Token{TokenAssignment, char}
			continue
		}

		if isConditionalOperator {
			if token.tokenType != TokenWhiteSpace {
				tokens.Push(token)
			}

			token = Token{TokenConditionalOperator, char}
			continue
		}

		return tokens.ToArr(), fmt.Errorf("invalid character %q found while tokenizing", string([]byte{char}))
	}

	if token.tokenType != TokenWhiteSpace {
		tokens.Push(token)
	}

	tokensArr := tokens.ToArr()

	for i, token := range tokensArr {
		if token.tokenType != TokenIdentifier {
			continue
		}

		if token.stringValue() == "while" {
			tokensArr[i].tokenType = TokenKeyword
		}
	}

	return tokensArr, nil
}
