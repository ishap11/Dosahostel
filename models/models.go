package models

type User struct {
	UserID   uint   `gorm:"primaryKey" json:"user_id"`
	FullName string `gorm:"not null" json:"full_name"`
	Email    string `gorm:"not null;unique" json:"email"`
	Phone    string `gorm:"not null" json:"phone"`
	Password string `gorm:"not null" json:"password"`
	Type     string `gorm:"not null" json:"type"`
	BlockID  uint   `gorm:"foreignKey:BlockID;references:BlockID;onDelete:SET NULL" json:"block_id"`
	USN      string `gorm:"unique" json:"usn"`
	Room     string `json:"room"`
}

type Block struct {
	BlockID   uint      `gorm:"primaryKey" json:"block_id"`
	BlockName string    `gorm:"not null" json:"block_name"`
	Users     []User    `gorm:"foreignKey:BlockID" json:"users"`
	Students  []Student `gorm:"foreignKey:BlockID" json:"students"`
	Wardens   []Warden  `gorm:"foreignKey:BlockID" json:"wardens"`
}

type Student struct {
	StudentID uint   `gorm:"primaryKey" json:"student_id"`
	UserID    uint   `gorm:"not null;unique" json:"user_id"` // Foreign key to User
	FullName  string `gorm:"not null" json:"full_name"`
	Email     string `gorm:"not null;unique" json:"email"`
	Phone     string `gorm:"not null" json:"phone"`
	USN       string `gorm:"not null;unique" json:"usn"`
	BlockID   uint   `gorm:"foreignKey:BlockID;references:BlockID;onDelete:CASCADE" json:"block_id"`
	Room      string `json:"room"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
	Block     Block  `gorm:"foreignKey:BlockID" json:"block"`
}

type Warden struct {
	WardenID uint  `gorm:"primaryKey" json:"warden_id"`
	UserID   uint  `gorm:"not null;unique" json:"user_id"` // Foreign key to User
	BlockID  uint  `gorm:"foreignKey:BlockID;references:BlockID;onDelete:CASCADE" json:"block_id"`
	User     User  `gorm:"foreignKey:UserID" json:"user"`
	Block    Block `gorm:"foreignKey:BlockID" json:"block"`
}
