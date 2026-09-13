package dtos

type AgentDto struct {
	Id      string `db:"id" json:"id"`
	Person  string `db:"person" json:"person"`
	Cobrand string `db:"cobrand" json:"cobrand"`
}

type PropertyAgentDto struct {
	Id       string `db:"id" json:"id"`
	Agent    string `db:"agent" json:"agent"`
	Property string `db:"property" json:"property"`
}
