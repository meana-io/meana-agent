package users

import (
	"io/ioutil"
	"log"
	"strings"
)

type UsersData struct {
	Users []User `json:"users"`
}

type User struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
}

func GetUsersData() (*UsersData, error) {
	var usersData UsersData

	passwdFile, err := ioutil.ReadFile("/etc/passwd")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	groupFile, err := ioutil.ReadFile("/etc/group")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, usr := range strings.Split(string((passwdFile)), "\n") {
		split := strings.Split(usr, ":")
		if len(split) > 1 {
			if strings.Contains(split[5], "home") {
				var user User
				user.Username = split[0]

				for _, group := range strings.Split(string((groupFile)), "\n") {
					split2 := strings.Split(group, ":")
					if len(split2) > 1 {

						if strings.Contains(split2[3], split[0]) {
							user.Groups = append(user.Groups, split2[0])
						}
					}
				}
				usersData.Users = append(usersData.Users, user)
			}
		}
	}

	return &usersData, nil
}
