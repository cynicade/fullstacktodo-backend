package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/cynicade/todo/database"
	"github.com/cynicade/todo/models"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func createToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    id,
		ExpiresAt: jwt.At(time.Now().Add(time.Hour * 48)),
	})

	ss, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	return ss, err
}

func createCookie(id string) (fiber.Cookie, error) {
	ss, err := createToken(id)

	if err != nil {
		return fiber.Cookie{}, err
	}

	return fiber.Cookie{
		Name:     "token",
		Value:    ss,
		Expires:  time.Now().Add(time.Hour * 48),
		HTTPOnly: true,
	}, nil
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": "could not process request",
		})
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not hash password",
		})
	}

	user := models.User{
		Username: data["username"],
		Email:    data["email"],
		Password: string(pwd),
		Id:       uuid.New(),
	}

	if res := database.DB.Table("users").Create(&user); res.Error != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not save user to database",
		})
	}

	cookie, err := createCookie(user.Id.String())
	if err != nil {
		// unsure which status code is the correct choice here
		c.Status(fiber.StatusTeapot)
		return c.JSON(fiber.Map{
			"message": "user saved to database, but could not generate cookie",
		})
	}

	c.Status(fiber.StatusCreated)
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "user registered successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	var user models.User

	if clientCookie := c.Cookies("token"); clientCookie != "" {

		clientToken, err := jwt.ParseWithClaims(clientCookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_SECRET")), nil
		})

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "error processing authentication cookie",
			})
		}

		if claims, ok := clientToken.Claims.(*jwt.StandardClaims); ok && clientToken.Valid {
			res := database.DB.Table("users").Where("id = ?", claims.Issuer).First(&user)

			if res.RowsAffected == 0 {
				c.Status(fiber.StatusUnauthorized)
				c.ClearCookie("token")
				return c.JSON(fiber.Map{
					"message": "user not found",
				})
			}

			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"message": "logged in successfully",
			})
		} else {
			c.Status(fiber.StatusUnauthorized)
			c.ClearCookie("token")
			fmt.Println(fmt.Errorf("error with cookie: %w", err))
			return c.JSON(fiber.Map{
				"message": "invalid authentication cookie",
			})
		}
	}

	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusUnprocessableEntity)
		return c.JSON(fiber.Map{
			"message": "could not process request",
		})
	}

	res := database.DB.Table("users").Where("email = ?", data["email"]).First(&user)

	if res.RowsAffected == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	cookie, err := createCookie(user.Id.String())

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not generate cookie",
		})
	}

	c.Status(fiber.StatusOK)
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "logged in successfully",
	})
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("token")
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "logged out",
	})
}
