package main

type Queue struct {
	server    *Server
	Name      string
	Messages  map[string]*Message
	Publish   chan *Message
	Consume   chan *Client
	Consumers map[string]*Client
}

func NewQueue(server *Server, name string) *Queue {
	queue := &Queue{
		server:    server,
		Name:      name,
		Messages:  make(map[string]*Message),
		Publish:   make(chan *Message),
		Consumers: make(map[string]*Client),
		Consume:   make(chan *Client),
	}

	go func() {
		for {
			select {
			case client := <-queue.Consume:
				queue.Consumers[client.ID] = client

				if len(queue.Messages) > 0 && len(queue.Consumers) == 1 {
					for id, message := range queue.Messages {
						go server.sendToClient(client, message.Body)
						delete(queue.Messages, id)
					}
				}
			case message := <-queue.Publish:
				queue.Messages[message.ID] = message

				if len(queue.Consumers) > 0 {
					for _, client := range queue.Consumers {
						go server.sendToClient(client, message.Body)
					}

					delete(queue.Messages, message.ID)
				}
			}
		}
	}()

	return queue
}
