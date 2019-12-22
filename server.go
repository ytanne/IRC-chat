package main

import (
    "fmt"
	"net"
	"os"
	//"strings"
)

type user struct {
	password string
	username string
	nickname string
}

func	main() {

	ln, err := net.Listen("tcp", ":8081")
	errorHandling(err)
	for {
		conn, err := ln.Accept()
		errorHandling(err)
		go handleConnection(conn)
	}

}

func	handleConnection(conn net.Conn) {
	var	credentials[3] string

	conn.Write([]byte("Enter username: "))
	credentials[0] = getData(conn)
	conn.Write([]byte("Enter password: "))
	credentials[1] = getData(conn)
	conn.Write([]byte("Enter nickname: "))
	credentials[2] = getData(conn)
	//fmt.Println(conn.RemoteAddr().String())
	conn.Close()
}

func	errorHandling(err error) {
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
}

func	getData(conn net.Conn) string {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	fmt.Printf("Length of the message is %d\n", reqLen)
	errorHandling(err)
	fmt.Printf("Message contents: %q\n", string(buf[:reqLen]))
	//conn.Write([]byte("Message received.\n"))
	return (string(buf[:reqLen]))
}