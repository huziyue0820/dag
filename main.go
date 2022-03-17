package main

import (
	"dag/CLI"
	"dag/blockchain"
	"fmt"
)

func main() {
	fmt.Println("dag begin...")

	bc := blockchain.NewBlockchain()
	defer bc.DB.Close()

	cli := CLI.CLI{BC: bc}
	cli.Run()
}
