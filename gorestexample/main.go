package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"github.com/rs/cors"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Pics     []*Pic `json:"pics"`
}

type Pic struct {
	Image       string `json:"image"`
	Description string `json:"description"`
}

type PicRequest struct {
	UserId      string `json:"userId"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type LoginRequest struct {
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

type ServerResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type RegisterRequest struct {
	Name            string `json:"name"`
	Password        string `json:"pwd"`
	ConfirmPassword string `json:"confirmPwd"`
}

var Users []*User

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		resp, err := Login(r)
		var respBody ServerResponse
		respBody.Data = resp
		if err != nil {
			respBody.Message = err.Error()
		} else {
			respBody.Message = "Login Successful"
		}
		Responser(w, respBody, err != nil)
		return
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		resp, err := Register(r)
		var respBody ServerResponse
		respBody.Data = resp
		if err != nil {
			respBody.Message = err.Error()
		} else {
			respBody.Message = "Registration Successful"
		}
		Responser(w, respBody, err != nil)
		return
	})

	mux.HandleFunc("/postpic", func(w http.ResponseWriter, r *http.Request) {
		resp, err := AddPic(r)
		var respBody ServerResponse
		respBody.Data = resp
		if err != nil {
			respBody.Message = err.Error()
		} else {
			respBody.Message = "Successfully Added Picture"
		}
		Responser(w, respBody, err != nil)
		return
	})

    handler := cors.Default().Handler(mux)

	fmt.Println("I'm running at http://localhost:8081")
	http.ListenAndServe(":8081", handler)
}

func AddPic(r *http.Request) (*Pic, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var pic PicRequest

	err = json.Unmarshal(body, &pic)
	if err != nil {
		return nil, err
	}

	var savedUser *User
	var found bool = false
	for _, user := range Users {
		if user.ID == pic.UserId {
			savedUser = user
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("Please Login to add your pics")
	}

	if pic.Image == "" || pic.Description == "" {
		return nil, errors.New("Please provide both image url and description for your post")
	}
	newPic := &Pic{Image: pic.Image, Description: pic.Description}
	savedUser.Pics = append(savedUser.Pics, newPic)
	return newPic, nil
}

func Register(r *http.Request) (*User, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var reg RegisterRequest

	err = json.Unmarshal(body, &reg)
	if err != nil {
		return nil, err
	}

	if reg.Password != reg.ConfirmPassword {
		return nil, errors.New("Please check both password and confirm password")
	}

	for _, user := range Users {
		if user.Name == reg.Name {
			return nil, errors.New("User already exists")
		}
	}

	var savedUser *User = &User{
		ID:       fmt.Sprint(rand.Int()),
		Name:     reg.Name,
		Password: reg.Password,
	}

	Users = append(Users, savedUser)

	return savedUser, nil
}

func Login(r *http.Request) (*User, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var log LoginRequest

	err = json.Unmarshal(body, &log)
	if err != nil {
		return nil, err
	}

	var savedUser *User
	var found bool = false
	for _, user := range Users {
		if user.Name == log.Name {
			savedUser = user
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("Please register first")
	}

	if log.Pwd != savedUser.Password {
		return nil, errors.New("Wrong Password")
	}
	return savedUser, nil
}

func Responser(w http.ResponseWriter, content interface{}, isError bool) {
	if isError {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	body, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
