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
	MongoID primitive.ObjectID `bson:"_id"`
	ID      int

	Manager *Manager
}

func (entity Entity) GetID() int {
	return entity.ID
}

func (entity *Entity) SetManager(m *Manager) {
	entity.Manager = m
	entity.Manager.entities = append(entity.Manager.entities, entity)
	entity.Manager.entitiesByID[EntityID(entity.ID)] = entity
}

func (entity *Entity) AddComponent(componentdata IComponent) *Entity {
	component_type := reflect.TypeOf(componentdata).String()
	//Remove package name and * from the name
	component_type = component_type[strings.IndexByte(component_type, '.')+1:]

	if entity.Manager.components[ComponentType(component_type)] == nil {
		entity.Manager.components[ComponentType(component_type)] = &Component{}
		entity.Manager.components[ComponentType(component_type)].data = make(map[EntityID]IComponent)
	}
	entity.Manager.components[ComponentType(component_type)].data[EntityID(entity.ID)] = componentdata
	componentdata.SetParent(entity)

	return entity
}

func (entity *Entity) RemoveComponent(comp_type ComponentType) *Entity {
	delete(entity.Manager.components[comp_type].data, EntityID(entity.ID))

	return entity
}

func (entity Entity) HasComponent(comp_type ComponentType) bool {
	if _, ok := entity.Manager.components[comp_type]; ok {
		if _, ok := entity.Manager.components[comp_type].data[EntityID(entity.ID)]; ok {
			return true
		}
	}
	return false
}

func (entity Entity) GetComponent(comp_type ComponentType) interface{} {
	data := entity.Manager.components[comp_type].data[EntityID(entity.ID)]

	return data
}
