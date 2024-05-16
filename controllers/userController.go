package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/initializers"
	// "github.com/keshav/fiber/maill"
	"github.com/keshav/fiber/models"
	"github.com/o1egl/paseto"

	//"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

// type User struct {
// 	Id             int    `json:"id"`
// 	First_name     string `json:"first_name"`
// 	Last_name      string `json:"last_name"`
// 	Email          string `json:"email"`
// 	Password       string `json:"password"`
// 	Age            int    `json:"age"`
// 	Phone_no       int    `json:"phone_no"`
// 	Secret_code    string `json:"secret_code"`
// 	Role_id        int    `json:"role_id"`
// }

func UserSignUp(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()

	// Parse the request body into a User struct
	var user models.User
	if err := c.BodyParser(&user); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to hash password",
		})
	}

	randomCode, err := generateRandomCode(6)
	if err != nil {
		fmt.Println("Error generating random code:", err)
	}
	// fmt.Println("randommmmm",randomCode)
	// Insert the user into the database
	_, err = db.Exec(context.Background(),
		"INSERT INTO users (first_name,last_name,email,password,age,phone_no,secret_code,role_id) VALUES ($1, $2, $3, $4, $5, $6, $7,$8)",
		user.First_name, user.Last_name, user.Email, string(hash), user.Age, user.Phone_no, randomCode, user.Role_id)

	err = db.QueryRow(context.Background(), "SELECT lastval()").Scan(&user.Id)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println("errrr", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error inserting data into the database",
		})
	}
	// var user1 User
	// query := `SELECT id FROM users WHERE email = $1`
	// row := db.QueryRow(context.Background(), query, user.Email)

	// err = row.Scan(&user1.Id)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": "Invalid username and password",
	// 	})}

	/// mail sendcode
	// err = maill.SendEmailWithGmail(user.Email, randomCode, user.First_name, user.Id)
	// if err != nil {
	// 	fmt.Println("error in sending the mail")
	// }

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": user,
	})
}

func UserLogin(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()

	var body struct {
		Email    string
		Password string
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	fmt.Println("Boddddyyy data", body)
	var user1 models.User

	query := `SELECT id, password FROM users WHERE email = $1`
	row := db.QueryRow(context.Background(), query, body.Email)

	err := row.Scan(&user1.Id, &user1.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid username and password",
		})
	}

	if user1.Id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid username and password",
		})
	}

	userID := strconv.Itoa(user1.Id)

	err = bcrypt.CompareHashAndPassword([]byte(user1.Password), []byte(body.Password))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid username and password",
		})
	}

	symmetricKey := os.Getenv("SECRET")
	now := time.Now()
	exp := now.Add(24 * time.Hour)
	nbt := now

	jsonToken := paseto.JSONToken{
		Jti:        "123",
		Subject:    userID,
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}

	fmt.Println("jsonnnnTokennData",jsonToken);
	// Encrypt data
	jsonToken.Set("data", "this is a signed message")
	jsonToken.Set("email", user1.Email)

	footer := "some footer"
	tokenString, err := paseto.NewV2().Encrypt([]byte(symmetricKey), jsonToken, footer)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to create a token",
		})
	}

	//c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"messgae": "I m logged In",
		"token":   tokenString,
	})

}

func GetAllUser(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()

	type Data1 struct{
		Id int
		Email string
		RoleId int
	}

	// var UseddDtaa Data1
	// UseddDtaa = c.Locals("User")



	// fmt.Println("UsedddData",UseddDtaa);

	RoleId := c.Locals("roleId").(string)
	Id, err := strconv.Atoi(RoleId)

	if Id != 1 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "You are not authorized",
		})
	}

	row, err := db.Query(context.Background(), "select id,first_name,last_name,email,password,age,phone_no,role_id,secret_code from users ")

	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to get data",
		})
	}

	defer row.Close()

	var users []models.User

	for row.Next() {
		var user models.User

		if err = row.Scan(&user.Id, &user.First_name, &user.Last_name, &user.Email, &user.Password, &user.Age, &user.Phone_no, &user.Role_id, &user.Secret_code); err != nil {
			fmt.Println("erros", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan row data",
			})
		}

		users = append(users, user)
	}

	if err := row.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error during iteration of rows",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"data": users,
	})

}
func generateRandomCode(length int) (string, error) {
	// Calculate the number of bytes needed for the random code
	numBytes := (length * 6) / 8

	// Generate random bytes
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	randomCode := base64.URLEncoding.EncodeToString(randomBytes)

	// Trim the random code to the desired length
	randomCode = randomCode[:length]

	return randomCode, nil
}

func VerifyApi(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()
	id := c.Params("id")
	secretCode := c.Params("secret_code")

	var user1 models.User
	query := `Select secret_code from users where id = $1`
	row := db.QueryRow(context.Background(), query, id)
	err := row.Scan(&user1.Secret_code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid username and password",
		})
	}

	if user1.Secret_code == secretCode {

		query = `Update users
		set secret_code = $1 , email_verified = $2
		where id = $3
		`
		result, err := db.Exec(context.Background(), query,
			"null",
			"true",
			id,
		)
		if err != nil {
			return err
		}

		rowAffected := result.RowsAffected()

		if rowAffected == 0 {
			c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		} else {
			return c.Render("verification", fiber.Map{
				"Title": "Hello, World!",
			})

		}

	} else {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "secret code mismatch",
		})
	}

	return nil
}

func UpdateOneUser(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()
	id := c.Params("id")

	RoleId := c.Locals("roleId").(string)
	Id, err := strconv.Atoi(RoleId)

	if Id != 1 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "You are not authorized",
		})
	}

	var body models.User

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read body",
		})
	}

	query := `
	UPDATE users
	SET first_name = $1, last_name = $2, email = $3, password = $4, age = $5, phone_no = $6
	WHERE id = $7
	 `

	result, err := db.Exec(context.Background(), query,
		body.First_name,
		body.Last_name,
		body.Email,
		body.Password,
		body.Age,
		body.Phone_no,
		id,
	)

	if err != nil {
		// Handl
		return err
	}
	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "User Updated",
	})

	return nil
}

func Deletedata(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()

	RoleId := c.Locals("roleId").(string)
	Id, err := strconv.Atoi(RoleId)

	if Id != 1 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "You are not authorized",
		})
	}

	id := c.Params("id")

	query := `Delete from users where id = $1`

	result, err := db.Exec(context.Background(), query, id)

	if err != nil {
		c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "unable to execute the query",
		})
	}
	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "User Deleted",
	})

	return nil
}
