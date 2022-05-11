package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName      *string            `json:"firstName" validate:"required,min=2,max=30"`
	LastName       *string            `json:"lastName"  validate:"required,min=2,max=30"`
	Password       *string            `json:"password"  validate:"required,min=6"`
	Email          *string            `json:"email"     validate:"email, required"`
	AuthToken      *string            `json:"authToken"`
	RefreshToken   *string            `json:"refreshToken"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
	UserId         *string            `json:"userId"`
	UserCart       []ProductUser      `json:"userCart" bson:"userCart"`
	AddressDetails []Address          `json:"addressDetails" bson:"addressDetails"`
	OrderStatus    []Order            `json:"orderStatus" bson:"orderStatus"`
}

type Product struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id"`
	Name   *string            `json:"name"`
	Price  *uint64            `json:"price"`
	Rating *uint8             `json:"rating"`
	Image  *string            `json:"image"`
}

type ProductUser struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id"`
	Name   *string            `json:"name"`
	Price  *uint64            `json:"price"`
	Rating *uint8             `json:"rating"`
	Image  *string            `json:"image"`
}

type Address struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	Street     *string            `json:"street"`
	House      *string            `json:"house"`
	City       *string            `json:"city"`
	PostalCode *string            `json:"postalCode"`
}

type Order struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	Cart          []ProductUser      `json:"cart" bson:"cart"`
	OrderedAt     time.Time          `json:"orderedAt"`
	Price         *uint64            `json:"price"`
	Discount      *int               `json:"discount"`
	PaymentMethod Payment            `json:"paymentMethod" bson:"paymentMethod"`
}

type Payment struct {
	Id             primitive.ObjectID `json:"_id" bson:"_id"`
	Digital        bool               `json:"digital"`
	CashOnDelivery bool               `json:"cashOnDelivery"`
}
