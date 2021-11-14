package entity

type User struct {
	ID       int `gorm:"primaryKey;autoIncrement:false"`
	Username string
	Password string
	Wallets  string
}
