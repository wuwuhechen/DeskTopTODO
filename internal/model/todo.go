package model

type Todo struct {
	ID        string `gorm:"primaryKey" json:"id"`
	NoteID    string `gorm:"not null;default:1" json:"noteId"`
	Content   string `gorm:"not null" json:"content"`
	Completed bool   `gorm:"not null;default:false" json:"completed"`
	Priority  string `gorm:"not null;default:normal" json:"priority"`
	SortOrder int    `gorm:"not null;index" json:"sortOrder"`
	CreatedAt string `gorm:"not null" json:"createdAt"`
	UpdatedAt string `gorm:"not null;default:''" json:"updatedAt"`
}
