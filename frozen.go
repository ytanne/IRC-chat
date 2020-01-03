package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Request struct {
	Person   *User
	RoomName string
}

type User struct {
	Username  string
	Nick   string
	Pw     string
	Output chan Message
	CurrentChatRoom ChatRoom
}

type Message struct {
	Username string
	Text     string
}

type ChatServer struct {
	AddUser     chan User
	AddNick    chan User
	RemoveNick chan User
	NickMap    map[string]User
	Users      map[string]User
	Rooms      map[string]ChatRoom
	Create     chan ChatRoom
	Delete     chan ChatRoom
	UserJoin    chan Request
	UserLeave   chan Request
}

type ChatRoom struct {
	Name  string
	Users map[string]User
	Join  chan User
	Leave chan User
	Input chan Message
}

func joinChannel(server *ChatServer, request Request) {
	if chatRoom, test := server.Rooms[request.RoomName]; test {
		chatRoom.Join <- *(request.Person)
		request.Person.CurrentChatRoom = chatRoom
	} else {
		chatRoome := ChatRoom{
			Name:  request.RoomName,
			Users: make(map[string]User),
			Join:  make(chan User, 4),
			Leave: make(chan User, 4),
			Input: make(chan Message, 4),
		}
		server.Rooms[chatRoome.Name] = chatRoome
		server.Create <- chatRoome
		chatRoome.Join <- *(request.Person)
		request.Person.CurrentChatRoom = chatRoome
	}
}

func (server *ChatServer) Run() {
	for {
		select {
			case user := <-server.RemoveNick:
				delete(server.NickMap, user.Nick)
			case user := <-server.AddNick:
				server.NickMap[user.Nick] = user
			case user := <-server.AddUser:
				server.Users[user.Username] = user
				server.NickMap[user.Nick] = user
			case chatRoom := <-server.Create:
				server.Rooms[chatRoom.Name] = chatRoom
				go chatRoom.manage()
			case chatRoom := <-server.Delete:
				delete(server.Rooms, chatRoom.Name)
			case request := <-server.UserJoin:
				joinChannel(server, request)
			case request := <-server.UserLeave:
				room := server.Rooms[request.RoomName]
				room.Leave <- *(request.Person)
		}
	}
}

func sendToChannel(p User, sender, msg string) {
	p.Output <- Message{
		Username: sender,
		Text:     msg,
	}
}

func (room *ChatRoom) manage() {
	for {
		select {
		case user := <-room.Join:
			room.Users[user.Username] = user
			room.Input <- Message{
				Username: "SYSTEM",
				Text:     fmt.Sprintf("%s joined room %s", user.Nick, room.Name),
			}
		case user := <-room.Leave:
			delete(room.Users, user.Username)
			room.Input <- Message{
				Username: "SYSTEM",
				Text:     fmt.Sprintf("%s left room %s", user.Nick, room.Name),
			}
		case msg := <-room.Input:
			for _, user := range room.Users {
				select {
				case user.Output <- msg:
				default:
				}
			}
		}
	}
}

func getString(s string) string {
	s_ := ""
	l := len(s)
	for i := 0; i < l; i++ {
		if (s[i] >= 32 && s[i] <= 126) {
			s_ += string(s[i])
		}
	}
	return s_
}

