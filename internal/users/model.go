package users

import "time"

type User struct {
	ID           int64     `gorm:"primaryKey"`
	Username     string    `gorm:"type:varchar(50);not null;unique"`
	PasswordHash string    `gorm:"type:varchar(100);not null"`
	Email        *string   `gorm:"type:varchar(255)"`
	IsActive     bool      `gorm:"type:boolean;default:true"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`
}
