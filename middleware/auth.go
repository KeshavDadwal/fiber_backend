package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/initializers"
	"github.com/keshav/fiber/models"

	//	"github.com/keshav/fiber/models"
	"github.com/o1egl/paseto"
	//"github.com/pelletier/go-toml/query"
)

// type User struct {
// 	Id          int    `gorm:"id"`
// 	First_name  string `gorm:"firstname"`
// 	Last_name   string `gorm:"lastname"`
// 	Email       string `gorm:"unique"`
// 	Password    string `gorm:"password"`
// 	Age         int    `gorm:"age"`
// 	Phone_no    int    `gorm:"phone_no"`
// 	Secret_code string `gorm:"secret_code"`
// 	Role_id     int    `gorm:"role_id"`
// }

func RequireAuth(c *fiber.Ctx) error {

	tokenString := c.Get("x-token")

	db, _ := initializers.ConnectToDB()

	var decryptedToken paseto.JSONToken

	symmetricKey := os.Getenv("SECRET")

	var newFooter string
	err := paseto.NewV2().Decrypt(tokenString, []byte(symmetricKey), &decryptedToken, &newFooter)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to decrpt the token",
		})

	}

	if time.Now().After(decryptedToken.Expiration) {
		fmt.Println("Token Expired")

	}

	idd, err := strconv.Atoi(decryptedToken.Subject)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to convert to int",
		})

	}

	var Dataa models.User

	query := `Select id,email,role_id from users where id = $1`

	row := db.QueryRow(context.Background(), query, idd)

	err = row.Scan(&Dataa.Id,&Dataa.Email,&Dataa.Role_id)
	if err != nil {
		fmt.Println("errrrrr", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Data is not there",
		})
	}
	

	fmt.Println("Data of user",Dataa)

	var roleId int
	query = `Select role_id from users where id = $1`
	row = db.QueryRow(context.Background(), query, idd)

	err = row.Scan(&roleId)

	if err != nil {
		fmt.Println("errrrrr", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid username and password",
		})
	}

	c.Locals("roleId", strconv.Itoa(roleId))
	c.Locals("id", strconv.Itoa(idd))
	c.Locals("User",Dataa)

	c.Next()

	return nil
}
