package handlers

import (
	"net/http"
	"fmt"
	"io"
	"time"
	"crypto/aes"	
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"

	"exp/mngr/api/config"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/pusher/pusher-http-go"
    "github.com/satori/go.uuid"
)

func respond(statusCode int, responseMessage string, c *gin.Context, isError bool) {
	response := &Response{Message: responseMessage}
	c.JSON(statusCode,response)
	if (isError) {
		c.Abort()
	}
}

func jwtVerifier() gin.HandlerFunc {
	return func(c *gin.Context) {

		appToken := c.Request.Header.Get("Authorization")

		if appToken == "" {
			respond(http.StatusForbidden, "Authorization header is required", c, true)
		} else {
			respond(http.StatusBadRequest, fmt.Sprintf("Invalid token: %s", appToken), c, true)
		}
	}
}

type Response struct {
	Message string `json:"message"`
}

//generate JWT
func generateJWT(userId string) string {
	mySigningKey := []byte(config.GetString("TOKEN_KEY"))
    claims := &jwt.StandardClaims{
    	ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
    	Issuer:    userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString,_ := token.SignedString(mySigningKey)
    return tokenString
}

//change timezone of date
func changeTimeZone(t time.Time) time.Time {
	loc,_ := time.LoadLocation("Asia/Manila")
	newTime,_ := time.ParseInLocation(time.RFC3339,t.Format(time.RFC3339),loc)
	return newTime
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func SendPushNotification(pusher *pusher.Client, channel string, event string, message string) {
	data := map[string]string{"message": message}
	pusher.Trigger(channel,event,data)
}

func GetStartOfDay(t time.Time) time.Time {
    year, month, day := t.Date()
    return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func GetEndOfDay(t time.Time) time.Time {
    year, month, day := t.Date()
    return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}

func GenerateID() string {
	return fmt.Sprintf("%s", uuid.NewV4())
}

func GetCreator(c *gin.Context) string {
	tokenString := c.Request.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString[7 : len(tokenString)], func(token *jwt.Token) (interface{}, error) {
	    return []byte(config.GetString("TOKEN_KEY")), nil
	})
	if err == nil && token.Valid {
		claims, _ := token.Claims.(jwt.MapClaims)
		return fmt.Sprintf("%s", claims["iss"])
	} else {
		return ""
	}
}

func GetDeletedAt(c *gin.Context) *time.Time {
	layout := "2006-01-02T15:04:05Z"
	t, _ := time.Parse(layout, c.PostForm("deleted_at"))
	return &t
}