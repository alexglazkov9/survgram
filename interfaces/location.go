package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type ILocation interface {
	GetID() primitive.ObjectID
}
