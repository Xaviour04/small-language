package main

import (
	"bytes"
	"log"
	"os"
)

func ExitOn(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func DebugPrintCommand(data [][]byte) {
	str := ""
	for i := 0; i < len(data); i++ {
		str += string(data[i]) + " "
	}
	log.Println(str)
}

func Evaluate(data []byte) {
	lines := bytes.Split(data, []byte("\n"))

	variables := make(map[string]int)

	for lineNo, line := range lines {
		line = bytes.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		tokens, err := Tokenize(line)

		if err != nil {
			log.Println(err.Error())
			log.Fatalf("Tokenization error on line %d", lineNo)
		}

		tokens, err = convertExpressionToPostfix(tokens)

		if err != nil {
			log.Println(err.Error())
			log.Fatalf("Tokenization error on line %d", lineNo)
		}

		evaluateExpression(tokens, variables)
	}

	for k, v := range variables {
		log.Printf("%s = %d", k, v)
	}
}

func main() {
	if len(os.Args) < 1 {
		log.Println("Please give a file to run")
		return
	}

	filename := os.Args[1]

	data, err := os.ReadFile(filename)
	ExitOn(err)

	Evaluate(data)

}
