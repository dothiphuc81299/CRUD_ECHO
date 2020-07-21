package test

import (
	"bytes"
	"context"
	"echo/controllers"
	"echo/database"
	"echo/models"
	"encoding/json"
	"io"
	"log"
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

//SetupSuite to.....
func (s CreateModelSuite) SetupSuite() {
	database.Connectdb("todo-test")
	//removeOldData()
	//addRecord(idGet)
	//addRecord(idDelete)

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
	//e := echo.New()
	model := models.Model{
		Title: "title 1",
		Desc:  "aggd",
	}
	req, _ := http.NewRequest("POST", "/Create", ToIOReader(model))
	req.Header.Set("Content-Type", "application/json")
	//create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateModel)
	handler.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	var res models.Model
	json.Unmarshal([]byte(rec.Body.String()), &res)
	assert.Equal(s.T(), res.Title, model.Title)
	assert.Equal(s.T(), res.Desc, model.Desc)
}
func (s *CreateModelSuite) TestGetAllModel() {
	req, _ := http.NewRequest("GET", "/List", nil)
	req.Header.Set("Content-Type", "application/json")
	//create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetAllModel)
	handler.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	ctx := context.Background()
	cursor, err := database.DB.Collection("todos").Find(ctx, bson.M{})
	if err != nil {
		panic("err")
	}
	var todos []models.Model
	defer cursor.Close(ctx)
	cursor.All(ctx, &todos)
	var res []models.Model
	json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(s.T(), todos, res)

}
func (s *CreateModelSuite) TestGetModel() {
	e := echo.New()
	req, _ := http.NewRequest("GET", "/Get/{id}", nil)
	req.Header.Set("Content-Type", "application/json")
	//create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idGet.Hex())
	handler := http.HandlerFunc(controllers.GetModel)
	handler.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	// var todos models.Model
	// var res models.Model
	// json.Unmarshal(rec.Body.Bytes(), &res)
	// assert.Equal(s.T(), todos, res)
}
func (s *CreateModelSuite) TestDeleteModel() {
	//e := echo.New()
	req, _ := http.NewRequest("DELETE", "/Delete/5f16af3d0aded989447b177b", nil)
	req.Header.Set("Content-Type", "application/json")

	//create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rec := httptest.NewRecorder()
	//	c := e.NewContext(req, rec)
	// c.SetParamNames("id")
	// c.SetParamValues("5f16af3d0aded989447b177c")
	handler := http.HandlerFunc(controllers.DeleteModel)
	handler.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	x := struct {
		DeletedCount int `bson:"DeletedCount" json:"DeletedCount"`
	}{}

	json.Unmarshal([]byte(rec.Body.String()), &x)
	log.Printf(rec.Body.String())
	log.Printf("%+v", x)

	assert.Equal(s.T(), 1, x.DeletedCount)

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
