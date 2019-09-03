package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/KateGritsay/Telnet/scaner"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	adress := flag.String("adr", "127.0.0.1:4242", "adress for connect")
	timeout := flag.Int("timeout", 60, "")
	flag.Parse()

	dialer := &net.Dialer{}
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Duration(*timeout)*time.Second)
	conn, err := dialer.DialContext(ctx, "tcp", *adress)
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, syscall.SIGINT)
	go func() {
		<-exitCh
		log.Println("Terminating the process with signal")
		cancel()
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		read(ctx, conn)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		write(ctx, conn)
		wg.Done()
	}()

	wg.Wait()
	conn.Close()

}

func read(ctx context.Context, conn net.Conn) {
	scanner := scaner.NewScanner(conn)
	go scanner.Doing()
	for {
		select {
		case <-ctx.Done():
			return
		case text := <-scanner.Text():
			log.Println(text)
		}
	}
}

func write(ctx context.Context, conn net.Conn) {
	scanner := scaner.NewScanner(os.Stdin)
	go scanner.Doing()
	for {
		select {
		case <-ctx.Done():
			return

		case text := <-scanner.Text():
			text = fmt.Sprintf("%s\n", text)
			_, err := conn.Write([]byte(text))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
