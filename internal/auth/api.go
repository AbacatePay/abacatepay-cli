package auth

import "github.com/go-resty/resty/v2"

type User struct {

	Name  string

	Email string

}



// TODO: I`ll made this func when i get the right endpoint later

func ValidateToken(client *resty.Client, baseURL, token string) (*User, error) {

	return &User{Name: "Mock User", Email: "mock@example.com"}, nil

}
