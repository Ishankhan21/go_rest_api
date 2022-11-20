package main

import (
	"context"
	"log"
	"time"

	"github.com/Ishankhan21/go/first-web-app/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	var router = gin.Default()
	var address = ":3000"

	healthChecks := router.Group("/health")           // Route groups,
	healthChecks.GET("/hello", func(c *gin.Context) { // -> /health/hello
		c.String(200, "Hello world")
	})

	healthChecks.GET("/health_check", func(c *gin.Context) { // -> /health/health_check
		c.String(200, "Server is healthy")
	})

	dbClient := db.ConnectDB()

	router.POST("/products", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		productsDatabase := dbClient.Database("test")
		productsCollection := productsDatabase.Collection("products")

		product, err := productsCollection.InsertOne(ctx, bson.D{
			{Key: "name", Value: "Mobile"},
			{Key: "brand", Value: "Apple"},
		})
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Created product =", product)
		}

	})

	router.GET("/products", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		productsDatabase := dbClient.Database("test")
		productsCollection := productsDatabase.Collection("products")

		cursor, err := productsCollection.Find(ctx, bson.M{})
		defer cursor.Close(ctx)

		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Got products")
		}

		var products []bson.M

		if err = cursor.All(ctx, &products); err != nil {
			log.Fatal(err)
		}

		for _, product := range products {
			log.Println("Product", product)
		}
	})

	router.PATCH("/products", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		productsDatabase := dbClient.Database("test")
		productsCollection := productsDatabase.Collection("products")

		id, _ := primitive.ObjectIDFromHex("6379e57f27f48c3b124dfa81")

		result, err := productsCollection.UpdateByID(ctx, bson.M{"_id": id}, bson.D{
			{"$set", bson.D{{"name", "MobileNext"}}},
		})

		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Updated document", result.ModifiedCount)
		}
	})

	router.DELETE("/products", func(c *gin.Context) {
		log.Println("Detele handler called")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		productsDatabase := dbClient.Database("test")
		productsCollection := productsDatabase.Collection("products")

		result, err := productsCollection.DeleteOne(ctx, bson.D{{"name", "Mobile"}})

		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Updated document", result.DeletedCount)
		}

	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	router.Run(address)
}
