package main

func GetUserGroup(session string) (int, int, error) {
	var user_id, group_id int
	row := SQL.QueryRow("SELECT `user_id`, `group_id` FROM `websockets` WHERE `client_uuid`=?", session)
	err := row.Scan(&user_id, &group_id)
	if err != nil {
		return 0, 0, err
	}
	return user_id, group_id, err
}

func GetUserName(user_id int) (string, error) {
	var name string
	row := SQL.QueryRow("SELECT `name` FROM `users` WHERE `id`=?", user_id)
	err := row.Scan(&name)
	return name, err
}
