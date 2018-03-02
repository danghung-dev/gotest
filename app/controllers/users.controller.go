package controllers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	//"gotest/app/utils"
	"gotest/database"
	"gotest/app/models"
	"encoding/json"
)

type UserController struct {
	db database.MySQLDB
}

func New() *UserController {
	return &UserController{}
}

func (uc *UserController) HelloUser(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
    "omg":    true,
    "number": 122,
  }).Warn("Cuc log helloUser")
}

func (uc *UserController) Create(w http.ResponseWriter, r *http.Request) {
	// Validate the length of the body since some controllers could send a big payload
	/*required := []string{"name", "email", "password"}
	if len(params) != len(required) {
		err := NewAPIError(false, "Invalid request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}*/
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{"r.body": &user}).Info("test1323")

	//j, _ := GetJSON(r.Body)
	defer r.Body.Close()
	//if err != nil {
	//	NewAPIError(&APIError{false, "Invalid request", http.StatusBadRequest}, w)
	//	return
	//}
	//
	//name, err := j.GetString("name")
	//if err != nil {
	//	NewAPIError(&APIError{false, "Name is required", http.StatusBadRequest}, w)
	//	return
	//}
	//// TODO: Implement something like this and embed in a basecontroller https://stackoverflow.com/a/23960293/2554631
	//if len(name) < 2 || len(name) > 32 {
	//	NewAPIError(&APIError{false, "Name must be between 2 and 32 characters", http.StatusBadRequest}, w)
	//	return
	//}
	//
	//email, err := j.GetString("email")
	//if err != nil {
	//	NewAPIError(&APIError{false, "Email is required", http.StatusBadRequest}, w)
	//	return
	//}
	//if ok := utils.IsEmail(email); !ok {
	//	NewAPIError(&APIError{false, "You must provide a valid email address", http.StatusBadRequest}, w)
	//	return
	//}
	//exists := uc.UserRepository.Exists(email)
	//if exists {
	//	NewAPIError(&APIError{false, "The email address is already in use", http.StatusBadRequest}, w)
	//	return
	//}
	//pw, err := j.GetString("password")
	//if err != nil {
	//	NewAPIError(&APIError{false, "Password is required", http.StatusBadRequest}, w)
	//	return
	//}
	//if len(pw) < 6 {
	//	NewAPIError(&APIError{false, "Password must not be less than 6 characters", http.StatusBadRequest}, w)
	//	return
	//}
	//
	//u := &models.User{
	//	Name:      name,
	//	Email:     email,
	//	Admin:     false,
	//	CreatedAt: time.Now(),
	//}
	//u.SetPassword(pw)
	//
	//err = uc.UserRepository.Create(u)
	//if err != nil {
	//	NewAPIError(&APIError{false, "Could not create user", http.StatusBadRequest}, w)
	//	return
	//}
	//
	//defer r.Body.Close()
	NewAPIResponse(&APIResponse{Success: true, Message: "User created"}, w, http.StatusOK)
}
