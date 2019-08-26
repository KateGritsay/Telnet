package telnet

import "io"

type scanner struct{
	reader io.Reader
	text chan string
}

func NewScanner (reader io.Reader) scanner {
	return scanner{reader, make(chan string)}
}

func (scanner scanner) Text() <- chan string {
	return scanner.text
}
