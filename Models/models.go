package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Role       string    `json:"role"`
	Handle     string    `json:"handle"`
	Name       string    `json:"name"`
	Password   string    `json:"pass"`
	Image      Image     `json:"image"`
	FaceVector []float32 `json:"face_vector"`
}

type Image struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Filename string             `bson:"filename"`
	Data     []byte             `bson:"data"`
}

type RecognitionRequest struct {
	Image string `json:"image"`
}
