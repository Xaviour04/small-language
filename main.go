package main

import (
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
	program, err := tokenizeSourceCode(data)

	if err != nil {
		log.Println("error: while tokenizing code")
		log.Fatalln(err)
	}

	for program.IsRunning() {
		err = program.RunLine()
		if err != nil {
			log.Printf("error: line %d", program.lineNo+1)
			log.Fatalln(err)
		}
	}

	for k, v := range program.variables {
		log.Printf("%s = %d", k, v.intValue())
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
