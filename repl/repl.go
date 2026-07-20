package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

func Start(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)

	// Loop indefinitely consuming parsing and interpreting the source code
	for {
		fmt.Fprintf(writer, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Get a line of text from the scanner
		line := scanner.Text()

		// Pass the line of text to the Lexer to tokenize
		lxr := lexer.NewLexer(line)

		// Grab the next token, check for the EOF and print out until there are no more tokens
		for tok := lxr.NextToken(); tok.Type != token.EOF; tok = lxr.NextToken() {
			fmt.Fprintf(writer, "%+v \n", tok)
		}

	}
}
