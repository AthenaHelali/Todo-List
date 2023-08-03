package memorystore

import (
	"fmt"
	"todo-list/models"
)

type Category struct {
	categories []models.Category
}

func (c Category) DoesThisUserHasThisCategoryID(userID, categoryID int) (bool, error) {
	isFound := false
	for _, c := range c.categories {
		if c.ID == categoryID && c.UserID == userID {
			isFound = true
			break
		}
	}
	if !isFound {
		return false, fmt.Errorf("category-id is not found")
	}
	return true, nil
}
