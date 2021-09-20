package ecs

import (
	"sync"
)

type EntityID int

type ComponentID int

type Manager struct {
	lock         *sync.RWMutex
	nextEntityID int

	entities     []*Entity
	entitiesByID map[EntityID]*Entity
	components   []*Component
}

func NewManager() *Manager {
	return &Manager{
		entitiesByID: make(map[EntityID]*Entity),
	}
}

func (manager *Manager) NewEntity() *Entity {
	entity := &Entity{
		ID:      EntityID(manager.nextEntityID),
		manager: manager,
	}
	manager.nextEntityID++

	manager.lock.Lock()
	manager.entities = append(manager.entities, entity)
	manager.entitiesByID[entity.ID] = entity
	manager.lock.Unlock()

	return entity
}

type Entity struct {
	ID      EntityID
	manager *Manager
}

func (entity *Entity) GetID() EntityID {
	return entity.ID
}

type Component struct {
	id       ComponentID
	datalock *sync.RWMutex
	data     map[EntityID]interface{}
}

func (component *Component) GetID() ComponentID {
	return component.id
}
