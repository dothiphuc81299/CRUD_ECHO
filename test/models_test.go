package test

import (
	"bytes"
	"context"
	"echo/controllers"
	"echo/database"
	"echo/models"
	"encoding/json"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//CreateModelSuite to ....
type CreateModelSuite struct {
	suite.Suite
	Models []models.Model
}

var idGet = primitive.NewObjectID()
var idDelete = primitive.NewObjectID()
var idUpdate = primitive.NewObjectID()
var idCompleted = primitive.NewObjectID()

//SetupSuite to.....
func (s CreateModelSuite) SetupSuite() {
	database.Connectdb("todo-test")
	removeOldData()
	addRecord(idGet)
	addRecord(idDelete)
	addRecord(idUpdate)
	addRecord(idCompleted)

}

//TearDownSuite to clear data
func (s CreateModelSuite) TearDownSuite() {
	//removeOldData()
}

//removeOldData to ...
func removeOldData() {
	database.DB.Collection("todos").DeleteMany(context.Background(), bson.M{})

}

//TestCreateModel to ...
func (s *CreateModelSuite) TestCreateModel() {
	e := echo.New()
	model := models.Model{
		Title: "title 1",
		Desc:  "aggd",
	}
	req, _ := http.NewRequest(http.MethodPost, "/todos", ToIOReader(model))
	req.Header.Set("Content-Type", "application/json")
	//create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	controllers.CreateModel(c)
	assert.Equal(s.T(), http.StatusCreated, rec.Code)
	var res models.Model
	json.Unmarshal([]byte(rec.Body.String()), &res)
	assert.Equal(s.T(), res.Title, model.Title)
	assert.Equal(s.T(), res.Desc, model.Desc)
}

func (s *CreateModelSuite) TestGetAllModel() {
	e := echo.New()
	req, _ := http.NewRequest(http.MethodGet, "/todos", nil)
	//req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	controllers.GetAllModel(c)
	cursor, err := database.DB.Collection("todos").Find(context.TODO(), bson.M{})
	if err != nil {
		panic("err")
	}
	var list []models.Model
	cursor.All(context.TODO(), &list)
	var res []models.Model
	json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(s.T(), list, res)

}
func (s *CreateModelSuite) TestDeleteModel() {
	e := echo.New()
	req, _ := http.NewRequest(http.MethodDelete, "/todos/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idDelete.Hex())
	controllers.DeleteModel(c)
	assert.Equal(s.T(), http.StatusOK, rec.Code)

	x := struct {
		DeletedCount int `bson:"DeletedCount" json:"DeletedCount"`
	}{}
	json.Unmarshal([]byte(rec.Body.String()), &x)
	assert.Equal(s.T(), x.DeletedCount, 1)

}
func (s *CreateModelSuite) TestUpdateModel() {
	e := echo.New()
	req, _ := http.NewRequest(http.MethodPut, "/todos/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idUpdate.Hex())
	controllers.UpdateModel(c)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	x := struct {
		MatchedCount  int    `bson:"MatchedCount" json:"MatchedCount"`
		ModifiedCount int    `bson:"ModifiedCount" json:"ModifiedCount"`
		UpsertedCount int    `bson:"UpsertedCount" json:"UpsertedCount"`
		UpsertedID    string `bson:"UpsertedID" json:"UpsertedID"`
	}{}
	json.Unmarshal([]byte(rec.Body.String()), &x)
	assert.Equal(s.T(), x.MatchedCount, 1)
	assert.Equal(s.T(), x.ModifiedCount, 1)
	assert.Equal(s.T(), x.UpsertedCount, 0)
	assert.Equal(s.T(), x.UpsertedID, "")

}

func (s *CreateModelSuite) TestCompletedModel() {
	e := echo.New()
	req, _ := http.NewRequest(http.MethodPut, "/todos/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idCompleted.Hex())
	controllers.CompletedModel(c)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	x := struct {
		MatchedCount  int    `bson:"MatchedCount" json:"MatchedCount"`
		ModifiedCount int    `bson:"ModifiedCount" json:"ModifiedCount"`
		UpsertedCount int    `bson:"UpsertedCount" json:"UpsertedCount"`
		UpsertedID    string `bson:"UpsertedID" json:"UpsertedID"`
	}{}
	json.Unmarshal([]byte(rec.Body.String()), &x)
	assert.Equal(s.T(), x.MatchedCount, 1)
	assert.Equal(s.T(), x.ModifiedCount, 1)
	assert.Equal(s.T(), x.UpsertedCount, 0)
	assert.Equal(s.T(), x.UpsertedID, "")

}


func (s *CreateModelSuite) TestGetModelByID() {
	e := echo.New()
	req, _ := http.NewRequest(http.MethodPut, "/todos/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idGet.Hex())
	controllers.GetModelByID(c)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	// err := database.DB.Collection("todos").FindOne(context.TODO(), models.Model{ID: idGet})
	// if err != nil {
	// 	panic("err")
	// }
	// x := struct {
	// 	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	// 	Title string             `json:"title" bson:"title"`
	// 	Desc  string             `json:"desc" bson:"desc"`
	// }{}

	// json.Unmarshal(rec.Body.Bytes(), &x)
	// fmt.Println(x)
	// assert.Equal(s.T(), idGet, x.ID)
}

//ToIOReader to...
func ToIOReader(i interface{}) io.Reader {
	b, _ := json.Marshal(i)
	return bytes.NewReader(b)
}
func addRecord(id primitive.ObjectID) {
	todo := models.Model{
		ID:        id,
		Title:     "title 2",
		Desc:      "desc 2",
		Completed: false,
	}
	database.DB.Collection("todos").InsertOne(context.TODO(), todo)
}

//TestCreateModelSuite to..
func TestCreateModelSuite(t *testing.T) {
	suite.Run(t, new(CreateModelSuite))
}
