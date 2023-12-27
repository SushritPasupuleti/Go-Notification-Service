package models

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Notification struct {
	FromID  int   `json:"from_id"`
	ToID    int   `json:"to_id"`
	Message string `json:"message"`
}
