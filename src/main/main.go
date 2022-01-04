package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/dgrijalva/jwt-go"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"Type"`
}
type Dog struct {
	Name string `json:"name"`
	Type string `json:"Type"`
}
type Hamster struct {
	Name string `json:"name"`
	Type string `json:"Type"`
}
type JwtClaims struct{
	Name string `json:"name"`
	jwt.StandardClaims
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from server!!")
}
func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	return c.String(http.StatusOK, fmt.Sprintf("your cat name is %s.\n and type of cat is %s.", catName, catType))
}
func getById(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	datatype := c.Param("id")

	if datatype == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}
	if datatype == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is %s.\n and type of cat is %s.", catName, catType))
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to let us know if you want string or json datatype",
	})
}
func addCat(c echo.Context) error { //faster
	cat := Cat{}
	defer c.Request().Body.Close()
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body : %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(body, &cat)
	if err != nil {
		log.Printf("Failed Unmarshaling : %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your cat :%v", cat)
	return c.String(http.StatusOK, "we got your cat!!")
}
func addDog(c echo.Context) error {
	dog := Dog{}
	defer c.Request().Body.Close()
	// body,err:=ioutil.ReadAll(c.Request().Body)
	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Failed processing addDog request : %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	log.Printf("this is your cat :%v", dog)
	return c.String(http.StatusOK, "we got your dog!!")
}
func addHamster(c echo.Context) error { //slower
	hamster := Hamster{}
	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("Failed processing addHamster request : %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}
	log.Printf("this is your cat :%v", hamster)
	return c.String(http.StatusOK, "we got your hamster!!")
}
func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "welcome to the secret admin page.")
}
func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "you are on the yet secret page.")
}
func login(c echo.Context) error{
	username:=c.QueryParam("username")
	password:=c.QueryParam("password")
	//check username and password against DB after hashing password
	if username=="rohan" && password=="1234"{
		cookie:=&http.Cookie{}
		//this is same
		//cookie:=new(http.Cookie)
		cookie.Name="sessionID"
		cookie.Value="some_string"
		cookie.Expires=time.Now().Add(48*time.Hour)
		c.SetCookie(cookie)

		//TODO create jwt Token
		token,err:=createJwtToken()
		if err!=nil{
			log.Println("Error generating JWT Token")
			return c.String(http.StatusInternalServerError,"Something is wrong")
		}
		return c.JSON(http.StatusOK,map[string]string{
			"message":"you are logged in!",
			"token":token,
		})
	}
	return c.String(http.StatusUnauthorized,"your username or password were wrong!")
}
func createJwtToken() (string,error){
	claims:=JwtClaims{
		"Jack",
		jwt.StandardClaims{
			Id:"main_user_id",
			ExpiresAt: time.Now().Add(48*time.Hour).Unix(),
		},
	}
	rawToken:=jwt.NewWithClaims(jwt.SigningMethodHS512,claims)
	token,err:=rawToken.SignedString([]byte("mySecret"))
	if err!=nil{
		return "",err
	}
	return token,nil
}


//                                 --------Middleware--------------



func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc{
	return func(c echo.Context)error{
		c.Response().Header().Set(echo.HeaderServer,"BlueBot/1.0")
		c.Response().Header().Set("useless","MeaningNothing")
		return next(c)
	}
}
func checkCookie(next echo.HandlerFunc) echo.HandlerFunc{
	return func(c echo.Context)error{
		cookie,err:=c.Cookie("sessionID")
		if err!=nil{
			log.Println(err)
			return err
		}
		if cookie.Value=="some_string" { 
			return next(c)
		}
		return c.String(http.StatusUnauthorized,"you don't have right cookie")
	}
}

func main() {
	fmt.Println("Welcome to server!!")
	e := echo.New()
	// e.Use(ServerHeader)
	adminGroup := e.Group("/admin")
	cookieGroup:=e.Group("/cookie")
	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool,error) {
		//check in th DB
		if username == "rohan" && password == "1234" {
			return true,nil
		}
		return false,nil
	}))
	adminGroup.GET("/main", mainAdmin)
	cookieGroup.GET("/main",mainCookie)
	cookieGroup.Use(checkCookie)
	e.GET("/login",login)
	e.GET("/", hello)
	e.GET("/cats", getCats)
	e.GET("/cats/:id", getById)
	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	//Echo Method
	e.POST("/hamsters", addHamster)
	e.Start(":8000")
}
