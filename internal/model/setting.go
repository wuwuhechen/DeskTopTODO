package model

type Setting struct {
	Key   string `gorm:"primaryKey"`
	Value string `gorm:"not null"`
}
