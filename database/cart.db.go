package database

import (
	"GoEcommerceApp/models"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct   = errors.New("Product Cannot Be Found")
	ErrCantDecodeProduct = errors.New("Product Cannot Be Decoded")
	ErrCantUpdateUser    = errors.New("User Cannot Be Updated")
	ErrCantRemoveItem    = errors.New("Item Cannot Be Removed")
	ErrCantBuyItem       = errors.New("Item Cannot Be Bought")
	ErrInvalidUserId     = errors.New("Invalid User Id")
)

func AddToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	searchFromDb, err := prodCollection.Find(ctx, bson.M{"_id": productId})
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}
	var cart []models.ProductUser
	err = searchFromDb.All(ctx, &cart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProduct
	}

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserId
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "userCart", Value: bson.D{{Key: "$each", Value: cart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserId
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"userCart": bson.M{"_id": productId}}}

	_, err = UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItem
	}

	return nil
}

func BuyItemFromCart() string {

}

func InstantBuy() string {

}
