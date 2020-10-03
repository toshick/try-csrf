package util

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

var IsDev bool = false

func SetCookie(c echo.Context, name string, value string) error {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Secure = false
	cookie.HttpOnly = false
	cookie.Path = "/"
	c.SetCookie(cookie)
	return nil
}

func DumpCookies(c echo.Context) error {
	fmt.Println("\n--------------------- DumpCookies --------------------- ")

	for _, cookie := range c.Cookies() {

		fmt.Println("cookie: ", cookie.Name, cookie.Value)
	}
	fmt.Println("-------------------------------------------------------")
	return nil
}
