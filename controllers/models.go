package controllers

import (
	"context"
	"echo/database"
	"echo/models"

	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreateModel to ....
func CreateModel(c echo.Context) error {
	collection := database.DB.Collection("todos")
	model := new(models.Model)
	//c.Bind(model)
	model.ID = primitive.NewObjectID()
	result, _ := collection.InsertOne(context.TODO(), model)
	if err := c.Bind(model); err != nil {
		fmt.Println("ehrrr", err)
	}
	log.Println(result)
	return c.JSON(http.StatusCreated, model)
}

//GetAllModel to ...
func GetAllModel(c echo.Context) error {
	collection := database.DB.Collection("todos")
	var model1 []models.Model
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	for cursor.Next(context.TODO()) {
		var model models.Model
		err := cursor.Decode(&model)
		if err != nil {
			fmt.Println(err)
		}
		model1 = append(model1, model)
	}
	defer cursor.Close(context.TODO())
	return c.JSON(http.StatusOK, model1)

}

//DeleteModel to ...
func DeleteModel(c echo.Context) error {
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	collection := database.DB.Collection("todos")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		log.Fatal("DeleteOne() ERROR: ", err)

	}
	return c.JSON(http.StatusOK, result)
}

//UpdateModel  to ....
func UpdateModel(c echo.Context) error {

	collection := database.DB.Collection("todos")
	model := new(models.Model)
	if err := c.Bind(model); err != nil {
		log.Println(err)
	}
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"title":     model.Title,
		"desc":      model.Desc,
		"completed": model.Completed,
	}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	// } else {
	// 	fmt.Println("UpdateOne()", result)
	// }
	return c.JSON(http.StatusOK, result)
}

//CompletedModel to ...
func CompletedModel(c echo.Context) error {
	collection := database.DB.Collection("todos")
	model := new(models.Model)
	if err := c.Bind(model); err != nil {
		log.Println(err)
	}
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"completed": true,
	}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	// } else {
	// 	fmt.Println("UpdateOne()", result)
	// }
	return c.JSON(http.StatusOK, result)

}

//GetModelByID to...
func GetModelByID(c echo.Context) error {
	collection := database.DB.Collection("todos")
	var model models.Model
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(context.TODO(), models.Model{ID: objectID}).Decode(&model)
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, model)

}
