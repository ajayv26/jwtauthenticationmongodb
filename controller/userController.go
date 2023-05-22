package controller

import (
	"context"
	"fmt"
	"jwtauth/database"
	"jwtauth/helpers"
	"jwtauth/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email or password not correct")
		check = false
	}
	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		validateErr := validate.Struct(user)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occure while checking error"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occur while checking phone"})
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone already exsist"})
		}
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAlltokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken

		resultNo, inseretErr := userCollection.InsertOne(ctx, user)
		if inseretErr != nil {
			msg := fmt.Sprintf("user item ws not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultNo)
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email password is incorrect"})
			return
		}
		passValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true
		c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
		return
	}

	if foundUser.Email == nil{
		 c.JSON(http.StatusInternalServerError, gin.H{"error":"user Not found"})
	}
	 token, refreshToken, _:=helper.GenerateAlltokens(*foundUser.Email,*foundUser.FirstName,*foundUser.Lastname,*foundUser.UserType,*foundUser.UserID)
    helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)
	 err =userCollection.FindOne(ctx,bson.M{"user_id":foundUser.UserID}).Decode(&foundUser)
     if err != nil{
		 c.JSON(http.Status.IntertnalServerError, gin.H{"error":err.Error()})
		return 
	 }
	 c.JSON(http.StatusOK, foundUser)
}

func GetUsers() gin.HandlerFunc{  
	 return func(c *gin.Context){
		 helper.CheckUserType(c, "ADMIN"); err != nil{
			 c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
             return
		 }
		 var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		  recordPerPage, err :=   strconv.Atoi(c.Query("recordPerPage"))
		  if err != nil || recordPerPage <1{  
			 recordPerPage = 10
		  }
         page, err1:=  strconv.Atoi(c.Query("page "))
		 if err1!= nil || page<1{
			page =1
		 }
		 startIndex := (page -1)* recordPerPage
		 startIndex, err = strconv.Atoi(c.Query("startIndex")) 

		 matchStage := bson.D{{"$match",bson.D{{}}}}
		 groupStage := bson.D{{"$group",bson.D{
			{"_id",  bson.D{{"_id", "null"}}}, 
			{"total_count", bson.D{{"$sum", 1}}},
			{"data", bson.D{{"$push", "$$ROOT"}}}
		}}}
         projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"total_conut", 1},
			{"user_items",bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
		   }}
		 }       
         userCollection.Aggregrate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage
		 })
		 defer cancel()
		 if err != nil {
			 c.JSON{http.InternalServerError, gin.H{"error":"error occured while listing user"}}
		 }
		 var allUsers []bson.M
        if err = result.All(ctx, &allusers); err != nil {
			 log.Fatal(http.StatusOk, allusers[0])
		}

		}
	 }

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"userID": userID}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
