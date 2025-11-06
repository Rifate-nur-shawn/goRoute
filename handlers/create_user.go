package handlers

import (
	"encoding/json"
	"fmt"
	"goRoute/database"
	"goRoute/util"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser database.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "invalid Request Data", http.StatusBadRequest)
	  	return
	}

	createdUser := newUser.Store()
	util.SendData(w,createdUser,201)

}
