package entity

type Entity interface {
	GetComponent(comp_type string) interface{}
}
