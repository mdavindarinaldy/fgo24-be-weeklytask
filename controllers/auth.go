package controllers

import (
	"backend3/models"
	"backend3/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func AuthRegister(c *gin.Context) {
	user := models.User{}
	c.ShouldBind(&user)
	err := HandleRegister(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to register user",
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Create user success!",
	})
}

func AuthLogin(c *gin.Context) {
	user := struct {
		Email    string `form:"email" json:"email" db:"email" binding:"required,email"`
		Password string `form:"password" json:"password" db:"password" binding:"required"`
		Pin      string `form:"pin" json:"pin" db:"pin" binding:"required"`
	}{}
	c.ShouldBind(&user)
	userData, err := CheckUser(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Internal server error",
		})
		return
	}
	if userData == (models.User{}) {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "User is not registered",
		})
		return
	}
	if userData.Email == user.Email && userData.Password == user.Password && userData.Pin == user.Pin {
		// token := models.Token
		// *token, err = GenerateToken()

		token, err := GenerateToken()

		if err != nil {
			c.JSON(http.StatusBadRequest, utils.Response{
				Success: false,
				Message: "Failed to generate token",
			})
			return
		}
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Message: "Login success!",
			Result:  token,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Password and/or PIN is wrong!",
		})
	}
}

func CheckUser(email string) (models.User, error) {
	conn, err := utils.DBConnect()
	if err != nil {
		return models.User{}, err
	}
	defer conn.Close()
	rows, err := conn.Query(context.Background(),
		`SELECT * FROM users WHERE email=$1`, email)

	if err != nil {
		return models.User{}, err
	}
	user, err := pgx.CollectOneRow[models.User](rows, pgx.RowToStructByName)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func HandleRegister(user models.User) error {
	conn, err := utils.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(
		context.Background(),
		`
		INSERT INTO users (name, email, phone_number, password, pin)
		VALUES
		($1,$2,$3,$4,$5);
		`,
		user.Name, user.Email, user.PhoneNumber, user.Password, user.Pin)
	if err != nil {
		return err
	}
	return nil
}

func GenerateToken() (string, error) {
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": 0,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(15 * time.Minute),
	})
	token, err := generateToken.SignedString([]byte(os.Getenv("APP_SECRET")))
	if err != nil {
		return token, err
	}
	return token, nil
}

func VerifyToken(token string) bool {
	godotenv.Load()
	secretKey := os.Getenv("APP_SECRET")
	rawToken, _ := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	fmt.Println(rawToken.Claims.(jwt.MapClaims)["userId"])
	fmt.Println(rawToken.Claims.(jwt.MapClaims))
	return rawToken.Claims.(jwt.MapClaims)["userId"] == float64(0)
}

func HandleToken(c *gin.Context) {
	token := struct {
		Value string `form:"token"`
	}{}
	c.ShouldBind(&token)
	fmt.Println(VerifyToken(token.Value))
	if VerifyToken(token.Value) {
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Message: "berhasil yeay",
		})
	} else {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "gagal lah kocak",
		})
	}
}
