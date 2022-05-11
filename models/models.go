package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	id             primitive.ObjectID
	firstName      *string
	lastName       *string
	password       *string
	email          *string
	authToken      *string
	refreshToken   *string
	createdAt      time.Time
	updatedAt      time.Time
	userId         *string
	userCart       []ProductUser
	addressDetails []Address
	orderStatus    []Order
}

type Product struct {
	id     primitive.ObjectID
	name   *string
	price  *uint64
	rating *uint8
	image  *string
}

type ProductUser struct {
	id     primitive.ObjectID
	name   *string
	price  *uint64
	rating *uint8
	image  *string
}

type Address struct {
	id         primitive.ObjectID
	street     *string
	house      *string
	city       *string
	postalCode *string
}

type Order struct {
	id            primitive.ObjectID
	cart          []ProductUser
	orderedAt     time.Time
	price         *uint64
	discount      *int
	paymentMethod Payment
}

type Payment struct {
	id             primitive.ObjectID
	digital        bool
	cashOnDelivery bool
}
