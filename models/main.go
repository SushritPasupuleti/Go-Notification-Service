package models

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Notification struct {
	FromID  User   `json:"from_id"`
	ToID    User   `json:"to_id"`
	Message string `json:"message"`
}
