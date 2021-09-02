package entity

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IComponent interface {
	Update(dt float64)
}

/*
	Use entity.New() to create new entity
*/
type Entity struct {
	ID         primitive.ObjectID    `bson:"_id"`
	Components map[string]IComponent `bson:"components"`
}

func New() *Entity {
	return &Entity{
		Components: make(map[string]IComponent),
	}
}

func (e Entity) GetComponent(comp_type string) IComponent {
	component := e.Components[comp_type]
	return component
}

func (e *Entity) AddComponent(components ...IComponent) {
	for _, c := range components {
		component_type := reflect.TypeOf(c).String()
		component_type = component_type[strings.IndexByte(component_type, '.')+1:]
		e.Components[component_type] = c
	}

}

func (e *Entity) Update(dt float64) {
	for _, c := range e.Components {
		c.Update(dt)
	}
}
