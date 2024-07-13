package controllers

import (
	"strconv"
	"time"
	"github.com/devmehta22/JWT-Auth/database"
	"github.com/devmehta22/JWT-Auth/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)




func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to parse request body",
        })
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?",data["email"]).First( &existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already exists",
		})
	}

	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}
	user := models.User{
		Username: data["username"],
		Email:    data["email"],
		Password: hashedPassword,
		}

	if err := database.DB.Create(&user).Error; err!=nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	var user models.User
	if err := database.DB.Where("email = ?",data["email"]).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}
	
	//GENERATE JWT-TOKEN
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"sub":strconv.Itoa(int(user.Id)),
		"exp":time.Now().Add(time.Hour*24).Unix(),
	})
	token,err := claims.SignedString([]byte(database.DotEnv("SECRETKEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
		})
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour*24),
		HTTPOnly: true,
		Secure: true,
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Logged in successfully",
	})	
}

func Protected(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, 
		func(token *jwt.Token) (interface{}, error) {
			return []byte(database.DotEnv("SECRETKEY")), nil
			})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	claims,ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	id,_ := strconv.Atoi((*claims)["sub"].(string))
	user := models.User{Id: int(id)}

	database.DB.Where("id=?",id).First( &user)

	return c.Status(fiber.StatusOK).JSON(user)
}

func Logout(c *fiber.Ctx) error  {
	cookie := fiber.Cookie{
		Name: "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour*24),
		HTTPOnly: true,
		Secure: true,
	}
	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
		})
}