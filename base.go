package main

type BaseCommand struct {
	server *Server
}

func NewBaseCommand(server *Server) *BaseCommand {
	return &BaseCommand{
		server: server,
	}
}

func (command *BaseCommand) FindArgument(args []*Argument, name string) *Argument {
	for _, arg := range args {
		if arg.Name == name {
			return arg
		}
	}

	return nil
}
