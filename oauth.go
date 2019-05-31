package main

import (
	//	"bytes"

	//"encoding/json"
	"fmt"
	"log"
	"time"

	//	"github.com/GoogleIdTokenVerifier/GoogleIdTokenVerifier"
	"github.com/dgrijalva/jwt-go"
	//	fb "github.com/huandu/facebook"
	//	"github.com/labstack/echo"
	// "github.com/labstack/echo/v4"
	// "github.com/labstack/echo/middleware"
	fb "github.com/huandu/facebook"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createAccessToken(accessKey string, id string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	exp := time.Now().Add(time.Hour * 72).Unix()
	claims["id"] = id
	claims["type"] = "access"
	claims["exp"] = exp
	t, err := token.SignedString([]byte(accessKey + "___Good###" + id))
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func createRefreshToken(accessKey string, id string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	exp := time.Now().Add(time.Hour * 720).Unix()
	claims["type"] = "refresh"
	claims["id"] = id
	claims["exp"] = exp
	t, err := token.SignedString([]byte(accessKey + "____bad###" + id))
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func getIdAndIsVaildToken(tokenStr string) primitive.ObjectID {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		var id string = token.Claims.(jwt.MapClaims)["id"].(string)
		accessKeys := accessKey + "___Good###" + id
		return []byte(accessKeys), nil
	})

	if err != nil {
		return primitive.NilObjectID
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["type"] == "access" {
		fmt.Println(claims, "   ::: CLiammmmmmmmms true  ::::", ok, token.Valid)
		//		idtext := strings.TrimSuffix(strings.TrimPrefix(claims["id"].(string), "ObjectID(\""), "\")")
		idc, err := primitive.ObjectIDFromHex(claims["id"].(string))
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

func faceVerify(token string) *Usertokeninfo {
	debugToken, err := fb.Get("/debug_token", fb.Params{ // to get info about token to verify it
		"input_token":  token,
		"access_token": id + "|" + secid,
	})
	if err != nil {
		return nil
	}
	var resultdebug struct {
		Data debuginfo `json:"data"`
	}
	debugToken.Decode(&resultdebug)
	if resultdebug.Data.AppId == id && resultdebug.Data.IsValid {
		userToken, _ := fb.Get("/me", fb.Params{
			"fields":       "id,name,email,picture{url}",
			"access_token": token,
		})
		var resultInfo *Usertokeninfo
		userToken.Decode(&resultInfo)
		return resultInfo
	}
	return nil
}

func createToken(id string) authres {
	accessToken := createAccessToken(accessKey, id)
	refreshJWT := createRefreshToken(refreshKey, id)
	return authres{Refresh: refreshJWT, Access: accessToken}
}

func IsVaildRefreshToken(tokenStr string) primitive.ObjectID {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		var id string = token.Claims.(jwt.MapClaims)["id"].(string)
		keys := refreshKey + "____bad###" + id
		return []byte(keys), nil
	})

	if err != nil {
		return primitive.NilObjectID
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["type"] == "refresh" {
		fmt.Println(claims, "   ::: CLiammmmmmmmms true  ::::", ok, token.Valid)
		//		idtext := strings.TrimSuffix(strings.TrimPrefix(claims["id"].(string), "ObjectID(\""), "\")")
		idc, err := primitive.ObjectIDFromHex(claims["id"].(string))
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
