package controller

import (
	"blog-website/database"
	"blog-website/models"
	"blog-website/util"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`)
	return re.MatchString(email)
}

func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User

	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	// check if password is less than 6 characters
	if len(data["password"].(string)) < 6 || data["password"] == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 6 characters long."})
	}

	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format."})
	}

	// check if email already exists in database
	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.Id > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "This email address has been registered."})
	}

	// create new user and save it into the database
	user := models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
		Phone:     data["phone"].(string),
	}

	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user)
	if err != nil {
		log.Println(err)
	}

	return c.Status(200).JSON(fiber.Map{"success": "Account created successfully.", "user": user})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.SendString("Unable to parse body")
	}

	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{"error": "Email does not exist, please create a account."})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{"error": "Incorrect password"})
	}

	token, err := util.GenerateJwt(strconv.Itoa(int(user.Id)))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"success": "You have successfully logged in", "user": user})
}

type Claims struct {
	jwt.StandardClaims
}
