package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Application start up call made with number of arguments: %s and arguments: %s",len(os.Args),os.Args)
	if len(os.Args)>1 {
		startRouter(os.Args[1])
	} else {
		startRouter("")
	}
}







