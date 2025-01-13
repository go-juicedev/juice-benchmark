package benchmark

import "time"

type JuiceUser struct {
	ID        int       `column:"id" autoincr:"true" param:"id"`
	Name      string    `column:"name" param:"name"`
	Age       int       `column:"age" param:"age"`
	Email     string    `column:"email" param:"email"`
	CreatedAt time.Time `column:"created_at" param:"created_at"`
	UpdatedAt time.Time `column:"updated_at" param:"updated_at"`
}

type GormUser struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Name      string
	Age       int
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (GormUser) TableName() string {
	return "tbl_user"
}
