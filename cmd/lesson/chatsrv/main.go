package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	go func() {
		var mes string
		for {
			fmt.Fscan(os.Stdin, &mes)
			messages <- fmt.Sprintf("%s %s: %s", time.Now().String()[0:19], "Admin ", mes)
		}

	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)

	}

}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	nikName := conn.RemoteAddr().String()
	inputNikName := make([]byte, 50)
	lenNikName, err := conn.Read(inputNikName)
	if err != nil {
		log.Println("Error read NickName")
	} else {
		nikName = string(inputNikName[0:lenNikName])
	}
	ch <- fmt.Sprintf("%s You are  %s ", time.Now().String()[0:19], nikName)
	messages <- fmt.Sprintf("%s: %s has arrived", time.Now().String()[0:19], nikName)
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- fmt.Sprintf("%s %s: %s", time.Now().String()[0:19], nikName, input.Text())
	}
	leaving <- ch
	messages <- fmt.Sprintf("%s: %s has left", time.Now().String()[0:19], nikName)
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
