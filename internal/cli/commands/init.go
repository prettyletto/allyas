package commands

import "fmt"

type InitCommand struct {}


func NewInitCommand() *InitCommand  {
	return &InitCommand{}
}

func (c *InitCommand) Names() []string {
	return []string{"init", "-i"}
}

func (c *InitCommand) Usage() string {
	return "init"

}

func (c *InitCommand) Description() string {
	return `Configure the files needed to the app run locally 
	and prepare the file to be injected in the shell`
}

func (c *InitCommand) Execute(ctx CommandContext, args []string)  error { 
	fmt.Println("Not able to create files here")
	return nil
}


