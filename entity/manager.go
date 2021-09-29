package entity

type EntityID int

type ComponentID int

type ComponentType string

type Component struct {
	_type ComponentType
	data  map[EntityID]IComponent
}

type Manager struct {
	entities     []*Entity
	entitiesByID map[EntityID]*Entity
	components   map[ComponentType]*Component

	nextId int
}

func NewManager() *Manager {
	mgr := &Manager{
		entitiesByID: make(map[EntityID]*Entity),
		components:   make(map[ComponentType]*Component),
		nextId:       1,
	}
	return mgr
}

func (manager *Manager) NewEntity() *Entity {
	entity := &Entity{
		ID:      manager.nextId,
		Manager: manager,
	}
	manager.nextId++

	manager.entities = append(manager.entities, entity)
	manager.entitiesByID[EntityID(entity.ID)] = entity

	return entity
}

// func (manager *Manager) GetEntityByTelegramID(telegramId int) *Entity {
// 	players := manager.QueryEntities("PlayerComponent")
// 	for _, p := range players {
// 		player_C := p.GetComponent("PlayerComponent").(*components.PlayerComponent)
// 		if player_C.TelegramID == telegramId {
// 			return p
// 		}
// 	}
// 	return nil
// }

func (manager *Manager) GetEntityByID(id int) *Entity {
	return manager.entitiesByID[EntityID(id)]
}

func (manager *Manager) QueryEntities(comp_types ...ComponentType) []*Entity {
	result := []*Entity{}

	for _, entity := range manager.entities {
		hasAllComponents := true
		for _, comp_type := range comp_types {
			if !entity.HasComponent(comp_type) {
				hasAllComponents = false
				break
			}
		}
		if hasAllComponents {
			result = append(result, entity)
		}
	}

	return result
}

func (manager *Manager) QueryComponents(comp_type ComponentType) []IComponent {
	result := []IComponent{}

	for _, comp := range manager.components[comp_type].data {
		result = append(result, comp)
	}

	return result
}

func (manager *Manager) DisposeEntities(entities ...*Entity) {
	for _, entity := range entities {
		manager.DisposeEntity(entity)
	}
}

func (manager *Manager) DisposeEntity(entity *Entity) {
	for _, component := range manager.components {
		if entity.HasComponent(component._type) {
			entity.RemoveComponent(component._type)
		}
	}
	delete(manager.entitiesByID, EntityID(entity.ID))
	for i, e := range manager.entities {
		if e.ID == entity.ID {
			manager.entities = append(manager.entities[:i], manager.entities[i+1:]...)
			break
		}
	}
}
