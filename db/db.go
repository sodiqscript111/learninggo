package db

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"password"`
	Role      string    `gorm:"type:text;default:'user'" json:"role"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
	Events    []Event   `gorm:"foreignKey:CreatedBy" json:"events,omitempty"`
}

type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `gorm:"index:idx_events_start_time_desc" json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	MaxCapacity int       `json:"maxCapacity"`
	IsPaid      bool      `json:"isPaid"`
	Price       float64   `gorm:"type:decimal(10,2)" json:"price"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null;index" json:"createdBy"`
	User      User      `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

type Registration struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index:idx_registration_user_event,priority:1" json:"userId"`
	User          User      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	EventID       uuid.UUID `gorm:"type:uuid;not null;index:idx_registration_user_event,priority:2" json:"eventId"`
	Event         Event     `gorm:"foreignKey:EventID;references:ID" json:"-"`
	RegisteredAt  time.Time `gorm:"autoCreateTime;index" json:"registeredAt"`
	Status        string    `gorm:"type:text;index" json:"status"`
	TicketToken   string    `gorm:"type:text;uniqueIndex" json:"ticketToken"`
	PaymentStatus string    `gorm:"type:text;index" json:"paymentStatus"`
}

type EventLink struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"eventId"`
	Event     Event      `gorm:"foreignKey:EventID;references:ID" json:"-"`
	Code      string     `gorm:"type:text;uniqueIndex" json:"code"`
	ExpiresAt *time.Time `gorm:"index" json:"expiresAt,omitempty"`
	MaxUses   *int       `json:"maxUses,omitempty"`
	UsesCount int        `json:"usesCount"`

	CreatedBy uuid.UUID `gorm:"type:uuid;not null;index" json:"createdBy"`
	User      User      `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	IsActive  bool      `gorm:"default:true" json:"isActive"`
}

type Waitlist struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_waitlist_user_event,priority:1" json:"userId"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	EventID   uuid.UUID `gorm:"type:uuid;not null;index:idx_waitlist_user_event,priority:2" json:"eventId"`
	Event     Event     `gorm:"foreignKey:EventID;references:ID" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
}

type Attendances struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_attendances_user_event,priority:1" json:"userId"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	EventID   uuid.UUID `gorm:"type:uuid;not null;index:idx_attendances_user_event,priority:2" json:"eventId"`
	Event     Event     `gorm:"foreignKey:EventID;references:ID" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"createdAt"`
}

func ConnectDatabase() {
	dsn := "host=db user=postgres password=password dbname=Event port=5432 sslmode=disable"
	var err error

	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	DB, err = gorm.Open(postgres.Open(dsn), dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB: ", err)
	}

	sqlDB.SetMaxOpenConns(550) // Adjusted for 3000 req/s with scalability
	sqlDB.SetMaxIdleConns(150) // Balanced for load spikes
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(3 * time.Minute) // Reduced for resource efficiency

	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		log.Fatal("Failed to enable uuid-ossp extension: ", err)
	}

	err = CreateTable()
	if err != nil {
		log.Fatal("Failed to create users table: ", err)
	}

	// Ensure optimized index for GetAllEvents
	err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_events_start_time_desc ON events (start_time DESC)").Error
	if err != nil {
		log.Printf("Failed to create index idx_events_start_time_desc: %v", err)
	}

	// Remove conflicting composite index
	err = DB.Exec("DROP INDEX IF EXISTS idx_event_time").Error
	if err != nil {
		log.Printf("Failed to drop index idx_event_time: %v", err)
	}

	// Optimize table performance
	err = DB.Exec("ANALYZE VERBOSE events").Error
	if err != nil {
		log.Printf("Failed to analyze events table: %v", err)
	}
	err = DB.Exec("VACUUM VERBOSE events").Error
	if err != nil {
		log.Printf("Failed to vacuum events table: %v", err)
	}

	// Enable pg_stat_statements for diagnostics
	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS pg_stat_statements").Error
	if err != nil {
		log.Printf("Failed to enable pg_stat_statements: %v", err)
	}
}

func CreateTable() error {
	return DB.AutoMigrate(
		&User{}, &Event{}, &Registration{},
		&EventLink{}, &Waitlist{}, &Attendances{},
	)
}
