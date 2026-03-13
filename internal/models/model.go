package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Student struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Course struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Fee       float64   `gorm:"type:decimal(10,2);not null" json:"fee"`
	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Enrollment struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StudentID      uuid.UUID `gorm:"type:uuid;not null" json:"student_id"`
	CourseID       uuid.UUID `gorm:"type:uuid;not null" json:"course_id"`
	EnrollmentDate time.Time `gorm:"type:date;not null" json:"enrollment_date"`
	FeePaid        float64   `gorm:"type:decimal(10,2);default:0" json:"fee_paid"`
	PaymentStatus  string    `gorm:"type:varchar(10);not null;default:'unpaid'" json:"payment_status"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Course  Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

type Payment struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EnrollmentID uuid.UUID `gorm:"type:uuid;not null" json:"enrollment_id"`
	Amount       float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaidAt       time.Time `gorm:"autoCreateTime" json:"paid_at"`

	Enrollment Enrollment `gorm:"foreignKey:EnrollmentID" json:"enrollment,omitempty"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Student{}, &Course{}, &Enrollment{}, &Payment{})
}
