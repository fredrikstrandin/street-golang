package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type StructCity struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Name   string             `bson:"name,omitempty"`
	Street []StructStreet     `bson:"street,omitempty"`
}

type StructStreet struct {
	Name     string           `bson:"name,omitempty"`
	StreetNo []StructStreetNo `bson:"streetNo,omitempty"`
}

type StructStreetNo struct {
	Name   string             `bson:"name,omitempty"`
	Detail StructStreetDetail `bson:"detial,omitempty"`
}

type StructStreetDetail struct {
	Lon      float64 `bson:"lon,omitempty"`
	Lat      float64 `bson:"lat,omitempty"`
	Postcode string  `bson:"postcode,omitempty"`
}
