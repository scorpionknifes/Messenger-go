package main

//Join
type Join struct {
	Session string `json:"session"`
}

type Message struct {
	Session string `json:"session"`
	Message string `json:"message"`
}

type Push struct {
	User_id  int    `json:"user_id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

//Join Reply
type JoinReply struct {
	Group_id   int        `json:"group_id"`
	Group_name string     `json:"group_name"`
	Creator_id int        `json:"creator_id"`
	Users      UsersReply `json:"users"`
}

type UsersReply struct {
	User_id  int    `json:"user_id"`
	Username string `json:"username"`
}
