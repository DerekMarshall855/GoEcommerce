package database

import (
	"GoEcommerceApp/models"
	"context"
	"errors"
	"log"
	"time"

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

	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItem
	}

	return nil
}

func BuyItemFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserId
	}

	var getCartItems models.User
	var orderCart models.Order

	orderCart.Id = primitive.NewObjectID()
	orderCart.OrderedAt = time.Now()
	orderCart.Cart = make([]models.ProductUser, 0)
	orderCart.PaymentMethod.CashOnDelivery = true

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$userCart"}}}}
	group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$userCart.price"}}}}}}

	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, group})
	ctx.Done()
	if err != nil {
		panic(err)
		return ErrCantBuyItem
	}

	var getUserCart []bson.M
	if err = currentResults.All(ctx, &getUserCart); err != nil {
		panic(err)
		return ErrCantBuyItem
	}

	var totalPrice int32

	for _, userItem := range getUserCart {
		price := userItem["total"]
		totalPrice = price.(int32)
	}
	orderCart.Price = int(totalPrice)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantBuyItem
	}
	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getCartItems)
	if err != nil {
		return ErrCantBuyItem
	}

	orderUpdate := bson.M{"$push": bson.M{"orders.$[].orderList": bson.M{"$each": getCartItems.UserCart}}}
	_, err = userCollection.UpdateOne(ctx, filter, orderUpdate)
	if err != nil {
		return ErrCantBuyItem
	}

	emptyUserCart := make([]models.ProductUser, 0)
	emptyCartUpdate := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "userCart", Value: emptyUserCart}}}}
	_, err = userCollection.UpdateOne(ctx, filter, emptyCartUpdate)

	return nil
}

func InstantBuy(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Println(err)
		return ErrInvalidUserId
	}

	var productDetails models.ProductUser
	var orderDetails models.Order

	orderDetails.Id = primitive.NewObjectID()
	orderDetails.OrderedAt = time.Now()
	orderDetails.Cart = make([]models.ProductUser, 0)
	orderDetails.PaymentMethod.CashOnDelivery = true

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productId}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
	}
	orderDetails.Price = int(*productDetails.Price)

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	orderDetailsUpdate := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderDetails}}}}
	_, err = userCollection.UpdateOne(ctx, filter, orderDetailsUpdate)
	if err != nil {
		log.Println(err)
	}

	orderListUpdate := bson.M{"$push": bson.M{"orders.$[].orderList": productDetails}}
	_, err = userCollection.UpdateOne(ctx, filter, orderListUpdate)
	if err != nil {
		log.Println(err)
	}

	return nil
}
