package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/keshav/fiber/initializers"
	"github.com/keshav/fiber/models"

	"github.com/keshav/fiber/routes.js"
)

func init() {
	initializers.LoadVariable()
	initializers.ConnectToDB()
}

func main() {
	db, _ := initializers.ConnectToDB()

	defer db.Close(context.Background())
	// Initialize Fiber

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/role", func(c *fiber.Ctx) error {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			return err
		}
		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c.Response().BodyWriter(), nil)
	})

	app.Post("/roleSubmit", func(c *fiber.Ctx) error {

		role := c.FormValue("role")
		//email := c.FormValue("email")

		_, err := db.Exec(context.Background(), "INSERT INTO role (role_name) VALUES ($1)", role)

		if err != nil {
			fmt.Println("errrr", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error inserting data into the database",
			})
		}

		fmt.Println("Role ", role)

		return c.Redirect("/rolelisting")

	})

	app.Get("/rolelisting", func(c *fiber.Ctx) error {

		query := `Select role_name from role`
		row, err := db.Query(context.Background(), query)

		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to get data",
			})
		}

		//var RolesDatas []models.Role
		var RolesDatas []string
		defer row.Close()

		for row.Next() {

			//var RoleData models.Role
			var RoleData string
			if err = row.Scan(&RoleData); err != nil {
				fmt.Println("erros", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to scan row data",
				})
			}

			RolesDatas = append(RolesDatas, RoleData)
		}

		fmt.Println("ROwwwws", RolesDatas)
		return c.Render("listingPage", fiber.Map{
			"Title": "Listing Data",
			"Data":  RolesDatas,
		})

	})

	app.Get("/permission", func(c *fiber.Ctx) error {
		tmpl, err := template.ParseFiles("Permission.html")
		if err != nil {
			return err
		}
		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c.Response().BodyWriter(), nil)
	})

	app.Post("/permissionSubmit", func(c *fiber.Ctx) error {
		permi := c.FormValue("permission")
		//email := c.FormValue("email")

		_, err := db.Exec(context.Background(), "INSERT INTO permissions (permission_name) VALUES ($1)", permi)

		if err != nil {
			fmt.Println("errrr", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error inserting data into the database",
			})
		}
		return c.Redirect("/permissionlisting")

	})

	app.Get("/permissionlisting", func(c *fiber.Ctx) error {

		query := `Select permission_name from permissions`
		row, err := db.Query(context.Background(), query)

		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to get data",
			})
		}

		// type Permission struct {
		// 	Name string `json:"permission_name"`
		// }
		// var PermissionNames []Permission
		var PermissionNames []string

		defer row.Close()

		for row.Next() {

			// var PermissionName Permission
			var PermissionName string
			if err = row.Scan(&PermissionName); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to scan row data",
				})
			}

			PermissionNames = append(PermissionNames, PermissionName)
		}

		fmt.Println("ROwwwws", PermissionNames)
		return c.Render("listingPage", fiber.Map{
			"Title": "Listing Data",
			"Data":  PermissionNames,
		})

	})

	app.Get("/roles_permission", func(c *fiber.Ctx) error {

		///// Roles Area
		query := `Select id,role_name from role`
		row, err := db.Query(context.Background(), query)

		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to get data",
			})
		}

		var RolesDatas []models.Role
		// type Roles struct{
		// 	Id int `json:"id"`
		// 	Name string `json:"name"`
		// }
		//var RolesDatas []string
		//var RolesDatas []Roles
		defer row.Close()
		var RoleData models.Role
		for row.Next() {

			//var RoleData models.Role
			//var RoleData string

			if err = row.Scan(&RoleData.Id, &RoleData.Role); err != nil {
				fmt.Println("erros", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to scan row data",
				})
			}

			RolesDatas = append(RolesDatas, RoleData)
		}

		fmt.Println("ROles rows ", RolesDatas)

		///Permission Area
		query = `Select id,permission_name from permissions`
		row, err = db.Query(context.Background(), query)

		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to get data",
			})
		}

		var PermissionNames []models.Permissions

		defer row.Close()

		for row.Next() {

			// var PermissionName Permission
			var PermissionName models.Permissions
			if err = row.Scan(&PermissionName.Id, &PermissionName.Permission); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to scan row data",
				})
			}

			PermissionNames = append(PermissionNames, PermissionName)
		}

		fmt.Println("Permisiions Rows", PermissionNames)

		return c.Render("role_permission", fiber.Map{
			"Title": "Listing Data",
			"Data1": RolesDatas,
			"Data2": PermissionNames,
		})
	})

	// app.Get("/role_permission", func(c *fiber.Ctx) error {

	// 	tmpl, err := template.ParseFiles("role_permission.html")

	// 	if err != nil {
	// 		return err
	// 	}
	// 	c.Set("Content-Type", "text/html")
	// 	return tmpl.Execute(c.Response().BodyWriter(), nil)
	// })

	app.Post("/role_permissionSubmit", func(c *fiber.Ctx) error {

		// roleMap := map[string]int{
		// 	"option1": 1, // Replace with actual IDs
		// 	"option2": 2,
		// 	"option3": 3,
		// }

		// permissionMap := map[string]int{
		// 	"option1": 1, // Replace with actual IDs
		// 	"option2": 2,
		// 	"option3": 3,
		// }

		roleval := c.FormValue("rolevalue")
		permival := c.FormValue("permissionvalue")

		fmt.Println("Rolllleeval", roleval)

		fmt.Println("PermissionVal", permival)

		// roleID, roleExists := roleMap[roleval]
		// permID, permExists := permissionMap[permival]

		// if !roleExists || !permExists {
		// 	// Handle error: Invalid role/permission value
		// 	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		// 		"error": "Invalid role/permission value",
		// 	})
		// }

		_, err := db.Exec(context.Background(), "INSERT INTO role_permission (role_id,permission_id) values($1,$2)", roleval, permival)

		if err != nil {
			fmt.Println("----error------", err)
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Error inserting into database in permission submit",
			})
		}

		return c.Redirect("/rolepermissionlisting")
	})

	app.Get("/rolepermissionlisting", func(c *fiber.Ctx) error {
		query := `SELECT id,role_id,permission_id from role_permission`
		row, err := db.Query(context.Background(), query)

		if err != nil {
			c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "error in database in role_permisiionListing",
			})
		}

		var rpDatas []models.Role_Permission
		defer row.Close()

		for row.Next() {
			var rpData models.Role_Permission
			if err = row.Scan(&rpData.Id, &rpData.Role_id, &rpData.Permission_id); err != nil {
				fmt.Println("errorrrr", err)
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "error from getting the data from database in rolepermilisting",
				})
			}
			rpDatas = append(rpDatas, rpData)
		}

		fmt.Println("ROwwwws", rpDatas)

		return c.Render("listingPage", fiber.Map{
			"Title": "Listing Data",
			"Data":  rpDatas,
		})
	})

	routes.SetupAdminRoutes(app)
	routes.SetupUserRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
