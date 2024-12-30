package user

import (
	"github.com/google/uuid"
)

type model struct {
	id   uuid.UUID
	name string
	age  int
}

type jsonBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func NewModel(name string, age int) (*model, error) {
	return &model{id: uuid.Must(uuid.NewRandom()), name: name, age: age}, nil
}

func (m *model) jsonBody() *jsonBody {
	return &jsonBody{Name: m.name, Age: m.age}
}

func (m *model) Rename(name string) {
	m.name = name
}
