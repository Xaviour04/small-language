package main

import "errors"

type IntStack struct {
	arr []int
}

func (stack *IntStack) Push(token int) {
	stack.arr = append(stack.arr, token)
}

func (stack *IntStack) Pop() (int, error) {
	length := len(stack.arr)

	if length == 0 {
		return -1, errors.New("cannot pop from empty stack")
	}

	elem := stack.arr[length-1]
	stack.arr = stack.arr[:length-1]
	return elem, nil
}

func (stack *IntStack) Len() int {
	return len(stack.arr)
}

type TokenStack struct {
	arr []Token
}

func (stack TokenStack) ToString() string {
	result := ""
	for _, token := range stack.arr {
		result += token.ToString() + " "
	}
	return result
}

func (stack *TokenStack) Push(token Token) {
	stack.arr = append(stack.arr, token)
}

func (stack *TokenStack) Pop() (Token, error) {
	length := len(stack.arr)

	if length == 0 {
		return Token{}, errors.New("cannot pop from empty stack")
	}

	elem := stack.arr[length-1]
	stack.arr = stack.arr[:length-1]
	return elem, nil
}

func (stack *TokenStack) Peek() (Token, error) {
	length := len(stack.arr)

	if length == 0 {
		return Token{}, errors.New("cannot peek into empty stack")
	}

	return stack.arr[length-1], nil
}

func (stack *TokenStack) Len() int {
	return len(stack.arr)
}

func (stack TokenStack) ToArr() []Token {
	return stack.arr
}
