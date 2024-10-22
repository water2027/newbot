package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"testsse/config"
)

func LoginSSEReq(config *config.BotConfig) (*http.Request, error) {
	//login
	loginData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, config.Email, config.Password)
	loginReq, err := http.NewRequest("POST", "https://ssemarket.cn/api/auth/login", bytes.NewBuffer([]byte(loginData)))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	loginReq.Header.Set("Content-Type", "application/json")
	return loginReq, nil
}

func GetPostsReq(config *config.BotConfig) (*http.Request, error) {
	//get posts
	getPostsData := fmt.Sprintf(`{"limit":5,"offset":0,"partition":"主页","searchsort":"home","userTelephone":"%s"}`, config.Telephone)
	req, err := http.NewRequest("POST", "https://ssemarket.cn/api/auth/browse", bytes.NewBuffer([]byte(getPostsData)))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func GetPostContentReq(id int, config *config.BotConfig) (*http.Request, error) {
	//get post content
	phonenum, err := strconv.Atoi(config.Telephone)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	str := fmt.Sprintf(`{"userTelephone": "%d","postID": %d}`, phonenum, id)
	req, err := http.NewRequest("POST", "https://ssemarket.cn/api/auth/showDetails", bytes.NewBuffer([]byte(str)))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	return req, nil
}

func GetHeatPostsReq(config *config.BotConfig) (*http.Request, error) {
	//get heat post
	req, err := http.NewRequest("GET", "https://ssemarket.cn/api/auth/calculateHeat", nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
