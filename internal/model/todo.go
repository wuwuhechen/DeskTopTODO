package model

type Todo struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	Priority  string `json:"priority"`
	SortOrder int    `json:"sortOrder"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
