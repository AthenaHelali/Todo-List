package filestore

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo-list/constant"
	"todo-list/models"
)

type FileStorage struct {
	filePath          string
	serializationMode string
}

func New(path, serializationMode string) FileStorage {
	return FileStorage{
		filePath:          path,
		serializationMode: serializationMode,
	}
}

func (f FileStorage) Save(u models.User) {
	f.writeUserToFile(u)
}
func (f FileStorage) Load() []models.User {
	var uStore []models.User

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
		var userStruct models.User
		switch f.serializationMode {
		case constant.ManualSerializationMode:
			var dErr error
			userStruct, dErr = deserializeUserFromManual(u)
			if dErr != nil {
				fmt.Println("can't deserialize user record to user struct ")
			}

		case constant.JsonSerializationMode:
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
func (f FileStorage) writeUserToFile(user models.User) {
	var file *os.File

	file, err := os.OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("file does not exist", err)

		return
	}
	defer file.Close()
	var data []byte
	if f.serializationMode == constant.ManualSerializationMode {
		data = []byte(fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n", user.ID, user.Name,
			user.Email, user.Password))
	} else if f.serializationMode == constant.JsonSerializationMode {

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
func deserializeUserFromManual(userStr string) (models.User, error) {
	userFields := strings.Split(userStr, ",")
	var user models.User

	for _, field := range userFields {
		values := strings.Split(field, ":")
		fieldName := strings.ReplaceAll(values[0], " ", "")
		fieldValue := strings.ReplaceAll(values[1], " ", "")
		switch fieldName {
		case "id":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {
				return models.User{}, fmt.Errorf("strconv error")
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
