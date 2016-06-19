package main

import (
	"fmt"
)

type ConsumeCommand struct {
	*BaseCommand
}

func NewConsumeCommand(server *Server) *ConsumeCommand {
	return &ConsumeCommand{
		BaseCommand: NewBaseCommand(server),
	}
}

func (command *ConsumeCommand) GetName() string {
	return "consume"
}

func (command *ConsumeCommand) Parse(body []byte) ([]*Argument, error) {
	arguments := []*Argument{
		{
			Name:  "queue",
			Value: body,
		},
	}

	return arguments, nil
}

func (command *ConsumeCommand) Run(client *Client, args []*Argument) *Result {
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

	result.Output = []byte("true")
	queue.Consume <- client

	return result
}
