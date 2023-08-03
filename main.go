package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"todo-list/constant"
	"todo-list/contract"
	"todo-list/models"
	"todo-list/repository/filestore"
	"todo-list/repository/memorystore"
	"todo-list/service/task"
)

var (
	authenticatedUser *models.User
	serializationMode string
	userStorage       []models.User
	categoryStorage   []models.Category
)

var userFileStore filestore.FileStorage

func main() {
	taskMemoryRepo := memorystore.NewTaskStore()
	taskService := task.NewService(taskMemoryRepo)
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
		runCommand(*command, taskService)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("please enter another command")
		scanner.Scan()
		*command = scanner.Text()
	}

}

func runCommand(command string, taskService *task.Service) {
	if command != "register-user" && command != "exit" && authenticatedUser == nil {
		login()

		if authenticatedUser == nil {
			return
		}
	}
	switch command {
	case "create-task":
		createTask(taskService)
	case "list-task":
		listTask(taskService)
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
func createTask(taskService *task.Service) {
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
	fmt.Println("please enter the task due date")
	scanner.Scan()
	duedate = scanner.Text()

	response, cErr := taskService.CreateTask(task.CreateTaskRequest{
		Title:               title,
		DueDate:             duedate,
		CategoryID:          categoryID,
		AuthenticatedUserID: authenticatedUser.ID,
	})
	if cErr != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Println("created task:", response)
	return

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
func listTask(taskService *task.Service) {
	userTasks, err := taskService.ListTask(task.ListRequest{authenticatedUser.ID})
	if err != nil {
		fmt.Println("error", err)

		return
	}
	fmt.Println("user tasks", userTasks)

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
