package app

import "github.com/Aktollkynn/GoProject.git/app/models"

type Model struct {
	Model interface{}
}

func RegisterModels() []Model {
	return []Model{
		{Model: models.User{}},
		{Model: models.Address{}},
	}
}
