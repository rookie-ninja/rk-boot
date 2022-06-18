package main

import (
	"context"
	"github.com/gin-gonic/gin"
	rkboot "github.com/rookie-ninja/rk-boot/v2"
	rkmongo "github.com/rookie-ninja/rk-db/mongodb"
	rkgin "github.com/rookie-ninja/rk-gin/v2/boot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var (
	todoCollection *mongo.Collection
)

// @title Swagger API With MongoDB
// @version 1.0
// @description This is Example MongoDB.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	boot := rkboot.NewBoot()
	boot.Bootstrap(context.TODO())

	db := rkmongo.GetMongoDB("todo-mongo", "tododb")
	todoCollection = db.Collection("todos")

	// Register APIs
	ginEntry := rkgin.GetGinEntry("todo-service")
	ginEntry.Router.GET("/", HelloWorld)
	ginEntry.Router.GET("/v1/todos", ListTodos)
	ginEntry.Router.GET("/v1/todo/:id", GetTodo)
	ginEntry.Router.POST("/v1/todo", CreateTodo)
	ginEntry.Router.PUT("/v1/todo/:id", UpdateTodo)
	ginEntry.Router.DELETE("/v1/todo/:id", DeleteTodo)

	boot.WaitForShutdownSig(context.TODO())
}

// *************************************
// *************** Model ***************
// *************************************

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title"`
	Body      string             `json:"body"`
	Completed bool               `json:"completed"`
	CreatedAt time.Time          `json:"created_at"`
}

// @Summary Hello world
// @Tags Hello world
// @version 1.0
// @produce application/json
// @Success 200
// @Router / [get]
func HelloWorld(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"Message": "Hello World",
	})
}

// @Summary Get List Todo
// @ID get-todo-list
// @Tags Todos
// @version 1.0
// @produce application/json
// @Success 200
// @Failure 500
// @Router /v1/todos [get]
func ListTodos(ctx *gin.Context) {
	t := []Todo{}
	cursor, err := todoCollection.Find(context.Background(), bson.M{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if err = cursor.All(context.TODO(), &t); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, t)
}

// @Summary Get Todo
// @ID get-todo-by-id
// @Param id path string true "todo ID"
// @Tags Todos
// @version 1.0
// @produce application/json
// @Success 200
// @Failure 500
// @Router /v1/todo/{id} [get]
func GetTodo(ctx *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	res := todoCollection.FindOne(context.Background(), bson.M{"_id": id})

	if res.Err() != nil {
		ctx.AbortWithError(http.StatusInternalServerError, res.Err())
		return
	}

	t := Todo{}
	err := res.Decode(&t)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, t)
}

// @Summary Post Todo
// @Tags Todos
// @version 1.0
// @produce application/json
// @Success 201
// @Failure 500
// @Router /v1/todo [post]
func CreateTodo(ctx *gin.Context) {
	t := &Todo{}
	ctx.BindJSON(&t)
	t.Completed = false
	t.CreatedAt = time.Now()

	_, err := todoCollection.InsertOne(context.Background(), t)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.String(http.StatusCreated, "created")
}

// @Summary Update Todo
// @ID update-todo-by-id
// @Param id path string true "todo ID"
// @Tags Todos
// @version 1.0
// @produce application/json
// @Success 200
// @Failure 404
// @Router /v1/todo/{id} [put]
func UpdateTodo(ctx *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	t := &Todo{}
	ctx.BindJSON(&t)
	t.CreatedAt = time.Now()

	res, err := todoCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{
		"$set": t,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if res.MatchedCount < 1 {
		ctx.JSON(http.StatusNotFound, "todo id not found")
		return
	}

	ctx.String(http.StatusOK, "success")
}

// @Summary Delete Todo
// @ID delete-todo-by-id
// @Param id path string true "todo ID"
// @Tags Todos
// @version 1.0
// @produce application/json
// @Success 200
// @Failure 404
// @Router /v1/todo/{id} [delete]
func DeleteTodo(ctx *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	res, err := todoCollection.DeleteOne(context.Background(), bson.M{
		"_id": id,
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if res.DeletedCount < 1 {
		ctx.JSON(http.StatusNotFound, "todo id not found")
		return
	}

	ctx.String(http.StatusOK, "success")
}
