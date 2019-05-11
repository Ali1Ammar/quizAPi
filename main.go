package main

import (
	//"go/token"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"golang.org/x/crypto/pbkdf2"
)

var key string = "SecretKeyForJWT#$*(FJIEOF$FEFK#"

type dataReq struct {
	Word  string `json:"word" form:"word" query:"word"`
	Learn int    `json:"learn" form:"learn" query:"learn"`
	Qtype []int  `json:"Qtype" form:"Qtype" query:"Qtype"`
}
type jsonreq struct {
	Token string    `json:"token" form:"token" query:"token"`
	Data  []dataReq `json:"data" form:"data" query:"data"`
}

func main() {

	e := echo.New()
	e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	version1(e)
	e.Logger.Fatal(e.Start(":1323"))
}
func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}
func version1(e *echo.Echo) {
	v1 := e.Group("/api/v1")
	v1.POST("/user/login", loginUser)
	v1.POST("/user/signin", signinUser)
	v1.POST("/word/add", addWord)
	v1.POST("/word/all", allWord)
	v1.POST("/word/sync", syncword)
	v1.POST("/word/learn", learnword)

}
func loginUser(c echo.Context) error {
	email := c.FormValue("email")
	if !validateEmail(email) {
		return c.String(401, "Email or Password is Not correct")
	}
	fmt.Println("Email : ", email, "...")
	emailExist, result := findByEmail(email)

	fmt.Println("Email : ", email, "...")
	newPass := hashPAssfunc(c.FormValue("pass"), result.UuidHash, result.ID)
	fmt.Println(result.Pass, "  ::: ", newPass, " ==  ", bytes.Equal(newPass, result.Pass))
	if emailExist && bytes.Equal(newPass, result.Pass) {
		newJWT := createToken(key, result.ID.String())
		return c.JSON(202, newJWT)
	} else {
		return c.String(401, "Email or Password is Not correct")
	}
}
func createToken(key string, id string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	exp := time.Now().Add(time.Hour * 72).Unix()
	claims["id"] = id
	claims["exp"] = exp
	t, err := token.SignedString([]byte(key + "___Good###" + id))
	if err != nil {
		log.Fatal(err)
	}
	return t
}
func signinUser(c echo.Context) error {
	email := c.FormValue("email")
	if IsEmailExist(email) {
		return c.String(409, "this Email is already used")
	} else {
		user := UserInfoset{Name: c.FormValue("name"), Email: email}
		id := user.insertUser()
		token := createToken(key, id.String())
		uuidHash, _ := uuid.New()
		hashPass := hashPAssfunc(c.FormValue("pass"), uuidHash, id)
		addHashPass(hashPass, uuidHash, id)
		//	return c.JSON(201, map[string]string{"messege": "sign in suecceful", "token": token})
		return c.JSON(201, token)
	}
}
func addWord(c echo.Context) error {
	token := c.FormValue("token")
	id := getIdAndIsVaildToken(token)
	if id == primitive.NilObjectID {
		return c.String(401, "Invalid JWT Token or exp")
	} else {
		words := c.FormValue("words")
		var wordsStruct []WordInfo
		err := json.Unmarshal([]byte(words), &wordsStruct)
		if err != nil {
			log.Fatal(err)
		}
		//	fmt.Println(wordsStruct[0].date, "   ", wordsStruct[0].date.IsZero(), "   ", time.Now())
		// if wordsStruct[0].date.IsZero() {
		// 	for index := range wordsStruct {
		// 		wordsStruct[index].date = time.Now()
		// 		fmt.Println(wordsStruct[index].date)
		// 	}
		// }
		fmt.Println("it is date : ", wordsStruct)
		errr := insertManyWord(wordsStruct, id)
		if errr != nil {
			return c.String(400, "error for formatted words")
		}
		return c.String(200, "correct done")
	}
}
func getIdAndIsVaildToken(tokenStr string) primitive.ObjectID {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		var id string = token.Claims.(jwt.MapClaims)["id"].(string)
		keys := key + "___Good###" + id
		return []byte(keys), nil
	})

	if err != nil {
		return primitive.NilObjectID
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims, "   ::: CLiammmmmmmmms true  ::::", ok, token.Valid)
		idtext := strings.TrimSuffix(strings.TrimPrefix(claims["id"].(string), "ObjectID(\""), "\")")
		idc, err := primitive.ObjectIDFromHex(idtext)
		fmt.Println(idtext, "  ::  ", idc)
		if err != nil {
			fmt.Println(err)
		}
		return idc
	} else {
		fmt.Println(claims, "   ::: CLiammmmmmmmms  ::::", ok, token.Valid)
		log.Printf("Invalid JWT Token")
		return primitive.NilObjectID
	}
}
func hashPAssfunc(pass string, key [16]byte, id primitive.ObjectID) []byte {
	fmt.Println(pass, " :: ", key, " :: ", id.String())

	return pbkdf2.Key([]byte(pass), append([]byte(id.String()), key[:]...), 10, 64, sha256.New)

}

func allWord(c echo.Context) error {
	token := c.FormValue("token")
	id := getIdAndIsVaildToken(token)
	if id == primitive.NilObjectID {
		return c.String(401, "Invalid JWT Token or exp")
	} else {
		mywords := getAllWords(id)
		if mywords.Words == nil {
			return c.String(204, "no contect")
		}

		return c.JSON(200, mywords)
	}
}

func syncword(c echo.Context) error {
	token := c.FormValue("token")
	id := getIdAndIsVaildToken(token)
	if id == primitive.NilObjectID {
		return c.String(401, "Invalid JWT Token or exp")
	}
	datesync := c.FormValue("datesync")
	timeDate, err := time.Parse(time.RFC3339, datesync)
	if err != nil {
		return c.String(400, "date format wrong")
	}
	newword := syncDate(id, timeDate)
	if len(newword.Words) == 0 {
		return c.String(204, "every things Ok no new updates")
	}
	fmt.Println("HERE LEN not 0")
	return c.JSON(200, newword)
}

func learnword(c echo.Context) error {

	var jsoneReq jsonreq
	if err := c.Bind(&jsoneReq); err != nil {
		return err
	}
	fmt.Println(jsoneReq)
	token := jsoneReq.Token
	id := getIdAndIsVaildToken(token)
	if id == primitive.NilObjectID {
		return c.String(401, "Invalid JWT Token or        expccc")
	}
	fmt.Println(token, id)

	if err := learn(id, jsoneReq.Data); err != nil {
		return err
	}
	return c.String(http.StatusOK, fmt.Sprintf("%v", jsoneReq.Data))
}
