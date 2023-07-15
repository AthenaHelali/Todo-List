package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type User struct {
	ID       int
	Email    string
	Password string
}

var userStorage []User

func main() {
	fmt.Println("Hello to TODO app")

	command := flag.String("command", "no-command", "command to run")
	flag.Parse()
	for {
		runCommand(*command)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
		fmt.Printf("userStorage: %v\n", userStorage)
	}

}
func runCommand(command string) {
	switch command {
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser()
	case "login":
		login()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command is not valid", command)

	}
}
func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var name, duedate, category string

	fmt.Println("please enter the task title")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter the task category")
	scanner.Scan()
	category = scanner.Text()

	fmt.Println("please enter the task due date")
	scanner.Scan()
	duedate = scanner.Text()

	fmt.Println("task:", name, category, duedate)

}
func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Println("please enter the category title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the category color")
	scanner.Scan()
	color = scanner.Text()

	fmt.Println("category:", title, color)

}
func registerUser() {
	scanner := bufio.NewScanner(os.Stdin)
	var email, password string
	fmt.Println("please enter the user email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the user password")
	scanner.Scan()
	password = scanner.Text()

	fmt.Printf("user:%v %v %v \n", len(userStorage)+1, email, password)

	user := User{
		ID:       len(userStorage) + 1,
		Email:    email,
		Password: password,
	}
	userStorage = append(userStorage, user)

}
func login() {
	var id, email, password string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("please enter the user email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the user password")
	scanner.Scan()
	password = scanner.Text()

	fmt.Println("user", id, email, password)

}
