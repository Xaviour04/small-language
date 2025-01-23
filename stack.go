package main

type Stack struct {
	arr []Token
}

func (stack *Stack) Push(token Token) {
	stack.arr = append(stack.arr, token)
}

func (stack *Stack) Pop() Token {
	length := len(stack.arr)
	elem := stack.arr[length-1]
	stack.arr = stack.arr[:length-1]
	return elem
}

func (stack *Stack) Peek() Token {
	return stack.arr[len(stack.arr)-1]
}

func (stack *Stack) Len() int {
	return len(stack.arr)
}

func (stack Stack) ToArr() []Token {
	return stack.arr
}
