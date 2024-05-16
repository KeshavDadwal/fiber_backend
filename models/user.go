package models

type User struct {
	Id             int    `json:"id"`
	First_name     string `json:"first_name"`
	Last_name      string `json:"last_name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Age            int    `json:"age"`
	Phone_no       int    `json:"phone_no"`
	Secret_code    string `json:"secret_code"`
	Role_id        int    `json:"role_id"`
}

type Role struct{
	Id int `json:"id"`
	Role string `json:"role_name"`
}
type Permissions struct{
	Id int `json:"id"`
	Permission string `json:"permission_name"`
}
type Role_Permission struct {
	Id int `json:"id"`
	Role_id int `json:"role_id"`
	Permission_id int `json:"permission_id"`
}
type Images struct{
	Image_url string `gorm:"image_url"`
}