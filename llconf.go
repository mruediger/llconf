package main

import (
	"fmt"
	"os"
	"bufio"
	"github.com/mruediger/llconf/parse"
	"github.com/mruediger/llconf/promise"
)

func main() {
	in := os.Stdin
	bufin := bufio.NewReader( in )

	goal := ""
	if len(os.Args) > 1 {
		goal = os.Args[1]
		fmt.Println("executing goal " + goal)
	} else {
		fmt.Println("need goal")
		os.Exit(1)
	}
	
	promises := parse.ParsePromises( bufin )
	for k,v := range(promises) {
		if k == goal {
			fmt.Printf("evaluating %s\n", v)
			if v.Eval([]promise.Constant{}) {
				fmt.Println("success")
			} else {
				fmt.Println("failure")
			}
		} 
	}
}
