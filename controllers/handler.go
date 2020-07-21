package controllers

import (
	"context"
	"echo/database"
	"echo/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	//	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Handler struct {
// 	dt map[string]*models.Model
// }

//CreateModel to insert model into data
func CreateModel(response http.ResponseWriter, request *http.Request) {
	//setting the header "Content-type" to application/json
	response.Header().Set("content-type", "application/json")
	collection := database.DB.Collection("todos")
	var model models.Model
	model.ID = primitive.NewObjectID()
	//decode the data that is sent with the request and insert it into the post
	_ = json.NewDecoder(request.Body).Decode(&model)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("ehrrr", err)
	}
	log.Println(result)

	//use the encoding package to encode all the posts as well as returning it at the same line
	//json.NewEncoder(response).Encode(result)
	json.NewEncoder(response).Encode(model)

}

//GetAllModel to get list
func GetAllModel(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var model1 []models.Model
	collection := database.DB.Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {

		var model models.Model
		cursor.Decode(&model)
		model1 = append(model1, model)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(model1)
}

//GetModel to get model by id
func GetModel(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var model models.Model
	collection := database.DB.Collection("todos")
	//getting a specific model we need ID from url
	params := mux.Vars(request)
	//extract params
	id, _ := primitive.ObjectIDFromHex(params["id"])
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, models.Model{ID: id}).Decode(&model)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(model)
}

//DeleteModel to delete model by id
func DeleteModel(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	collection := database.DB.Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	log.Println(id)
	defer cancel()
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Fatal("DeleteOne() ERROR: ", err)
	} else {
		if result.DeletedCount == 0 {
			fmt.Println("DeleteOne() document not found: ", result)
		} else {
			fmt.Println("DeleteOne result: ", result)
			fmt.Println("Delete One", reflect.TypeOf(result))
		}
	}
	json.NewEncoder(response).Encode(result)
}

//UpdateModel to update model by id
func UpdateModel(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := database.DB.Collection("todos")
	var model models.Model
	_ = json.NewDecoder(request.Body).Decode(&model)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"title":     model.Title,
		"desc":      model.Desc,
		"completed": model.Completed,
	}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("UpdateOne()", result)
	}
	json.NewEncoder(response).Encode(result)
}
