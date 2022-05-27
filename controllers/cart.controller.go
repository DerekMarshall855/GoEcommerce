package controllers

import (
	"GoEcommerceApp/database"
	"GoEcommerceApp/models"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
}

func NewApplication(productCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		productCollection: productCollection,
		userCollection:    userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product Id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user Id is empty"))
			return
		}

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("user Id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product Id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		err = database.AddToCart(ctx, app.productCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Successfully added item to cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product Id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user Id is empty"))
			return
		}

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("user Id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product Id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		err = database.RemoveFromCart(ctx, app.productCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Successfully removed item from cart")
	}
}

func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Id"})
			c.Abort()
			return
		}

		userHexId, _ := primitive.ObjectIDFromHex(userQueryId)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userHexId}}).Decode(&filledCart)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Not Found")
			return
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userHexId}}}}
		unwindStage := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$UserCart"}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$UserCart.price"}}}}}}
		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindStage, groupStage})

		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = pointCursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledCart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("user Id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user Id is empty"))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		err := database.BuyItemFromCart(ctx, app.productCollection, app.userCollection, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product Id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user Id is empty"))
			return
		}

		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("user Id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product Id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		err = database.InstantBuy(ctx, app.productCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Successfully bought item")
	}
}
