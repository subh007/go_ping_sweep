package main

import (
	"fmt"
	"github.com/subh007/goodl/go_ping_sweep"
)

func main() {
	t := new(go_ping_sweep.Table)
	t.SetTitle("Programmming")
	t.SetHeader("rank", "lang")
	if t.AddData("1", "go") != nil {
		fmt.Println("some error happened")
	} else {
		fmt.Println("row added")
	}

	t.CreateTable()
}
