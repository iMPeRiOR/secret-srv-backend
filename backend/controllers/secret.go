package controllers

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"backend/database"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDataBasedOnResponseType(context *gin.Context, input *models.Data) error {
	contentType := context.ContentType()
	switch contentType {
	case "text/xml":
		if err := context.BindXML(&input); err != nil {
			return err
		}
	default:
		if err := context.BindJSON(&input); err != nil {
			return err
		}
	}
	return nil
}

func ErrorResponseStatus(context *gin.Context, status int, message string) {
	cacceptedHeader := context.Request.Header["Accept"][0]
	switch cacceptedHeader {
	case "text/xml":
		context.XML(status, gin.H{"message": message})
	default:
		context.JSON(status, gin.H{"message": message})
	}
}

func SuccessResponseStatus(context *gin.Context, status int, data string) {
	cacceptedHeader := context.Request.Header["Accept"][0]
	switch cacceptedHeader {
	case "text/xml":
		context.XML(status, gin.H{"message": "sucess", "data": data})
	default:
		context.JSON(status, gin.H{"message": "success", "data": data})
	}
}

func SuccessResponseStatusForToken(context *gin.Context, status int, data interface{}, object models.Secret) {
	cacceptedHeader := context.Request.Header["Accept"][0]
	switch cacceptedHeader {
	case "text/xml":
		context.XML(status, gin.H{"data": data, "object": object, "message": "success"})
	default:
		context.JSON(status, gin.H{"data": data, "object": object, "message": "success"})
	}
}

// Generate new secret
// @Summary      Create a new secret
// @Description  This route generates a new secret have the user's data
// @Tags         token
// @Accept       json
// @Accept		 xml
// @Produce      json
// @Produce      xml
// @Param        body  body      models.Data  true  "Generate a secret"
// @Success      200  {object}  models.ResultToken
// @Failure      400  {object}	models.ErrorModel
// @Failure      500  {object}	models.ErrorModel
// @Router       /generate [post]
func GenerateToken(context *gin.Context) {
	var input models.Data
	if GetDataBasedOnResponseType(context, &input) != nil {
		ErrorResponseStatus(context, http.StatusBadRequest, "wrong data type or some required fields is missing")
		return
	}

	var mySigningKey = []byte(os.Getenv("SECRET_TOKEN"))
	expirationTime := time.Now().Add(time.Minute * time.Duration(input.Expire))

	// add object to database
	result, err := database.SecretsCollection.InsertOne(database.Ctx, bson.D{
		{Key: "expire_date", Value: expirationTime},
		{Key: "views", Value: input.Views},
	})
	if err != nil {
		ErrorResponseStatus(context, http.StatusInternalServerError, "database error")
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = result.InsertedID
	claims["data"] = input.Data
	claims["exp"] = expirationTime.Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		ErrorResponseStatus(context, http.StatusInternalServerError, "cannot generate a secret")
		return
	}

	SuccessResponseStatus(context, http.StatusOK, tokenString)
}

// Get the secret Information
// @Summary      analyze the secret
// @Description  This routes generate new secret have the user's data
// @Tags         token
// @Accept       json
// @Accept       xml
// @Produce      json
// @Produce      xml
// @Param        token  path      string  true  "get the secret info"
// @Success      200  {object}    models.ResponseData
// @Failure      400  {object}	  models.ErrorModel
// @Failure      500  {object}	  models.ErrorModel
// @Router       /get/{token} [post]
func GetToken(context *gin.Context) {
	tokenString := context.Param("token")

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_TOKEN")), nil
	})

	if err != nil {
		ErrorResponseStatus(context, http.StatusBadRequest, "invalid token")
		return
	}

	info := token.Claims.(*jwt.MapClaims)
	id := (*info)["id"]
	secretData := (*info)["data"]
	objectId, err := primitive.ObjectIDFromHex(fmt.Sprintf("%v", id))
	if err != nil {
		ErrorResponseStatus(context, http.StatusBadRequest, "invalid id")
		return
	}
	// get the token object from database
	var object models.Secret
	err = database.SecretsCollection.
		FindOne(database.Ctx, bson.D{{Key: "_id", Value: objectId}}).
		Decode(&object)

	if err != nil {
		ErrorResponseStatus(context, http.StatusBadRequest, "invalid token")
		return
	}

	if object.Views <= 0 {
		// delete expired object
		_, err := database.SecretsCollection.DeleteOne(database.Ctx, bson.D{{Key: "_id", Value: objectId}})
		if err != nil {
			ErrorResponseStatus(context, http.StatusInternalServerError, "database error")
			return
		}
		ErrorResponseStatus(context, http.StatusBadRequest, "No views available")
		return
	}
	// update the number of views
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "views", Value: object.Views - 1}}}}
	_, err = database.SecretsCollection.UpdateOne(database.Ctx, filter, update)

	if err != nil {
		ErrorResponseStatus(context, http.StatusInternalServerError, "database error")
		return
	}
	SuccessResponseStatusForToken(context, http.StatusOK, secretData, object)
}
