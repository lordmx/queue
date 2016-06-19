package main

import (
	"bufio"
	"net"
)

type Client struct {
	ID     string
	conn   net.Conn
	server *Server
}

func NewClient(server *Server, conn net.Conn) *Client {
	return &Client{
		ID:     RandString(32),
		server: server,
		conn:   conn,
	}
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) read() {
	reader := bufio.NewReader(client.conn)
	server := client.server

	for {
		data, err := reader.ReadBytes('\n')

		if err != nil {
			server.clientDisconnected(client, err)

			break
		}

		server.Income <- &RawCommand{client, data}
	}

	client.Close()
}
