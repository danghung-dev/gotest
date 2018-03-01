package users

import (
	log "github.com/sirupsen/logrus"
	"net/http"
  )

func HelloUser(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
    "omg":    true,
    "number": 122,
  }).Warn("Cuc log helloUser")
}