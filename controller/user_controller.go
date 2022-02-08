package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/programmerug/fibermongo/config"
	"github.com/programmerug/fibermongo/model"
	"github.com/programmerug/fibermongo/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "user")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user model.User
	defer cancel()

	// validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	}

	newUser := model.User{
		ID:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(response.UserResponse{Status: http.StatusCreated, Message: "success", Data: result})
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user model.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(response.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"user": user}})
}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user model.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: err.Error()})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(response.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: validationErr.Error()})
	}

	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	//get updated user details
	var updatedUser model.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}
	}

	return c.Status(http.StatusOK).JSON(response.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"user": updatedUser}})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			response.UserResponse{Status: http.StatusNotFound, Message: "error", Data: "User with specified ID not found!"},
		)
	}

	return c.Status(http.StatusOK).JSON(
		response.UserResponse{Status: http.StatusOK, Message: "success", Data: "User successfully deleted!"},
	)
}

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []model.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser model.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: err.Error()})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		response.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"users": users}},
	)
}
