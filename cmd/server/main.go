package main

import (
	"AAStarCommunity/EthPaymaster_BackService/router"
	"fmt"
)

func init() {
	//init global variables when service start

}

func main() {
	//use InitRouter
	router.InitRouter()
	fmt.Printf("Server now running on 0.0.0.0:%d", 8080)
	router.Engine.Run(":8080")
}
