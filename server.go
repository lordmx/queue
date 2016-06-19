package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Server struct {
	Protocol     *Protocol
	Queues       map[string]*Queue
	Clients      map[string]*Client
	clientsCount int
	Joins        chan net.Conn
	Income       chan *RawCommand
	Results      chan *Result
	Errors       chan error
}

func NewServer(protocol *Protocol) *Server {
	return &Server{
		Protocol: protocol,
		Queues:   make(map[string]*Queue),
		Clients:  make(map[string]*Client),
		Results:  make(chan *Result),
		Errors:   make(chan error),
		Joins:    make(chan net.Conn),
		Income:   make(chan *RawCommand),
	}
}

func (server *Server) FindQueue(name string, orCreate bool) (*Queue, error) {
	if queue, ok := server.Queues[name]; ok {
		return queue, nil
	}

	if orCreate {
		return server.CreateQueue(name)
	}

	return nil, fmt.Errorf("Could not found the queue")
}

func (server *Server) CreateQueue(name string) (*Queue, error) {
	if _, ok := server.Queues[name]; !ok {
		queue := NewQueue(server, name)
		server.Queues[name] = queue

		return queue, nil
	}

	return nil, fmt.Errorf("Queue already exists")
}

func (server *Server) Run(host string, port int) error {
	go func() {
		for {
			select {
			case err := <-server.Errors:
				log.Printf("[%v] (!) %s", time.Now(), err.Error())
			case conn := <-server.Joins:
				server.clientConnected(conn)
			case raw := <-server.Income:
				command, arguments, err := server.Protocol.Parse(raw.Body)

				if err != nil {
					server.Errors <- err
				} else {
					go func() {
						server.Results <- command.Run(raw.Client, arguments)
					}()
				}
			case result := <-server.Results:
				if result.Error != nil {
					server.Errors <- result.Error
				}
			}
		}
	}()

	host = fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			server.Errors <- err
			continue
		}

		server.Joins <- conn
	}
}

func (server *Server) sendToClient(client *Client, data []byte) {
	_, err := client.conn.Write(data)

	if err != nil {
		server.Errors <- err
	}
}

func (server *Server) clientConnected(conn net.Conn) {
	server.clientsCount++

	client := NewClient(server, conn)
	server.Clients[client.ID] = client

	go client.read()
}

func (server *Server) clientDisconnected(client *Client, err error) {
	server.Errors <- err

	delete(server.Clients, client.ID)
	server.clientsCount--
}