func authentication(chatServer *ChatServer, conn net.Conn) User{
	var user User

	io.WriteString(conn, "Enter your username: ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	name := getString(scanner.Text())
	if tmp, test := chatServer.Users[name]; test {
		user = tmp

		io.WriteString(conn, "Enter your password: ")
		scanner.Scan()
		pass := getString(scanner.Text())
		for pass != user.Pw {
			io.WriteString(conn, "Wrong. Try again: ")
			scanner.Scan()
			pass = getString(scanner.Text())
		}

	} else {
		io.WriteString(conn, "Enter nickname: ")
		scanner.Scan()
		nickname := getString(scanner.Text())

		for {
			if _, test := chatServer.NickMap[nickname]; test {
				io.WriteString(conn, "Try again. This nickname is already taken\n")
				scanner.Scan()
				nickname = getString(scanner.Text())
			} else {
				break
			}
		}

		io.WriteString(conn, "Enter a password for your account: ")
		scanner.Scan()
		pass := getString(scanner.Text())
		tmp := User{
			Username:  name,
			Output: make(chan Message, 10),
			Nick:   nickname,
			Pw:     pass,
		}
		chatServer.AddUser <- tmp
		user = tmp
	}
	return (user)
}

func handleConnection(chatServer *ChatServer, conn net.Conn) {
	var user User

	defer conn.Close()
	user = authentication(chatServer, conn)
	scanner := bufio.NewScanner(conn)
	io.WriteString(conn, "Enter a chat Room: ")
	scanner.Scan()

	request := Request{
		Person:   &user,
		RoomName: getString(scanner.Text()),
	}
	chatServer.UserJoin <- request

	defer func() {
		chatServer.UserLeave <- request
	}()

	go func() {
		for scanner.Scan() {
			ln := getString(scanner.Text())
			args := strings.Split(ln, " ")
			if args[0] == "NICK" && len(args) > 1 {
				i := 0
				for _, p := range chatServer.Users {
					if i != 0 {
						break
					} else if p.Nick == args[1] {
						sendToChannel(user, "SYSTEM", "nickname \""+args[1]+"\" taken")
						i = 1
					}
				}

				if _, test := chatServer.NickMap[args[1]]; test {
					i = 2
				}
				if i == 0 {
					chatServer.RemoveNick <- user
					delete(chatServer.NickMap, user.Nick)
					chatServer.NickMap[args[1]] = user
					user.Nick = args[1]
					chatServer.RemoveNick <- user
				}
			} else if ln == "WHOAMI" {
				sendToChannel(user, "SYSTEM\n", "username: "+user.Username+"\nnickname: "+user.Nick+"\ncurrent room: "+user.CurrentChatRoom.Name)
			} else if ln == "NAMES" {
				for person := range chatServer.Users {
					sendToChannel(user, "SYSTEM", person)
				}
			} else if ln == "ROOMMATES" {
				for _, person := range user.CurrentChatRoom.Users {
					sendToChannel(user, "SYSTEM", person.Nick)
				}
			} else if args[0] == "PRIVMSG" && len(args) > 2 {
				if args[1] == "USER" {
					usr, ok := chatServer.Users[args[2]]
					if ok {
						usr.Output <- Message{
							Username: user.Username,
							Text:     fmt.Sprintf("%s", ln),
						}
					} else {
						user.Output <- Message{
							Username: "SYSTEM",
							Text:     fmt.Sprintf("User not found"),
						}
					}
				} else if args[1] == "CHAN" {
					room, ok := chatServer.Rooms[args[2]]
					if ok {
						room.Input <- Message{
							Username: user.Username,
							Text:     ln,
						}
					} else {
						user.Output <- Message{
							Username: user.Username,
							Text:     fmt.Sprintf("Room not found"),
						}
					}
				} else {
					user.Output <- Message{
						Username: "SYSTEM",
						Text:     fmt.Sprintf("Invalid option"),
					}
				}
			} else if ln == "LIST" {
				for room := range chatServer.Rooms {
					sendToChannel(user, "SYSTEM", room)
				}
			} else if args[0] == "JOIN" && len(args) > 1 {
				request = Request{
					Person:   &user,
					RoomName: user.CurrentChatRoom.Name,
				}
				chatServer.UserLeave <- request
				request = Request{
					Person:   &user,
					RoomName: args[1],
				}
				chatServer.UserJoin <- request
			} else if ln == "PART" {
				request = Request{
					Person:   &user,
					RoomName: user.CurrentChatRoom.Name,
				}
				chatServer.UserLeave <- request
				request = Request{
					Person:   &user,
					RoomName: "lobby",
				}
				chatServer.UserJoin <- request
			} else {
				user.CurrentChatRoom.Input <- Message{user.Nick, ln}
			}
		}
	}()

	for msg := range user.Output {
		if msg.Username != user.Username {
			_, err := io.WriteString(conn, msg.Username + ": " + msg.Text + "\n")
			if err != nil {
				break
			}
		}
	}
}

func createServer() *ChatServer {
	chatServer := &ChatServer{
		AddUser:     make(chan User, 4),
		AddNick:    make(chan User, 4),
		RemoveNick: make(chan User, 4),
		NickMap:    make(map[string]User),
		Users:      make(map[string]User),
		Rooms:      make(map[string]ChatRoom),
		Create:     make(chan ChatRoom, 4),
		Delete:     make(chan ChatRoom, 4),
		UserJoin:    make(chan Request, 4),
		UserLeave:   make(chan Request, 4),
	}
	return chatServer
}

func main() {
	server, err := net.Listen("tcp", ":6667")
	defer server.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	chatServer := createServer()	
	go chatServer.Run()
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go handleConnection(chatServer, conn)
	}
}