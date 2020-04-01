package types

type Task struct {
	ID            string `json:"id"`
	DesiredStatus string `json:"desired_status"`
	Status        string `json:"status"`
}
