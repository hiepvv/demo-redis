package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiepvv/auth-no-redis/models"
)

type user struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}
type login struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,gte=6"`
}

// AddNewUserEndPoint allow add new user data into User Collection
func AddNewUserEndPoint(context *gin.Context) {
	var newUser user
	if e := context.Bind(&newUser); e != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Err"})
		panic(e)
	}
	newUserData, errorAdding := models.AddNewUser(newUser.Email)
	if errorAdding != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		panic(errorAdding)
	}
	context.JSON(http.StatusOK, newUserData)
}

// LoginNewSession ...
func LoginNewSession(context *gin.Context) {
	var loginInfo login
	if context.Bind(&loginInfo) != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Err"})
		panic("Err")
	}
	userData, errorVerifying := models.VerifyUser(loginInfo.Email, loginInfo.Password)
	if errorVerifying != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Not found"})
		panic("Not found")
	}

	tokenString := models.GenerateHMACCrypto(userData.Salt, userData.Email)
	context.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"data":    tokenString,
	})
}

// FindUser ...
func FindUser(context *gin.Context) {
	var user user
	if e := context.Bind(&user); e != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Err"})
		panic(e)
	}

	userDataR, err := models.Client.Get(user.Email).Result()
	if err == nil {
		context.JSON(http.StatusOK, userDataR)
		return
	}

	userDataD, err := models.FindUser(map[string]interface{}{"email": user.Email})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		panic(err)
	}

	err = models.Client.Set(user.Email, userDataD.ID.Hex(), 0).Err()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
		panic(err)
	}

	context.JSON(http.StatusOK, user)
}

// AddMulUserEndPoint allow add new user data into User Collection
func AddMulUserEndPoint(context *gin.Context) {
	for i := 0; i < 1000; i++ {
		_, errorAdding := models.AddNewUser("user" + strconv.Itoa(i) + "@email.com")
		if errorAdding != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error"})
			panic(errorAdding)
		}
	}
	context.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

// FindMul allow add new user data into User Collection
func FindMul(context *gin.Context) {
	for i := 0; i < 500; i++ {
		_, err := models.Client.Get("user" + strconv.Itoa(i) + "@email.com").Result()
		if err != nil {

			userDataD, err := models.FindUser(map[string]interface{}{"email": "user" + strconv.Itoa(i) + "@email.com"})
			if err != nil {
			}

			err = models.Client.Set("user"+strconv.Itoa(i)+"@email.com", userDataD.ID.Hex(), 0).Err()
			if err != nil {

			}
		}
	}
	context.JSON(http.StatusOK, gin.H{"message": "Ok"})
}

// EraseRedis ...
func EraseRedis(context *gin.Context) {
	for i := 0; i < 500; i++ {
		_ = models.Client.Del("user" + strconv.Itoa(i) + "@email.com")
	}
	context.JSON(http.StatusOK, gin.H{"message": "Ok"})
}
