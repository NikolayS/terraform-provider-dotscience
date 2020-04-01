package types

type Runner struct {
	ID            string  `json:"id"`
	AccountID     string  `json:"account_id"`
	Name          string  `json:"name"`
	Status        string  `json:"status"`
	Managed       bool    `json:"managed"`
	ServerState   string  `json:"server_state"`
	StatusMessage string  `json:"status_message"`
	RunnerProfile string  `json:"runner_profile"`
	Tasks         []*Task `json:"tasks"`
}
