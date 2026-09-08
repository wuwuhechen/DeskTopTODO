package model

// Setting stores application preferences as key/value pairs.
type Setting struct {
	Key   string `gorm:"primaryKey"`
	Value string `gorm:"not null"`
}
