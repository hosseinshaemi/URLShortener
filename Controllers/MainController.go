package controllers

import (
	repository "URLShortner/Repository"
	models "URLShortner/Repository/Models"
	utils "URLShortner/Utils"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	db, err := repository.DbConn.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in connect to database"})
		return
	}
	userRepo := repository.NewUserRepository(db)
	users, _ := userRepo.List()
	db.Close()
	c.JSON(http.StatusOK, users)
}

func GetUserById(c *gin.Context) {
	db, err := repository.DbConn.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in connect to database"})
		return
	}

	idNum, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not valid"})
		db.Close()
		return
	}
	userRepo := repository.NewUserRepository(db)
	user, _ := userRepo.GetById(int64(idNum))

	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This is user does not exist"})
		db.Close()
		return
	}

	db.Close()
	c.JSON(http.StatusOK, user)
}

func GetLinks(c *gin.Context) {
	db, err := repository.DbConn.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in connect to database"})
		return
	}
	linkRepo := repository.NewLinkRepository(db)
	links, _ := linkRepo.List()
	db.Close()
	c.JSON(http.StatusOK, links)
}

func AddUser(c *gin.Context) {
	db, err := repository.DbConn.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in connect to database"})
		return
	}

	userRepo := repository.NewUserRepository(db)

	var firstname string = c.PostForm("firstname")
	var lastname string = c.PostForm("lastname")
	var email string = c.PostForm("email")

	var user models.User = models.User{Firstname: firstname, Lastname: lastname, Email: email}
	fmt.Println(user)
	er := userRepo.CreateUser(&user)
	if er != nil {
		if strings.Contains(er.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This user already exists"})
			db.Close()
			return
		}
	}
	db.Close()
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func AddLink(c *gin.Context) {
	db, err := repository.DbConn.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in connect to database"})
		return
	}

	linkRepo := repository.NewLinkRepository(db)

	userId, err := strconv.Atoi(c.PostForm("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		db.Close()
		return
	}
	var name string = c.PostForm("name")
	var url string = c.PostForm("url")
	var short string = "http://localhost:8080/short/" + utils.RandomStringGenerator()

	var link models.Link = models.Link{UserId: int64(userId), Name: name, Url: url, ShortLink: short}
	er := linkRepo.CreateLink(&link)
	if er != nil {
		if strings.Contains(er.Error(), "wrong UserId") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
			db.Close()
			return
		}
	}
	db.Close()
	c.JSON(http.StatusOK, gin.H{"message": "link successfully created", "link": link})
}

func GetLinkByUser(c *gin.Context) {
	db, err := repository.DbConn.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in connect to database"})
		return
	}

	idNum, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not valid"})
		db.Close()
		return
	}
	linkRepo := repository.NewLinkRepository(db)
	links, err := linkRepo.GetByUserId(idNum)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred during fetching data"})
		db.Close()
		return
	}
	db.Close()
	c.JSON(http.StatusOK, links)
}

func RedirectShortLink(c *gin.Context) {
	db, err := repository.DbConn.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in connect to database"})
		return
	}

	var shorted string = "http://localhost:8080/short/" + c.Param("shorted")
	linkRepo := repository.NewLinkRepository(db)
	link, err := linkRepo.GetLinkByShortLink(shorted)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occurred during fetching data"})
		db.Close()
		return
	}
	db.Close()
	c.Redirect(http.StatusMovedPermanently, link.Url)
}
