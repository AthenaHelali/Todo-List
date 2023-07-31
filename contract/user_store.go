package contract

import "todo-list/models"

type UserWriteStore interface {
	Save(u models.User)
}
type UserReadStore interface {
	Load() []models.User
}
