package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type cmdFlags struct {
	Add    string
	Del   int 
	Edit   string
	Toggle int 
	List   bool
  Quit   string
}

func NewCmdFlags() *cmdFlags {
	cf := cmdFlags{}
	flag.StringVar(&cf.Add, "add", "", "Add a new todo specify title")
  	flag.StringVar(&cf.Edit, "Edit", "", "Edit todo")
    flag.IntVar(&cf.Del, "del", -1, "Delete todo")
    flag.IntVar(&cf.Toggle, "toggle", -1, "Toggle todo")
    flag.BoolVar(&cf.List, "list", false, "List todos")
    flag.StringVar(&cf.Quit, "quit", "", "Quit todos")
    flag.Parse()
  
return &cf 
}

func (cf *cmdFlags) execute(todos *Todos){
  switch {
  case cf.List:
    todos.print()

  case cf.Add != "":
    todos.add(cf.Add)
    todos.print()
    
case cf.Edit != "":
  parts := strings.SplitN(cf.Edit, ":", 2)
  if len(parts) != 2{
    fmt.Println("Invalid Format for edit")
    os.Exit(1)
  }
  index , err := strconv.Atoi(parts[0])
  if err != nil {
    fmt.Println("Error: invalid index")
    os.Exit(1)
  }
  todos.edit(index, parts[1] )
  case cf.Toggle != -1: 
  todos.toggle(cf.Toggle)
  todos.print()

case cf.Del != -1:
  todos.delete(cf.Del)
  todos.print()
case cf.Quit != "":

default:
  fmt.Println("Invalid Command")
}

}


