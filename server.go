package main

import (
    "fmt"
	"net"
	"os"
	"encoding/json"
	//"bufio"
	//"io"

	"github.com/gorilla/websocket"
	"strings"
)

type User struct {
	Password 	string
	Username 	string
	Nickname 	string
}

type Message struct {
	Text string `json:"text"`
}

type hub struct {
	clients          map[string]*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan    chan Message
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
	getUserCredentials(conn)
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

func	getUserCredentials(conn net.Conn) {
	var	person User

	var command = getData(conn)
	person.Password = checkCommand(command, "PASS")
	person.Password = person.Password[:len(person.Password) - 1]
	if (person.Password == "") {
		go closeConn(conn, person)
	}
	command = getData(conn)
	person.Nickname = checkCommand(command, "NICK")
	person.Nickname = person.Nickname[:len(person.Nickname) - 1]
	if (person.Password == "") {
		go closeConn(conn, person)
	}
	command = getData(conn)
	person.Username = checkCommand(command, "USER")
	person.Username = person.Username[:len(person.Username) - 1]
	if (person.Password == "") {
		go closeConn(conn, person)
	}
	if (person != User{}) {
		b, err := json.Marshal(person)
		errorHandling(err)
		saveCredentials(b)
	}
	/*
	conn.Write([]byte("Enter username: "))
	person.Username = getData(conn)
	conn.Write([]byte("Enter password: "))
	person.Password = getData(conn)
	conn.Write([]byte("Enter nickname: "))
	person.Nickname = getData(conn)
	checkCredentials(person)
	b, err := json.Marshal(person)
	errorHandling(err)
	saveCredentials(b)
	*/
}

func	closeConn(conn net.Conn, person User) {
	conn.Write([]byte("Wrong input\n"))
	person.Password = ""
	person.Nickname = ""
	person.Username = ""
}

func	checkCommand(input string, goal string) string {
	var array = strings.Split(input, " ")
	
	if array != nil && len(array) == 2 && array[0] == goal {
		return array[1]
	} else {
		return ""
	}
}


func	saveCredentials(credentials []byte) {
	f, err := os.OpenFile("users.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorHandling(err)
	}
	defer f.Close()
	if _, err := f.WriteString(string(credentials) + "\r\n"); err != nil {
		errorHandling(err)
	}
}
/*
func	checkCredentials(person User) bool {

	file, err := os.Open("users.log")
	if err != nil {
		errorHandling(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			return (false)
		}
		if err := json.Unmarshal(line, &person); err != nil {
			errorHandling(err)
		}
		fmt.Printf("%s \n", line)
	}
}
*/