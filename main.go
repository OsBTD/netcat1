package main

import (
	"io"
	"log"
	"net"
	"sync"
)

var (
	clients []net.Conn
	mutex   sync.Mutex
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

	clientgone := "client has disconnected"
	// bufferName := make([]byte, 1024)
	// var users map[string]net.Conn
	// m, _ := conn.Read(bufferName)

	// username := bufferName[:m]
	// users = append(users, username)
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
					client.Write([]byte(clientgone))
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
				continue
			}
			_, err3 := client.Write([]byte(string(message)))
			if err3 != nil {
				log.Println("error sending message to client", err)
			}

		}
		mutex.Unlock()

	}
}
