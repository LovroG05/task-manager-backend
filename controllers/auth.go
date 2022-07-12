package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovrog05/task-manager-backend/models"
	"github.com/lovrog05/task-manager-backend/models/inputs"
	"github.com/lovrog05/task-manager-backend/utils/token"
	"golang.org/x/crypto/bcrypt"
)

func CurrentUser(c *gin.Context) {

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func Login(c *gin.Context) {

	var input inputs.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password
	token, err := models.LoginCheck(u.Username, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func UpdateFmcToken(c *gin.Context) {
	var input inputs.FmcTokenInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.FmcToken = input.Token
	if err := u.UpdateFmcToken(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Register(c *gin.Context) {

	var input inputs.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	passwd, errr := bcrypt.GenerateFromPassword([]byte(models.PreparePasswordInput(input.Password)), bcrypt.DefaultCost)
	if errr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errr.Error()})
		return
	}

	u.Username = input.Username
	u.Password = string(passwd)

	_, err := u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})

}
