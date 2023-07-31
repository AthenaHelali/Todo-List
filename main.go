package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"todo-list/constant"
	"todo-list/contract"
	"todo-list/filestore"
	"todo-list/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	authenticatedUser *models.User
	serializationMode string
	userStorage       []models.User
	taskStorage       []models.Task
	categoryStorage   []models.Category
)

var userFileStore filestore.FileStorage

func main() {
	fmt.Println("Hello to TODO app")

	serializeModeInput := flag.String("serialize-mode", constant.ManualSerializationMode, "serialization mode to write data to file")
	command := flag.String("command", "no-command", "command to run")
	flag.Parse()

	switch *serializeModeInput {
	case constant.ManualSerializationMode:
		serializationMode = constant.ManualSerializationMode
	default:
		serializationMode = constant.JsonSerializationMode
	}
	userFileStore = filestore.New(constant.UserStoragePath, serializationMode)

	userStorage = append(userStorage, userFileStore.Load()...)
	for {
		runCommand(*command)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}

}

func runCommand(command string) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		login()

		if authenticatedUser == nil {
			return
		}
	}
	switch command {
	case "create-task":
		createTask()
	case "list-task":
		listTask()
	case "create-category":
		createCategory()
	case "register-user":
		registerUser(userFileStore)
	case "exit":
		os.Exit(0)
	case "no-command":
		break
	default:
		fmt.Println("command is not valid", command)

	}
}
func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, duedate, category string

	fmt.Println("please enter the task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("please enter the task category-id")
	scanner.Scan()
	category = scanner.Text()
	categoryID, err := strconv.Atoi(category)
	if err != nil {
		fmt.Printf("category-id is not valid, %v\n", err)

		return
	}
	isFound := false
	for _, c := range categoryStorage {
		if c.ID == categoryID && c.UserID == authenticatedUser.ID {
			isFound = true
			break
		}
	}
	if !isFound {
		fmt.Printf("category-id is not found\n")
		return
	}

	fmt.Println("please enter the task due date")
	scanner.Scan()
	duedate = scanner.Text()
	task := models.Task{
		ID:         len(taskStorage) + 1,
		Title:      title,
		DueDate:    duedate,
		CategoryID: categoryID,
		IsDone:     false,
		UserID:     authenticatedUser.ID,
	}
	taskStorage = append(taskStorage, task)

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

	c := models.Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}
	categoryStorage = append(categoryStorage, c)

}
func registerUser(store contract.UserWriteStore) {
	scanner := bufio.NewScanner(os.Stdin)
	var email, password, name string

	fmt.Println("please enter the user name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("please enter the user email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the user password")
	scanner.Scan()
	password = scanner.Text()

	fmt.Printf("user:%v %v %v %v\n", len(userStorage)+1, email, password, name)

	hashedPassword := hashPassword(password)
	u := models.User{
		ID:       len(userStorage) + 1,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}
	userStorage = append(userStorage, u)
	store.Save(u)

}
func listTask() {
	for _, task := range taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Println(task)
		}
	}

}
func login() {
	fmt.Println("login process")
	var email, password string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("please enter the user email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("please enter the user password")
	scanner.Scan()
	password = scanner.Text()

	for _, user := range userStorage {
		if (user.Email == email) && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
			authenticatedUser = &user

			break
		}
	}
	if authenticatedUser == nil {
		fmt.Println("the email or password is incorrect")
		return
	}
	fmt.Println("you are logged in")

}

func hashPassword(password string) string {
	pass := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)

}
