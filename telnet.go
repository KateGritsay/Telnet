package telnet

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
	"github.com/KateGritsay/Telnet/async"
)

func main() {
	adress := flag.String("adr", "127.0.0.1:4242", "adress for connect")
	timeout := flag.Int("timeout", 60, "")
	flag.Parse()

dialer := &net.Dialer{}
ctx := context.Background()

ctx, cancel := context.WithTimeout(ctx, time.Duration(*timeout) * time.Second)
conn, err := dialer.DialContext(ctx, "tcp", *adress)
if err != nil {
log.Fatalf("Cannot connect: %v", err)
}



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


cancel()
wg.Wait()
conn.Close()


}

func read (ctx context.Context, conn net.Conn) {
	scanner := async.NewScanner(conn)
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
	scanner := async.NewScanner(os.Stdin)
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



