package controllers

import (
	"GoEcommerceApp/models"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Parameters"})
			c.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.IndentedJSON(500, "Internal Server")
		}

		var newAddress models.Address
		newAddress.Id = primitive.NewObjectID()

		if err := c.BindJSON(&newAddress); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		matchStage := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwindStage := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$Id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pointCursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{matchStage, unwindStage, groupStage})
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var addressInfo []bson.M
		if err = pointCursor.All(ctx, &addressInfo); err != nil {
			panic(err)
		}

		var size int32 = 0
		for _, addressNo := range addressInfo {
			count := addressNo["count"]
			size += count.(int32)
		}
		if size < 2 {
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "address", Value: newAddress}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				log.Println(err)
			}
		} else {
			c.IndentedJSON(400, "Not Allowed")
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successful Creation")
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Parameters"})
			c.Abort()
			return
		}
		userObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id", Value:userObjectId}}
		update := bson.D{{Key:"$set", Value:bson.D{primitive.E{Key:"Address.0.houseName", Value: editAddress.House}, {Key:"Address.0.streetName", Value: editAddress.Street}, {Key:"Address.0.cityName", Value: editAddress.City}, {Key:"Address.0.postalCode", Value:editAddress.postalCode}}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something Went Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successful Edit")
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Parameters"})
			c.Abort()
			return
		}
		userObjectId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key:"_id", Value:userObjectId}}
		update := bson.D{{Key:"$set", Value:bson.D{primitive.E{Key:"Address.1.houseName", Value: editAddress.House}, {Key:"Address.1.streetName", Value: editAddress.Street}, {Key:"Address.1.cityName", Value: editAddress.City}, {Key:"Address.1.postalCode", Value:editAddress.postalCode}}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something Went Wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successful Edit")
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Query("id")

		if userId == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Parameters"})
			c.Abort()
			return
		}

		userAddresses := make([]models.Address, 0)
		userObjectId, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: userObjectId}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: userAddresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Invalid update command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Deleted Addresses")
	}
}
