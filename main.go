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
	clientgone := "client has disconnected"
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
