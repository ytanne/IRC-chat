# Frozen (IRC Server project)
### Description
In this rush, we had the chance to (re)discover the amazing world of IRC in a new language, Go.
According to the definition of Wikipedia, an Internet Relay Chat (IRC) is an application layer protocol that facilitates communication in the form of text. The chat process works on a client/server networking model. IRC clients are computer programs that a user can install on their system. These clients communicate with chat servers to transfer messages to other clients. IRC is mainly designed for group communication in discussion forums, called channels, but also allows one-on-one communication via private messages as well as chat and data transfer, including file sharing.
Therefore, an IRC server was a great example of a concurrent program, and recreating one in Go language, which is very suitable for concurrency, is the perfect opportunity to take programming skills to the next level.

### Installation
Firstly, clone the repository:
```
git clone https://github.com/ytanne/IRC-chat && cd IRC-chat
```
Run server using
```
go run frozen.go
```
The server will be running on `localhost` (or IP of computer running server) on port `6667`. You can check the functionality of the server by accessing it using netcat:
```
nc localhost 6667
```

### Features
* NICK - Change nickname
* JOIN - Makes the user join a channel. If the channel doesn’t exist, it will be created.
* PART - Makes the user leave a channel.
* NAMES - Lists all users connected to the server (bonus: make it RFC com-
pliant with channel modes).
* LIST - Lists all channels in the server (bonus: make it RFC compliant with channel modes).
* PRIVMSG - Send a message to another user or a channel.



### The project is still on development stage
While waiting, please listen for the [Let it go (Rock version)](https://www.youtube.com/watch?v=GG31XuWPQQ4).
I will try to finish this project when my hands will be free.
Thanks for attention.
<img src="https://media0.giphy.com/media/Fjy5XItIvYjEQ/giphy.gif">
