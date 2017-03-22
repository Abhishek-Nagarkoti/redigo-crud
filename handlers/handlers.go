package handlers

import (
	"github.com/Abhishek-Nagarkoti/redigo-crud/models"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"os"
)

type Handler struct {
	DB redis.Conn
}

/*=================================
***   establish connection  ***
=================================*/
func (h *Handler) Connect() {
	var err error
	h.DB, err = redis.Dial("tcp", os.Getenv("HOST")+":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}
	user := models.User{}
	user.Automigration(h.DB)
}

/*-----  End of connect  ----*/

/*=================================
***   set value in database  ***
=================================*/

func (h *Handler) Set(ctx *gin.Context) {
	user := models.User{}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"Error": "Validation error"})
	} else {
		user.Id = uuid.NewV1().String()
		err, userData := user.Create(h.DB)
		if err == nil {
			ctx.JSON(200, gin.H{"user": userData})
		} else {
			ctx.JSON(500, gin.H{"error": err})
		}
	}
}

/*-----  End of Set  ----*/

/*===================================
***   get value from database  ***
===================================*/

func (h *Handler) Get(ctx *gin.Context) {
	value, _ := ctx.GetQuery("id")
	if value == "" {
		user := models.User{}
		err, users := user.GetALL(h.DB)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err})
		} else {
			ctx.JSON(200, gin.H{"users": users})
		}
	} else {
		user := models.User{}
		err, userData := user.Get(h.DB, value)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err})
		} else {
			ctx.JSON(200, gin.H{"user": userData})
		}
	}
}

/*-----  End of Get  ------*/

/*===================================
***   update value in database  ***
===================================*/

func (h *Handler) Update(ctx *gin.Context) {
	user := models.User{}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"Error": "Validation error"})
	} else {
		user.Id = ctx.Param("id")
		err, userData := user.Update(h.DB)
		if err == nil {
			ctx.JSON(200, gin.H{"user": userData})
		} else {
			ctx.JSON(500, gin.H{"error": err})
		}
	}
}

/*-----  End of Update  ------*/

/*===================================
***   delete value from database  ***
===================================*/

func (h *Handler) Delete(ctx *gin.Context) {
	user := models.User{}
	user.Id = ctx.Param("id")
	err := user.Delete(h.DB)
	if err == nil {
		ctx.JSON(200, gin.H{"All": "Good"})
	} else {
		ctx.JSON(500, gin.H{"error": err})
	}
}

/*-----  End of Delete  ------*/

/*===================================
***   find by name from database  ***
===================================*/

func (h *Handler) Find(ctx *gin.Context) {
	user := models.User{}
	user.FirstName = ctx.Param("name")
	err, users := user.Find(h.DB)
	if err == nil {
		ctx.JSON(200, gin.H{"users": users})
	} else {
		ctx.JSON(500, gin.H{"error": err})
	}
}

/*-----  End of Find  ------*/
