package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		repl()
	case 1:
		if err := runFile(args[0]); err != nil {
			fmt.Printf("error: %s\n", err)
		}
	default:
		fmt.Println("usage")
	}
}

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			if err := interpret(scanner.Text()); err != nil {
				fmt.Printf("error: %s\n", err)
			}
		}
	}
}

func runFile(filename string) error {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return interpret(string(source))
}

func interpret(source string) error {
	chunk, err := newCompiler().compile(source)
	if err != nil {
		return err
	}
	return newVM().run(chunk)
}
