package main

import (
	"fmt"
	"github.com/alexeykirinyuk/take-smaller-tasks-tool/command"
)

func main() {
	res, err := command.Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(res)
}
