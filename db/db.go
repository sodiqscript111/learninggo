package db

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string
	Email     string `gorm:"uniqueIndex"`
	Password  string
	Role      string `gorm:"type:text;default:'user'"`
	CreatedAt time.Time
	Events    []Event `gorm:"foreignKey:CreatedBy"`
}

type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedAt   time.Time
	Title       string
	Description string
	Location    string
	StartTime   time.Time
	EndTime     time.Time
	MaxCapacity int
	IsPaid      bool
	Price       float64 `gorm:"type:decimal(10,2)"`

	CreatedBy uuid.UUID `gorm:"type:uuid"`
	User      User      `gorm:"foreignKey:CreatedBy;references:ID"`
}

type Registration struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"` // foreign key field
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	EventID uuid.UUID `gorm:"type:uuid;not null"` // foreign key field
	Event   Event     `gorm:"foreignKey:EventID;references:ID"`

	RegisteredAt  time.Time `gorm:"autoCreateTime"`
	Status        string    `gorm:"type:text"`
	TicketToken   string    `gorm:"type:text;uniqueIndex"`
	PaymentStatus string    `gorm:"type:text"`
}

type EventLink struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	EventID uuid.UUID `gorm:"type:uuid;not null"`
	Event   Event     `gorm:"foreignKey:EventID;references:ID"`

	Code      string `gorm:"type:text;uniqueIndex"`
	ExpiresAt *time.Time
	MaxUses   *int
	UsesCount int

	CreatedBy uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:CreatedBy;references:ID"`

	CreatedAt time.Time
	IsActive  bool `gorm:"default:true"`
}

type Waitlist struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	EventID   uuid.UUID `gorm:"type:uuid;not null"`
	Event     Event     `gorm:"foreignKey:EventID;references:ID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Attendances struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	User        User      `gorm:"foreignKey:UserID;references:ID"`
	EventID     uuid.UUID `gorm:"type:uuid;not null"`
	Event       Event     `gorm:"foreignKey:EventID;references:ID"`
	CreatedInAt time.Time `gorm:"autoCreateTime"`
}

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=password dbname=Event port=5432 sslmode=disable"
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatal("Failed to enable uuid-ossp extension: ", err)
	}

	err = CreateTable()
	if err != nil {
		log.Fatal("Failed to create users table: ", err)
	}

}

func CreateTable() error {
	return DB.AutoMigrate(
		&User{}, &Event{}, &Registration{},
		&EventLink{}, &Waitlist{}, &Attendances{},
	)

}
