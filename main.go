package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/GoogleIdTokenVerifier/GoogleIdTokenVerifier"

	//	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//github.com\mongodb\mongo-go-driver
)

var accessKey string = "SecretaccessKeyForJWT#$*(FJIEOF$FEFK#"
var refreshKey string = "Secret#@DFWrefreshKeyForJWT#$*(FJIEOF$FEFK#"
var myaud string = "421652019678-c76vldjrurop7m3thl75msi805hdqcrb.apps.googleusercontent.com"
var id string = "1646756352305359"
var secid string = "52604b40a547e38ed33754a6e5ccaa5a"

type debuginfo struct {
	AppId   string `json:"app_id"`
	IsValid bool   `json:"is_valid"`
	UserId  string `json:"user_id"`
}
type Usertokeninfo struct {
	UserId  string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture struct {
		Data struct {
			Url string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

type dataReq struct {
	Word  string `json:"word" form:"word" query:"word"`
	Learn int    `json:"learn" form:"learn" query:"learn"`
	Qtype []int  `json:"Qtype" form:"Qtype" query:"Qtype"`
}
type jsonreq struct {
	Token string    `json:"token" form:"token" query:"token"`
	Data  []dataReq `json:"data" form:"data" query:"data"`
}
type authres struct {
	Refresh string `json:"refresh" form:"refresh" query:"refresh"`
	Access  string `json:"access" form:"access" query:"access"`
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

func version1(e *echo.Echo) {
	v1 := e.Group("/api/v1")
	//	v1.POST("/user/login", loginUser)
	//	v1.POST("/user/signin", signinUser)
	v1.POST("/word/add", addWord)
	v1.POST("/word/all", allWord)
	v1.POST("/word/sync", syncword)
	v1.POST("/word/learn", learnword)
	v1.GET("/oauth/google", oauthGetGoogle)
	v1.GET("/oauth/facebook", oauthGetFacebook)
	v1.GET("/oauth/refresh", refreshToken)
}

func refreshToken(c echo.Context) error {
	token := c.QueryParam("token")
	id := IsVaildRefreshToken(token)
	if id == primitive.NilObjectID {
		return c.String(200, "exp resfresh token or invaild")
	}
	newAccess := createAccessToken(accessKey, id.Hex())
	return c.String(400, newAccess)
}

func oauthGetFacebook(c echo.Context) error {
	fmt.Println("handle Correct GOOOOOOOOOD Facebook")
	token := c.QueryParam("token")
	if resultInfo := faceVerify(token); resultInfo != nil {
		/////token verify corrcet
		if isexist, result := findBysubFacbookId(resultInfo.UserId); isexist {
			tokenObject := createToken(result.ID.Hex())
			return c.JSON(201, tokenObject)
		} else {
			user := UserInfoset{Name: resultInfo.Name, Email: resultInfo.Email, PhotoUrl: resultInfo.Picture.Data.Url, FacebookId: resultInfo.UserId}
			id := user.insertUser()
			tokenObject := createToken(id.Hex())
			return c.JSON(201, tokenObject)
		}
	} else {
		return c.String(401, "Invalid JWT Token or exp")
	}
}

func oauthGetGoogle(c echo.Context) error {
	fmt.Println("handle Correct GOOOOOOOOOD")
	token := c.QueryParam("token")
	tokenInfo := GoogleIdTokenVerifier.Verify(token, myaud)
	if tokenInfo == nil {
		fmt.Println("nil")
		return c.String(401, "Invalid JWT Token or exp")
	}
	if isexist, result := findBysubGoogleId(tokenInfo.Sub); isexist {
		tokenObject := createToken(result.ID.Hex())
		fmt.Println(tokenObject.Refresh, "   2 ", tokenObject.Access)
		fmt.Println(tokenObject)
		return c.JSON(201, tokenObject)
	} else {
		user := UserInfoset{Name: tokenInfo.Name, Email: tokenInfo.Email, PhotoUrl: tokenInfo.Picture, Local: tokenInfo.Local, SubGoogleId: tokenInfo.Sub}
		id := user.insertUser()
		tokenObject := createToken(id.Hex())
		fmt.Println(tokenObject)
		fmt.Println(tokenObject.Refresh, "   2 ", tokenObject.Access)
		return c.JSON(201, tokenObject)
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
		return c.String(400, err.Error())
	}
	return c.String(200, "Done")
}

// func hashPAssfunc(pass string, accessKey [16]byte, id primitive.ObjectID) []byte {
// 	fmt.Println(pass, " :: ", accessKey, " :: ", id.String())

// 	return pbkdf2.accessKey([]byte(pass), append([]byte(id.String()), accessKey[:]...), 10, 64, sha256.New)

// }

// func signinUser(c echo.Context) error {
// 	email := c.FormValue("email")
// 	if IsEmailExist(email) {
// 		return c.String(409, "this Email is already used")
// 	} else {
// 		user := UserInfoset{Name: c.FormValue("name"), Email: email}
// 		id := user.insertUser()
// 		token := createAccessToken(accessKey, id.String())
// 		uuidHash, _ := uuid.New()
// 		hashPass := hashPAssfunc(c.FormValue("pass"), uuidHash, id)
// 		addHashPass(hashPass, uuidHash, id)
// 		//	return c.JSON(201, map[string]string{"messege": "sign in suecceful", "token": token})
// 		return c.JSON(201, token)
// 	}
// }

// func validateEmail(email string) bool {
// 	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
// 	return Re.MatchString(email)
// }
