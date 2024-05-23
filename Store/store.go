package Store

import (
	"attdapp/Models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	Collection *mongo.Collection
}

const uri = "mongodb+srv://medha:drumDRO67%23%24@cluster0.qj0tdiv.mongodb.net/"

func (m *MongoStore) OpenConnectionWithMongoDB() {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		return
	}

	fmt.Println("Connected to MongoDB!")
	m.Collection = client.Database("attdapp").Collection("users")
}

func (m *MongoStore) StoreUserData(user Models.User, imageData []byte, filename string) error {
	fmt.Println("Trying to insert user data into MongoDB")

	// Create the image struct
	image := Models.Image{
		Filename: filename,
		Data:     imageData,
	}

	// Assign image to user
	user.Image = image

	_, err := m.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Println("Error inserting user data:", err)
		return err
	}
	fmt.Println("Insertion of user data successful")
	return nil
}

func (m *MongoStore) UserLogin(handle string, password string) (bool, string) {
	var foundUser Models.User
	err := m.Collection.FindOne(context.TODO(), bson.M{
		"handle":   handle,
		"password": password,
	}).Decode(&foundUser)

	if err != nil {
		fmt.Println("wrong credentials(error at userlogin function): ", err)
		return false, ""
	}
	if foundUser.Role == "admin" {
		return true, "admin"
	} else if foundUser.Role == "student" {
		return true, "student"
	}
	return false, ""
}

func (m *MongoStore) GetAllUsers() []Models.User {
	var users []Models.User
	cursor, err := m.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println("Error retrieving users:", err)
		return nil
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user Models.User
		if err = cursor.Decode(&user); err != nil {
			fmt.Println("Error decoding user:", err)
			continue
		}
		users = append(users, user)
	}
	return users
}
