package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var nikName string

func main() {
	fmt.Println("Введите свой NikName")
	fmt.Scan(&nikName)
	conn, err := net.Dial("tcp", "localhost:8000")

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Write([]byte(nikName))
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		_, err = io.Copy(os.Stdout, conn)
		if err != nil {
			log.Fatal(err)
		}
	}()
	_, err = io.Copy(conn, os.Stdin) // until you send ^Z
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: exit", conn.LocalAddr())
}
