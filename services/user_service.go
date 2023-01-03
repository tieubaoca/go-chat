package services

import "log"

func FindOnlineUsers() ([]string, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	users := make([]string, 0)
	for username, cs := range wsClients {
		if len(cs) > 0 {
			users = append(users, username)
		}
	}
	return users, nil
}
