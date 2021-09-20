package entity

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IComponent interface {
	SetParent(e *Entity)
}

/*
	Use entity.New() to create new entity
*/
type Entity struct {
	MongoID    primitive.ObjectID `bson:"_id"`
	ID         int
	Components map[string]IComponent `bson:"components"`

	manager *Manager
}

//TODO Remove this function because manager creates entities
func New() *Entity {
	return &Entity{
		ID:         1,
		Components: make(map[string]IComponent),
	}
}

func (entity Entity) GetID() int {
	return entity.ID
}

func (entity *Entity) SetManager(m *Manager) {
	entity.manager = m
	entity.manager.entities = append(entity.manager.entities, entity)
	entity.manager.entitiesByID[EntityID(entity.ID)] = entity
}

func (entity *Entity) AddComponent(componentdata IComponent) *Entity {
	component_type := reflect.TypeOf(componentdata).String()
	//Remove package name and * from the name
	component_type = component_type[strings.IndexByte(component_type, '.')+1:]

	if entity.manager.components[ComponentType(component_type)] == nil {
		entity.manager.components[ComponentType(component_type)] = &Component{}
		entity.manager.components[ComponentType(component_type)].data = make(map[EntityID]IComponent)
	}
	entity.manager.components[ComponentType(component_type)].data[EntityID(entity.ID)] = componentdata
	componentdata.SetParent(entity)

	return entity
}

func (entity *Entity) RemoveComponent(comp_type ComponentType) *Entity {
	delete(entity.manager.components[comp_type].data, EntityID(entity.ID))

	return entity
}

func (entity Entity) HasComponent(comp_type ComponentType) bool {
	if _, ok := entity.manager.components[comp_type]; ok {
		if _, ok := entity.manager.components[comp_type].data[EntityID(entity.ID)]; ok {
			return true
		}
	}
	return false
}

func (entity Entity) GetComponent(comp_type ComponentType) interface{} {
	data := entity.manager.components[comp_type].data[EntityID(entity.ID)]

	return data
}
