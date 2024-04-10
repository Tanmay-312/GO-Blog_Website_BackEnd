package controller

import (
	"blog-website/database"
	"blog-website/models"
	"blog-website/util"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePost(c *fiber.Ctx) error {
	var blogpost models.Blog
	if err := c.BodyParser(&blogpost); err != nil {
		fmt.Println("Unable to parse body")
	}

	if err := database.DB.Create(&blogpost).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid payload"})
	}

	return c.JSON(fiber.Map{"Success": "Congratulations, your post is live."})
}

func AllPost(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 5
	offset := (page - 1) * limit

	var total int64
	var getBlog []models.Blog
	database.DB.Preload("User").Offset(offset).Limit(limit).Find(&getBlog)
	database.DB.Model(&models.Blog{}).Count(&total)
	return c.JSON(fiber.Map{
		"Data": getBlog, "Meta": fiber.Map{
			"Total":     total,
			"Page":      page,
			"Last_page": ceil(int(total) / limit),
		},
	})
}

func DetailPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var blogpost models.Blog
	database.DB.Where("id=?", id).Preload("User").First(&blogpost)
	return c.JSON(fiber.Map{"Data": blogpost})
}

func UpdatePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	blog := models.Blog{
		Id: uint(id),
	}

	if err := c.BodyParser(&blog); err != nil {
		fmt.Println("Unable to parse body")
	}

	database.DB.Model(&blog).Updates(blog)

	return c.JSON(fiber.Map{"Success": "Updated successfully", "Data": blog})
}

func UniquePost(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	id, _ := util.ParseJwt(cookie)
	var blog []models.Blog
	database.DB.Model(&blog).Where("user_id=?", id).Preload("User").Find(&blog)

	return c.JSON(blog)
}

func DeletePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	blog := models.Blog{
		Id: uint(id),
	}

	deleteQuery := database.DB.Delete(&blog)
	if errors.Is(deleteQuery.Error, gorm.ErrRecordNotFound) {
		c.Status(400)
		return c.JSON(fiber.Map{"Error": "Records not found"})
	}

	return c.Status(200).JSON(fiber.Map{"Success": "Post deleted successfully"})
}

func ceil(n int) int {
	x := float64(n)
	if x < 0 {
		return int(x - 0.5)
	}
	return int(x + 0.5)
}
