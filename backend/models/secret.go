package models

import (
	"time"
)

type Secret struct {
	ExpireDate time.Time `bson:"expire_date" xml:"expire_date" binding:"required"`
	Views      int       `bson:"views" xml:"views" binding:"required"`
}

type Data struct {
	Data   string `bson:"data" xml:"data" binding:"required"`
	Views  int    `bson:"views" xml:"views" binding:"required"`
	Expire int    `bson:"expire" xml:"expire" binding:"required"`
}

type ResponseData struct {
	Data    string `bson:"data" xml:"data" `
	Message string `bson:"message" xml:"message"`
	Object  Secret
}

type ResultToken struct {
	Message string `bson:"message" xml:"message"`
	Data    string `bson:"data" xml:"data"`
}

type ErrorModel struct {
	Message string `bson:"message" xml:"message"`
}
