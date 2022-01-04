package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)
type Cat struct{
	Name string `json:"name"`
	Type string `json:"Type"`
}
type Dog struct{
	Name string `json:"name"`
	Type string `json:"Type"`
}
func hello(c echo.Context)error{
	return c.String(http.StatusOK,"Hello from server!!")
}
func getCats(c echo.Context) error{
	catName:=c.QueryParam("name")
	catType:=c.QueryParam("type")

	return c.String(http.StatusOK,fmt.Sprintf("your cat name is %s.\n and type of cat is %s.",catName,catType))
}
func getById(c echo.Context) error{
	catName:=c.QueryParam("name")
	catType:=c.QueryParam("type")
	datatype:=c.Param("id")

	if datatype=="json"{
		return c.JSON(http.StatusOK,map[string]string{
			"name":catName,
			"type":catType,
		})
	}
	if datatype=="string"{
		return c.String(http.StatusOK,fmt.Sprintf("your cat name is %s.\n and type of cat is %s.",catName,catType))
	}
	return c.JSON(http.StatusBadRequest,map[string]string{
		"error":"you need to let us know if you want string or json datatype",
	})
}
func addCat(c echo.Context) error{
	cat:=Cat{}
	defer c.Request().Body.Close()
	body,err:=ioutil.ReadAll(c.Request().Body)
	if err!=nil{
		log.Printf("Failed reading the request body : %s",err)
		return c.String(http.StatusInternalServerError,"")
	}
	err=json.Unmarshal(body,&cat)
	if err!=nil{
		log.Printf("Failed Unmarshaling : %s",err)
		return c.String(http.StatusInternalServerError,"")
	}
	log.Printf("this is your cat :%v",cat)
	return c.String(http.StatusOK,"we got your cat!!")
}
func addDog(c echo.Context)error{
	dog:=Dog{}
	defer c.Request().Body.Close()
	body,err:=ioutil.ReadAll(c.Request().Body)
}

func main() {
	fmt.Println("Welcome to server!!")
	e:=echo.New()
	e.GET("/",hello)
	e.GET("/cats",getCats)
	e.GET("/cats/:id",getById)
	e.POST("/cats",addCat)
	e.POST("/dogs",addDog)
	e.Start(":8000")
}

