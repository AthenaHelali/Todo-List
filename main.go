package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type Task struct {
	ID         int
	Title      string
	DueDate    string
	CategoryID int
	IsDone     bool
	UserID     int
}
type Category struct {
	ID     int
	Title  string
	Color  string
	UserID int
}

var (
	authenticatedUser *User
	serializationMode string
	userStorage       []User
	taskStorage       []Task
	categoryStorage   []Category
)
var userFileStore = fileStorage{
	filePath: "user.txt",
}

const (
	userStoragePath         = "user.txt"
	manualSerializationMode = "manualSerializationMode"
	jsonSerializationMode   = "jsonSerializationMode"
)

func main() {
	fmt.Println("Hello to TODO app")

	serializeModeInput := flag.String("serialize-mode", manualSerializationMode, "serialization mode to write data to file")
	command := flag.String("command", "no-command", "command to run")
	flag.Parse()

	switch *serializeModeInput {
	case manualSerializationMode:
		serializationMode = manualSerializationMode
	default:
		serializationMode = jsonSerializationMode
	}

	LoadUserFromStorage(userFileStore)

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
	task := Task{
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

	c := Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}
	categoryStorage = append(categoryStorage, c)

}
func registerUser(store userWriteStore) {
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
	u := User{
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

type userWriteStore interface {
	Save(u User)
}
type userReadeStore interface {
	Load(serializationMode string) []User
}

func (f fileStorage) writeUserToFile(user User) {
	var file *os.File

	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("file does not exist", err)

		return
	}
	defer file.Close()
	var data []byte
	if serializationMode == manualSerializationMode {
		data = []byte(fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n", user.ID, user.Name,
			user.Email, user.Password))
	} else if serializationMode == jsonSerializationMode {

		var jErr error
		data, jErr = json.Marshal(user)
		if jErr != nil {
			fmt.Println("cant marshal user struct to json", err)

			return
		}
		data = append(data, []byte("\n")...)
	} else {
		fmt.Println("invalid serialization mode")
		return
	}
	file.Write(data)

}
func LoadUserFromStorage(store userReadeStore) {
	users := store.Load(serializationMode)
	userStorage = append(userStorage, users...)

}
func deserializeUserFromManual(userStr string) (User, error) {
	userFields := strings.Split(userStr, ",")
	var user User

	for _, field := range userFields {
		values := strings.Split(field, ":")
		fieldName := strings.ReplaceAll(values[0], " ", "")
		fieldValue := strings.ReplaceAll(values[1], " ", "")
		switch fieldName {
		case "id":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {
				return User{}, fmt.Errorf("strconv error")
			}
			user.ID = id
		case "name":
			user.Name = fieldValue
		case "email":
			user.Email = fieldValue
		case "password":
			user.Password = fieldValue
		}
	}
	return user, nil

}

func hashPassword(password string) string {
	pass := []byte(password)

	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)

}

type fileStorage struct {
	filePath string
}

func (f fileStorage) Save(u User) {
	f.writeUserToFile(u)
}
func (f fileStorage) Load(serializationMode string) []User {
	var uStore []User

	file, err := os.Open(f.filePath)
	if err != nil {
		fmt.Println("can't open the file", err)
		return nil
	}
	defer file.Close()

	var data = make([]byte, 10240)
	_, oErr := file.Read(data)
	if oErr != nil {
		fmt.Println("can't read from the file", oErr)
	}

	dataStr := string(data)

	userSlice := strings.Split(dataStr, "\n")

	for index, u := range userSlice {
		if index == len(userSlice)-1 {
			continue
		}
		var userStruct User
		switch serializationMode {
		case manualSerializationMode:
			var dErr error
			userStruct, dErr = deserializeUserFromManual(u)
			if dErr != nil {
				fmt.Println("can't deserialize user record to user struct ")
			}

		case jsonSerializationMode:
			uErr := json.Unmarshal([]byte(u), &userStruct)
			if uErr != nil {
				fmt.Println("can't deserialize user record to user struct from json", uErr)
				return nil
			}
		default:
			fmt.Println("invalid serialization mode")

		}
		uStore = append(uStore, userStruct)
	}
	return uStore

}
