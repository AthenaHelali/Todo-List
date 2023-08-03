package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"todo-list/delivery/deliveryparam"
)

func main() {
	var message string
	message = os.Args[1]
	connection, err := net.Dial("tcp", "127.0.0.1:2022")
	if err != nil {
		log.Fatalln("can't dial the given address", err)
	}
	defer connection.Close()

	fmt.Println("local address", connection.LocalAddr())

	req := deliveryparam.Request{
		Command: message,
	}
	if req.Command == "create-task" {
		req.CreateTaskRequest = deliveryparam.CreateTaskRequest{
			Title:      "test",
			DueDate:    "test",
			CategoryID: 1,
		}
	}
	serializedData, mErr := json.Marshal(req)
	if mErr != nil {
		log.Fatalln("can't marshal request", mErr)
	}

	numberOfWrittenByte, wErr := connection.Write(serializedData)
	if wErr != nil {
		log.Fatalln("can't write data to connection", wErr)
	}

	fmt.Println(numberOfWrittenByte)

	var data = make([]byte, 1024)
	_, rErr := connection.Read(data)
	if rErr != nil {
		log.Fatalln("can't read data from connection", rErr)
	}
	fmt.Println("server response:", string(data))

}
