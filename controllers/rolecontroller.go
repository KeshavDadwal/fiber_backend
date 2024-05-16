package controllers

// import (
// 	"context"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/log"
// 	"github.com/keshav/fiber/initializers"
// 	"github.com/keshav/fiber/models"
// )

// func VariousRole(c *fiber.Ctx) error{
// 	db,_ := initializers.ConnectToDB()

// 	var roles models.Role
// 	if err := c.BodyParser(&roles);  err!=nil{
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request body",
// 		})
// 	}

// 	_, err := db.Exec(context.Background(),"INSERT INTO role (role_name) VALUES ($1)",roles.Role)

// 	err = db.QueryRow(context.Background(),"SELECT lastval()").Scan(&roles.Id)
// 	if err!=nil{
// 		log.Fatal(err)
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"Message":"Role created",
// 	})

// }