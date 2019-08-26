package async

import (
	"bufio"
	"io"
	"log"
)

type scanner struct{
	reader io.Reader
	text chan string
}

func NewScanner (reader io.Reader) scanner {
	return scanner{reader, make(chan string)}
}

func (scan scanner) Doing() {
	scanner := bufio.NewScanner(scan.reader)
	for scanner.Scan(){
		scan.text <- scanner.Text()
	}
	if err := scanner.Err();
		err != nil {
		log.Println(err)
	}
}


func (scanner scanner) Text() <- chan string {
	return scanner.text
}
