// package main

// import (
// 	"fmt"

// 	fb "github.com/huandu/facebook"
// )

// var accesstoken string = "EAAXZAt1dpKM8BAMoXkJKUChk36sWw9I0ymogN2DHJzO8jwzZAVSY3xMiBub2p6zEeTPuJTW1K7ZBAEz9ZAnZAHNNByW2LOA6CHO274KTqgk7AZAxXL5Aq3NGZBZBc83JYxNhuI26h1SGsnMRMhs2ltt6rEDz7enpfxEZD"
// var id string = "1646756352305359"
// var secid string = "52604b40a547e38ed33754a6e5ccaa5a"

// // func debugFacebookToken(tokeninput, tokeAPP) string {
// // 	return "graph.facebook.com/debug_token?input_token=" + tokeninput + "&access_token=" + id + "|" + secid
// // }

// // func getinfoFacebookToken(tokeninput) string {
// // 	return "https://graph.facebook.com/me?fields=id,name,email&access_token=" + tokeninput
// // }
// type debuginfo struct {
// 	AppId   string `json:"app_id"`
// 	IsValid bool   `json:"is_valid"`
// 	UserId  string `json:"user_id"`
// }
// type Usertokeninfo struct {
// 	UserId  string `json:"id"`
// 	Name    string `json:"name"`
// 	Email   string `json:"email"`
// 	Picture struct {
// 		Data struct {
// 			Url string `json:"url"`
// 		} `json:"data"`
// 	} `json:"picture"`
// }

// // type pictureUrl struct {

// // }

// func main() {
// 	debugToken, _ := fb.Get("/debug_token", fb.Params{
// 		"input_token":  accesstoken,
// 		"access_token": id + "|" + secid,
// 	})
// 	var resultdebug struct {
// 		Data debuginfo `json:"data"`
// 	}
// 	debugToken.Decode(&resultdebug)
// 	fmt.Println(resultdebug.Data.UserId)
// 	if resultdebug.Data.AppId == id {
// 		if resultdebug.Data.IsValid {
// 			userToken, _ := fb.Get("/me", fb.Params{
// 				"fields":       "id,name,email,picture{url}",
// 				"access_token": accesstoken,
// 			})
// 			var resultInfo Usertokeninfo
// 			userToken.Decode(&resultInfo)
// 			fmt.Println(userToken)
// 			fmt.Println("resultInfo.email")
// 			fmt.Println(resultInfo)
// 			fmt.Println(resultInfo.Picture.Data.Url)

// 		} else {

// 		}

// 	} else {

// 	}

// 	res, _ := fb.Get("/me", fb.Params{
// 		"fields":       "id,name,email",
// 		"access_token": accesstoken,
// 	})
// 	fmt.Println("Here is my Facebook first name:", res)
// 	fmt.Println(res.UsageInfo())
// 	fmt.Println(res.GetField())
// 	fmt.Println(res.DebugInfo())
// 	fmt.Println("_______________________")
// }
