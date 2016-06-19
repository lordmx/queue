package main

import (
	"bytes"
	"fmt"
)

type Body []byte

type Argument struct {
	Name  string
	Value []byte
}

type Result struct {
	Output []byte
	Error  error
}

type RawCommand struct {
	Client *Client
	Body   []byte
}

type Command interface {
	GetName() string
	Run(client *Client, args []*Argument) *Result
	Parse(body []byte) ([]*Argument, error)
}

type Protocol struct {
	Commands map[string]Command
}

func NewProtocol() *Protocol {
	return &Protocol{
		Commands: make(map[string]Command),
	}
}

func (proto *Protocol) RegisterCommand(command Command) {
	proto.Commands[command.GetName()] = command
}

func (proto *Protocol) FindCommand(name string) (Command, error) {
	if command, ok := proto.Commands[name]; ok {
		return command, nil
	}

	return nil, fmt.Errorf("Could not found command")
}

func (proto *Protocol) Parse(body []byte) (Command, []*Argument, error) {
	arguments := []*Argument{}
	parts := bytes.SplitN(body, []byte(" "), 1)
	commandName, body := parts[0], parts[1]
	command, err := proto.FindCommand(string(commandName))

	if err != nil {
		return nil, arguments, err
	}

	arguments, err = command.Parse(body)

	if err != nil {
		return nil, arguments, err
	}

	return command, arguments, nil
}
