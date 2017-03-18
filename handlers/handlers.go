package handlers

import (
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"os"
)

type Handler struct {
	DB redis.Conn
}

type User struct {
	Id        string `json:"_id" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Gender    string `json:"gender" bson:"gender"`
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
}

/*-----  End of connect  ----*/

/*=================================
***   set value in database  ***
=================================*/

func (h *Handler) Set(ctx *gin.Context) {
	user := User{}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"Error": "Validation error"})
	} else {
		user.Id = uuid.NewV1().String()
		log.Println("user", user)
		_, err := redis.String(h.DB.Do("HSET", "user:"+user.Id, user))
		if err != nil {
			ctx.JSON(500, gin.H{"error": err})
		} else {
			ctx.JSON(200, gin.H{"All": "Good"})
		}
	}
}

/*-----  End of Set  ----*/

/*===================================
***   get value from database  ***
===================================*/

func (h *Handler) Get(ctx *gin.Context) {
	value, _ := ctx.GetQuery("key")
	if value == "" {
		ctx.JSON(400, gin.H{"Error": "Wrong query string"})
	} else {
		_, err := redis.String(h.DB.Do("HGETALL", "user:"+value))
		if err != nil {
			ctx.JSON(500, gin.H{"error": err})
		} else {
			ctx.JSON(200, gin.H{"All": "Good"})
		}
	}
}

/*-----  End of Get  ------*/

/*===================================
***   update value in database  ***
===================================*/

func (h *Handler) Update(ctx *gin.Context) {
	// value, _ := ctx.GetQuery("key")
	// if value == "" {
	// 	ctx.JSON(400, gin.H{"Error": "Wrong query string"})
	// } else {
	// 	reply, err := redis.String(h.DB.Do("get", value))
	// 	if err != nil {
	// 		ctx.JSON(500, gin.H{"error": err})
	// 	} else {
	ctx.JSON(200, gin.H{"All": "Good"})
	// 	}
	// }
}

/*-----  End of Update  ------*/

/*===================================
***   delete value from database  ***
===================================*/

func (h *Handler) Delete(ctx *gin.Context) {
	// value, _ := ctx.GetQuery("key")
	// if value == "" {
	// 	ctx.JSON(400, gin.H{"Error": "Wrong query string"})
	// } else {
	// 	reply, err := redis.String(h.DB.Do("get", value))
	// 	if err != nil {
	// 		ctx.JSON(500, gin.H{"error": err})
	// 	} else {
	ctx.JSON(200, gin.H{"All": "Good"})
	// 	}
	// }
}

/*-----  End of Delete  ------*/
