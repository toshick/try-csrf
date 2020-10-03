package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/toshick/ronten-maker/app/model"
)

/**
 * /GetUsers
 */
func GetUsers(c echo.Context) error {

	// パラメータ取得
	// hash := c.Param("hash")
	// if hash == "" {
	// 	return c.JSON(http.StatusBadRequest, &model.ApiError{Error: true, Message: "cant find hash"})
	// }

	// util.SetCookie(c, "get-users", "yes")

	userlist := []model.User{}
	r1 := model.User{ID: "001", Name: "タロー1", Email: "taro001@taro.com", Pass: "xxx"}
	userlist = append(userlist, r1)
	r2 := model.User{ID: "002", Name: "タロー2", Email: "taro002@taro.com", Pass: "xxx"}
	userlist = append(userlist, r2)

	ret := model.Users{Users: userlist}

	return c.JSON(http.StatusOK, ret)
}

/**
 * /PostUsers
 */
func PostUsers(c echo.Context) error {

	fmt.Printf("リクエスト成功（PostUsers） \n")

	return c.JSON(http.StatusOK, &model.ApiError{Error: false, Message: "ok"})

}
