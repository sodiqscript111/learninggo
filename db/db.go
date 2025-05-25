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
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"password"`
	Role      string    `gorm:"type:text;default:'user'" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	Events    []Event   `gorm:"foreignKey:CreatedBy" json:"events,omitempty"`
}

type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	MaxCapacity int       `json:"maxCapacity"`
	IsPaid      bool      `json:"isPaid"`
	Price       float64   `gorm:"type:decimal(10,2)" json:"price"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"createdBy"`
	User      User      `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

type Registration struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	User          User      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	EventID       uuid.UUID `gorm:"type:uuid;not null" json:"eventId"`
	Event         Event     `gorm:"foreignKey:EventID;references:ID" json:"-"`
	RegisteredAt  time.Time `gorm:"autoCreateTime" json:"registeredAt"`
	Status        string    `gorm:"type:text" json:"status"`
	TicketToken   string    `gorm:"type:text;uniqueIndex" json:"ticketToken"`
	PaymentStatus string    `gorm:"type:text" json:"paymentStatus"`
}

type EventLink struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventID   uuid.UUID  `gorm:"type:uuid;not null" json:"eventId"`
	Event     Event      `gorm:"foreignKey:EventID;references:ID" json:"-"`
	Code      string     `gorm:"type:text;uniqueIndex" json:"code"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	MaxUses   *int       `json:"maxUses,omitempty"`
	UsesCount int        `json:"usesCount"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"createdBy"`
	User      User      `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	IsActive  bool      `gorm:"default:true" json:"isActive"`
}

type Waitlist struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	EventID   uuid.UUID `gorm:"type:uuid;not null" json:"eventId"`
	Event     Event     `gorm:"foreignKey:EventID;references:ID" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type Attendances struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	User        User      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	EventID     uuid.UUID `gorm:"type:uuid;not null" json:"eventId"`
	Event       Event     `gorm:"foreignKey:EventID;references:ID" json:"-"`
	CreatedInAt time.Time `gorm:"autoCreateTime" json:"createdInAt"`
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
