package main

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	clients []net.Conn
	mutex   sync.Mutex
	users   = make(map[net.Conn]string)
	user    string
)

func main() {
	listener, err := net.Listen("tcp", ":8989")
	if err != nil {
		log.Println("error listening", err)
	}
	log.Println("server started at port 8989")
	defer listener.Close()
	for {
		conn, err2 := listener.Accept()
		if err2 != nil {
			log.Println("error accepting", err)
			continue
		}
		mutex.Lock()
		clients = append(clients, conn)
		mutex.Unlock()

		go Handle(conn)
	}
}

func Handle(conn net.Conn) {
	conn.Write([]byte("Welcome to TCP-Chat!\n         _nnnn_\n        dGGGGMMb\n       @p~qp~~qMb\n       M|@||@) M|\n       @,----.JM|\n      JS^\\__/  qKL\n     dZP        qKRb\n    dZP          qKKb\n   fZP            SMMb\n   HZM            MMMM\n   FqM            MMMM\n __| \".        |\\dS\"qML\n |    `.       | `' \\Zq\n_)      \\.___.,|     .'\n\\____   )MMMMMP|   .'\n     `-'       `--'\n[ENTER YOUR NAME]: "))
	bufferName := make([]byte, 1024)
	m, _ := conn.Read(bufferName)

	username := bufferName[:m]
	users[conn] = strings.TrimSpace(string(username))

	clientgone := "has disconnected"
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		message := buffer[:n]
		if err != nil {
			if err == io.EOF {
				log.Println(clientgone)
				mutex.Lock()
				for i, client := range clients {
					if client == conn {
						clients = append(clients[:i], clients[i+1:]...)
						break
					}
				}
				for _, client := range clients {
					if conn == client {
						user = users[conn]
					}
					client.Write([]byte(user + " " + clientgone))
				}

				mutex.Unlock()
				conn.Close()
				break

			}

			log.Println("error reading", err)
			return
		}
		log.Println("message received", string(message))

		mutex.Lock()

		for _, client := range clients {
			if conn == client {
				user = users[conn]
				continue
			}
			_, err3 := client.Write([]byte(user + ":" + " " + string(message)))
			if err3 != nil {
				log.Println("error sending message to client", err)
			}

		}
		mutex.Unlock()

	}
}
