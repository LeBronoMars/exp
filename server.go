package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"exp/mngr/api/config"
	m "exp/mngr/api/models"
	h "exp/mngr/api/handlers"
	"github.com/jinzhu/gorm"
	"github.com/dgrijalva/jwt-go"
	"github.com/itsjamie/gin-cors"
	_ "github.com/go-sql-driver/mysql"

	//"github.com/robfig/cron"
	//"github.com/pusher/pusher-http-go"
)

func main() {
	db := *InitDB()
	router := gin.Default()
	config := cors.Config{
		Origins:         "*",
		RequestHeaders:  "Authorization",
		Methods:         "GET, POST, PUT",
		Credentials:     true,
		ValidateHeaders: false,
		MaxAge:          24 * time.Hour,
	}
	router.Use(cors.Middleware(config))

	//readPollutants()
	//readStationCSV(&db)
	//readFTP("/WEBSITE", &db)
	//readFTPPerMinute("/WEBSITE", &db)
	LoadAPIRoutes(router, &db)
}

func LoadAPIRoutes(r *gin.Engine, db *gorm.DB) {
	//sher := *InitPusher()
	public := r.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		"admin@condobills.com" : "P@ssw0rd",
		}))
	
	private := r.Group("/api/v1")
	private.Use(Auth(config.GetString("TOKEN_KEY")))

	//manage users
	userHandler := h.NewUserHandler(db)
	public.POST("/register", userHandler.Create)
	public.POST("/login", userHandler.Auth)
	private.GET("/users", userHandler.Index)
	private.PUT("/change_password", userHandler.ChangePassword)
	private.PUT("/change_profile_pic", userHandler.ChangeProfilePic)
	public.POST("/forgot_password", userHandler.ForgotPassword)
	private.GET("/me", userHandler.GetUserInfo)

	//manage group
	groupHandler := h.NewGroupHandler(db)
	private.GET("/groups", groupHandler.Index)
	private.POST("/group", groupHandler.Create)
	private.GET("/groups/search/:group_name", groupHandler.Search)

	//manage group members
	groupMemberHandler := h.NewGroupMemberHandler(db)
	private.GET("/members", groupMemberHandler.Index)
	private.POST("/add_member", groupMemberHandler.Create)
	private.POST("/remove_member", groupMemberHandler.Delete)

	var port = os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	r.Run(fmt.Sprintf(":%s", port))
}

func InitDB() *gorm.DB {
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.GetString("DB_USER"), config.GetString("DB_PASS"),
		config.GetString("DB_HOST"), config.GetString("DB_PORT"),
		config.GetString("DB_NAME"))
	log.Printf("\nDatabase URL: %s\n", dbURL)

	_db, err := gorm.Open("mysql", dbURL)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database:  %s", err))
	}
	_db.DB()
	_db.LogMode(true)
	_db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&m.User{},
																&m.Group{},
																&m.GroupMember{})
	log.Printf("must create tables")
	return _db
}


func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != "" {
			tokenString := c.Request.Header.Get("Authorization")

			if (strings.Contains(tokenString, "Bearer")) {
				token, err := jwt.Parse(tokenString[7 : len(tokenString)], func(token *jwt.Token) (interface{}, error) {
				    return []byte(secret), nil
				})
				if err != nil || !token.Valid {
					response := &Response{Message: err.Error()}
					c.JSON(http.StatusUnauthorized, response)
					c.Abort()
				} 
			} else {
				response := &Response{Message: "Invalid token!"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
			}
		} else {
			response := &Response{Message: "Authorization is required"}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
		}
	}
}

// func InitPusher() *pusher.Client {
//     client := pusher.Client{
//       AppId: "247359",
//       Key: "59834d0368c9d87da49c",
//       Secret: "32e5ada51483143bd015",
//       Cluster: "ap1",
//     }
//     return &client
// }

type Response struct {
	Message string `json:"message"`
}