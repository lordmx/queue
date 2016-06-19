package main

import (
	"bytes"
	"fmt"
)

type PublishCommand struct {
	*BaseCommand
}

func NewPublishCommand(server *Server) *PublishCommand {
	return &PublishCommand{
		BaseCommand: NewBaseCommand(server),
	}
}

func (command *PublishCommand) GetName() string {
	return "publish"
}

func (command *PublishCommand) Parse(body []byte) ([]*Argument, error) {
	parts := bytes.SplitN(body, []byte(" "), 1)

	if len(parts) < 2 {
		return nil, fmt.Errorf("Empty or missing body of message")
	}

	queueName, body := parts[0], parts[1]
	arguments := []*Argument{
		{
			Name:  "queue",
			Value: queueName,
		},
		{
			Name:  "message",
			Value: body,
		},
	}

	return arguments, nil
}

func (command *PublishCommand) Run(client *Client, args []*Argument) *Result {
	result := &Result{}
	queueName := command.FindArgument(args, "queueName")

	if queueName == nil {
		result.Error = fmt.Errorf("Empty or missing queue name")
		return result
	}

	queueNameStr := string(queueName.Value)
	queue, err := command.server.FindQueue(queueNameStr, true)

	if err != nil {
		result.Error = err
		return result
	}

	body := command.FindArgument(args, "body")

	if body == nil && len(body.Value) == 0 {
		result.Error = fmt.Errorf("Empty or missing body of message")
		return result
	}

	result.Output = body.Value
	queue.Publish <- NewMessage(queue, client, body.Value)

	return result
}
