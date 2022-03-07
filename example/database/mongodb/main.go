package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rookie-ninja/rk-boot/v2"
	"github.com/rookie-ninja/rk-db/mongodb"
	"github.com/rookie-ninja/rk-gin/v2/boot"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

var (
	userCollection *mongo.Collection
)

func createCollection(db *mongo.Database, name string) {
	opts := options.CreateCollection()
	err := db.CreateCollection(context.TODO(), name, opts)
	if err != nil {
		fmt.Println("collection exists may be, continue")
	}
}

func main() {
	boot := rkboot.NewBoot()

	boot.Bootstrap(context.TODO())

	// Auto migrate database and init global userDb variable
	db := rkmongo.GetMongoDB("my-mongo", "users")
	createCollection(db, "meta")

	userCollection = db.Collection("meta")

	// Register APIs
	ginEntry := rkgin.GetGinEntry("user-service")
	ginEntry.Router.GET("/v1/user", ListUsers)
	ginEntry.Router.GET("/v1/user/:id", GetUser)
	ginEntry.Router.PUT("/v1/user", CreateUser)
	ginEntry.Router.POST("/v1/user/:id", UpdateUser)
	ginEntry.Router.DELETE("/v1/user/:id", DeleteUser)

	boot.WaitForShutdownSig(context.TODO())
}

// *************************************
// *************** Model ***************
// *************************************

type User struct {
	Id   string `bson:"id" yaml:"id" json:"id"`
	Name string `bson:"name" yaml:"name" json:"name"`
}

func ListUsers(ctx *gin.Context) {
	userList := make([]*User, 0)

	cursor, err := userCollection.Find(context.Background(), bson.D{})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if err = cursor.All(context.TODO(), &userList); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, userList)
}

func GetUser(ctx *gin.Context) {
	res := userCollection.FindOne(context.Background(), bson.M{"id": ctx.Param("id")})

	if res.Err() != nil {
		ctx.AbortWithError(http.StatusInternalServerError, res.Err())
		return
	}

	user := &User{}
	err := res.Decode(user)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func CreateUser(ctx *gin.Context) {
	user := &User{
		Id:   xid.New().String(),
		Name: ctx.Query("name"),
	}

	_, err := userCollection.InsertOne(context.Background(), user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func UpdateUser(ctx *gin.Context) {
	uid := ctx.Param("id")

	user := &User{
		Id:   uid,
		Name: ctx.Query("name"),
	}

	res, err := userCollection.UpdateOne(context.Background(), bson.M{"id": uid}, bson.D{
		{"$set", user},
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if res.MatchedCount < 1 {
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	bytes, _ := json.Marshal(user)
	fmt.Println(string(bytes))

	ctx.JSON(http.StatusOK, user)
}

func DeleteUser(ctx *gin.Context) {
	res, err := userCollection.DeleteOne(context.Background(), bson.M{
		"id": ctx.Param("id"),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if res.DeletedCount < 1 {
		ctx.JSON(http.StatusNotFound, "user not found")
		return
	}

	ctx.String(http.StatusOK, "success")
}
