package controllers

import (
	"GoEcommerceApp/database"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

		userQueryId := c.Query("id")
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

		userQueryId := c.Query("id")
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
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {

		userQueryId := c.Query("id")
		if userQueryId == "" {
			log.Println("user Id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user Id is empty"))
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId)

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

		userQueryId := c.Query("id")
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
