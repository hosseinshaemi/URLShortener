package main

import (
	controllers "URLShortner/Controllers"
	repository "URLShortner/Repository"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	repository.DbConn = repository.SQLiteContext{}
	repository.DbConn.Init()
	db, _ := repository.DbConn.Open()

	db.Exec("CREATE TABLE IF NOT EXISTS users (UserId INTEGER PRIMARY KEY AUTOINCREMENT, Firstname VARCHAR(255) NOT NULL, Lastname VARCHAR(255) NOT NULL, Email VARCHAR(255) NOT NULL UNIQUE);")
	db.Exec("CREATE TABLE IF NOT EXISTS links (LinkId INTEGER PRIMARY KEY AUTOINCREMENT, UserId INTEGER, Name VARCHAR(255) NOT NULL, URL VARCHAR(255) NOT NULL, ShortLink VARCHAR(255) NOT NULL, FOREIGN KEY (UserId) REFERENCES users(UserId) ON DELETE CASCADE ON UPDATE NO ACTION);")

	db.Close()

	r := gin.Default()
	r.GET("/GetUsers", controllers.GetUsers)
	r.GET("/GetLinks", controllers.GetLinks)
	r.GET("/GetUserById/:userId", controllers.GetUserById)
	r.GET("/GetLinkByUser/:userId")
	r.POST("/AddUser", controllers.AddUser)
	r.POST("/AddLink", controllers.AddLink)
	r.GET("/short/:shorted", controllers.RedirectShortLink)
	r.Run()
}
