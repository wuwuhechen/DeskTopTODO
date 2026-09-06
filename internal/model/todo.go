package model

type Todo struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	SortOrder int    `json:"sortOrder"`
	CreatedAt string `json:"createdAt"`
}
