package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"todo-list/delivery/deliveryparam"
	"todo-list/repository/memorystore"
	"todo-list/service/task"
)

func main() {
	const (
		network = "tcp"
		address = "127.0.0.1:2022"
	)
	// create new listener
	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatalln("can't listen on given address", address, network)
	}
	defer listener.Close()

	fmt.Println("server listening on", listener.Addr())
	taskMemoryRepo := memorystore.NewTaskStore()
	taskService := task.NewService(taskMemoryRepo)
	for {
		// listen for new connection
		connection, aErr := listener.Accept()
		if aErr != nil {
			log.Println("can't listen to new connection", aErr)

			continue
		}

		fmt.Println("client address", connection.RemoteAddr())

		//process request

		var rawReq = make([]byte, 1024)
		numberOfReadBytes, rErr := connection.Read(rawReq)
		if rErr != nil {
			log.Println("can't read data from connection", rErr)
			continue
		}
		req := &deliveryparam.Request{}
		if uErr := json.Unmarshal(rawReq[:numberOfReadBytes], req); uErr != nil {
			log.Println("bad request", err)

			continue
		}
		switch req.Command {
		case "create-task":
			response, cErr := taskService.CreateTask(task.CreateTaskRequest{
				Title:               req.CreateTaskRequest.Title,
				DueDate:             req.CreateTaskRequest.DueDate,
				CategoryID:          req.CreateTaskRequest.CategoryID,
				AuthenticatedUserID: 0,
			})
			if cErr != nil {
				_, wErr := connection.Write([]byte(cErr.Error()))
				if wErr != nil {
					log.Println("can't write data to connection", wErr)

					continue
				}
			}
			data, mErr := json.Marshal(&response)
			if mErr != nil {
				_, wErr := connection.Write([]byte(cErr.Error()))
				if wErr != nil {
					log.Println("can't write data to connection", wErr)

					continue
				}

				continue
			}
			log.Printf("data:%+v", response)
			_, wErr := connection.Write(data)
			if wErr != nil {
				log.Println("can't write data to connection", wErr)

				continue
			}
		default:
			println("default")

		}
	}

}
