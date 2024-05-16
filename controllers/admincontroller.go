package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/keshav/fiber/initializers"
	"github.com/keshav/fiber/maill"
	"github.com/keshav/fiber/models"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

func AdminSignUp(c *fiber.Ctx) error {

	db, _ := initializers.ConnectToDB()
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

	err = maill.SendEmailWithGmail(user.Email, randomCode, user.First_name, user.Id)
	if err != nil {
		fmt.Println("error in sending the mail")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": user,
	})
}
func Adminlogin(c *fiber.Ctx) error {
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

	fmt.Println("Login")

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"messgae": "I m logged In",
		"token":   tokenString,
	})

}
func VariousRole(c *fiber.Ctx) error {
	db, _ := initializers.ConnectToDB()

	var roles models.Role
	if err := c.BodyParser(&roles); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, err := db.Exec(context.Background(), "INSERT INTO role (role_name) VALUES ($1)", roles.Role)

	err = db.QueryRow(context.Background(), "SELECT lastval()").Scan(&roles.Id)
	if err != nil {
		log.Fatal(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"Message": "Role created",
	})

}

func UploadHandler(c *fiber.Ctx) error{

        file,err := c.FormFile("video")
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
        }
		log.Println(file.Filename);

        // Create a new file to save the uploaded video
        uploadedFile, err := os.Create(file.Filename)
        if err != nil {
            c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"error in creating a file",
			})
        }
        defer uploadedFile.Close()

			newFileName := "./uploads/" + file.Filename
			newFile, err := os.Create(newFileName)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Error creating a file",
				})
			}
			defer newFile.Close()

			// Copy the uploaded content to the new file
			_, err = io.Copy(newFile, uploadedFile)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Error copying file contents",
				})
			}
			if _, err := os.Stat(newFileName); os.IsNotExist(err) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Video not found",
				})
			}
		
			return c.SendFile(newFileName)
    
}


func UploadImage(c *fiber.Ctx) error {
	// single file
	db, _ := initializers.ConnectToDB()

	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	var image struct {
		Image_url string
	}

	if err := c.BodyParser(&image); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file",
		})
	}

	randomName := generateRandomFileName(filepath.Ext(file.Filename))

	img := models.Images{
		Image_url: randomName,
	}

	_, err := db.Exec(context.Background(), "INSERT INTO images (image_url) VALUES ($1)", img.Image_url)

	if err != nil {
		log.Fatal(err)
	}

	f, _ := file.Open()

	tempFile, err := ioutil.TempFile("uploads", "upload-*.jpg")
	fmt.Println("-----temp-----", tempFile)
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "something went wrong",
		})
	}

	fmt.Println("Fileupload successfully")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"Message": "Image Uploaded !!",
	})
}



func generateRandomFileName(extension string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)

	timestamp := time.Now().UnixNano()
	timestampStr := strconv.FormatInt(timestamp, 16)

	return hex.EncodeToString(randBytes) + "-" + timestampStr + extension
}
