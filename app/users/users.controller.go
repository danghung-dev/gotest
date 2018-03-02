package users

import (
	log "github.com/sirupsen/logrus"
	"net/http"
  )

type UserController struct {

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