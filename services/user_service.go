package services

func FindOnlineUsers() ([]string, error) {
	users := make([]string, 0)
	for username, cs := range wsClients {
		if len(cs) > 0 {
			users = append(users, username)
		}
	}
	return users, nil
}
