package controllers

import (
	"os"

	"github.com/cynicade/todo/database"
	"github.com/cynicade/todo/models"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func getIdFromCookie(cookie string) string {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return "0"
	}

	var user models.User

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		res := database.DB.Table("users").Where("id = ?", claims.Issuer).First(&user)

		if res.RowsAffected == 0 {
			return "0"
		}

		return claims.Issuer
	}

	return "0"
}

func GetAll(c *fiber.Ctx) error {
	if userid := getIdFromCookie(c.Cookies("token")); userid != "0" {
		var todos []models.Todo

		res := database.DB.Table("todos").Where("User_ID = ?", userid).Find(&todos)

		if res.Error != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not load todos",
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"todos": todos,
		})
	}

	c.Status(fiber.StatusUnauthorized)
	return c.JSON(fiber.Map{
		"message": "you are not logged in",
	})
}

func NewTodo(c *fiber.Ctx) error {
	if userid := getIdFromCookie(c.Cookies("token")); userid != "0" {
		var data map[string]string
		userid, err := uuid.Parse(userid)

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "error while trying to parse user id",
			})
		}

		if err := c.BodyParser(&data); err != nil {
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(fiber.Map{
				"message": "could not process request",
			})
		}

		todo := models.Todo{
			Body:     data["body"],
			Complete: false,
			User_ID:  userid,
		}

		if res := database.DB.Table("todos").Create(&todo); res.Error != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not save todo to database",
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "saved todo to database successfully",
		})
	}

	c.Status(fiber.StatusUnauthorized)
	return c.JSON(fiber.Map{
		"message": "you are not logged in",
	})
}

func Update(c *fiber.Ctx) error {
	if userid := getIdFromCookie(c.Cookies("token")); userid != "0" {
		var todo models.Todo
		var data map[string]string

		res := database.DB.Table("todos").Where("User_ID = ? AND ID = ?", userid, c.Params("id")).Find(&todo)

		if res.Error != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not load todo",
			})
		}

		if err := c.BodyParser(&data); err != nil {
			c.Status(fiber.StatusUnprocessableEntity)
			return c.JSON(fiber.Map{
				"message": "could not process request",
			})
		}

		todo.Body = data["body"]
		if data["complete"] == "true" {
			todo.Complete = true
		} else {
			todo.Complete = false
		}

		updres := database.DB.Table("todos").Save(todo)

		if updres.Error != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not update todo",
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"todo": todo,
		})
	}

	c.Status(fiber.StatusUnauthorized)
	return c.JSON(fiber.Map{
		"message": "you are not logged in",
	})
}

func Delete(c *fiber.Ctx) error {
	if userid := getIdFromCookie(c.Cookies("token")); userid != "0" {

		res := database.DB.Table("todos").Delete(&models.Todo{}, c.Params("id"))

		if res.Error != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not delete todo",
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "deleted todo",
		})
	}

	c.Status(fiber.StatusUnauthorized)
	return c.JSON(fiber.Map{
		"message": "you are not logged in",
	})
}
