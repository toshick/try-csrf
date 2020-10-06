package routers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/toshick/try-csrf/app/context"
	"github.com/toshick/try-csrf/app/controllers"
	"github.com/toshick/try-csrf/app/model"
)

/**
 * SetRoutes
 */
func SetRoutes() {
	e := echo.New()
	e.Validator = model.NewValidator()

	// session
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	/**
	 * cookieチェック
	 */
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// util.DumpCookies(c)

			fmt.Println("\n--------------------- Cookie --------------------- ")
			fmt.Printf("c.Cookie() %v \n", c.Request().Cookies())
			fmt.Printf("Method %v  Referer %v \n", c.Request().Method, c.Request().Referer())
			return next(c)
		}
	})

	/**
	 * CORSの許可設定
	 */
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
		AllowHeaders:     []string{"*"}, //jsに特定のヘッダ情報取得を許可するための設定
	}))

	/**
	 * CSRFの対策
	 */
	// e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{}))

	/**
	 * HandlerFunc オリジナルコンテキストをセット
	 */
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &context.CustomContext{c}
			return next(cc)
		}
	})

	/**
	 * HandlerFunc CheckAuth
	 */
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// cc := c.(*context.CustomContext)

			// if cc.Path() == "/api/login" || cc.Path() == "/api/logout" {
			// 	return next(c)
			// }
			// secretHeader := cc.Request().Header.Get("X-XSRF-TOKEN")
			// if cc.Request().Method == "POST" && secretHeader != "SECRET" {
			// 	return cc.JSON(http.StatusUnauthorized, &model.ApiError{Error: true, Message: "It is CSRF!"})
			// }

			// fmt.Println("Method", cc.Request().Method, secretHeader)

			// 認証チェック
			// user := cc.LoginedUser()
			// fmt.Printf("認証チェック %+v  %v \n", user, cc.Path())
			// if user.ID == "" {
			// 	return cc.JSON(http.StatusUnauthorized, &model.ApiError{Error: true, Message: "need login"})
			// }

			return next(c)
		}
	})

	// var addr = flag.String("addr", ":8888", "http service address")
	// fmt.Printf("addr %v \n", addr)
	// err := http.ListenAndServe(*addr, nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
	// return

	// routing for front
	e.Static("/", "./public")

	// e.GET("/discussion/:hash", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "/discussion/:hash")
	// })

	//----------------------
	// api
	//----------------------
	g := e.Group("/api")
	// project
	g.GET("/users/", controllers.GetUsers)
	g.POST("/users/", controllers.PostUsers)

	port := "80"
	var p string
	p = os.Getenv("PORT")
	if p != "" {
		port = p
	}
	p = os.Getenv("HTTPPORT")
	if p != "" {
		port = p
	}

	fmt.Printf("port %v \n", port)

	port = fmt.Sprintf(":%s", port)

	e.Logger.Fatal(e.Start(port))
}
