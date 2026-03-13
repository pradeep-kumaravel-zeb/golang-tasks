package repository

import (
	"errors"
	"student-enrollment/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FetchAllEnrollments(courseName, paymentStatus, studentName string) ([]models.Enrollment, error)
	FetchEnrollmentByID(id uuid.UUID) (*models.Enrollment, error)
	CreateEnrollment(enrollment *models.Enrollment) error
	UpdateEnrollment(enrollment *models.Enrollment) error
	DeleteEnrollment(id uuid.UUID) error
	FindStudentByEmail(email string) (*models.Student, error)
	CreateStudent(student *models.Student) error
	UpdateStudent(student *models.Student) error
	FindCourseByName(name string) (*models.Course, error)
	CreateCourse(course *models.Course) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FetchAllEnrollments(courseName, paymentStatus, studentName string) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	query := r.db.Preload("Student").Preload("Course")

	if courseName != "" {
		query = query.Joins("JOIN courses ON courses.id = enrollments.course_id").
			Where("courses.name ILIKE ?", "%"+courseName+"%")
	}

	if paymentStatus != "" {
		query = query.Where("enrollments.payment_status = ?", paymentStatus)
	}

	if studentName != "" {
		query = query.Joins("JOIN students ON students.id = enrollments.student_id").
			Where("students.name ILIKE ?", "%"+studentName+"%")
	}

	if err := query.Find(&enrollments).Error; err != nil {
		return nil, err
	}

	return enrollments, nil
}

func (r *repository) FetchEnrollmentByID(id uuid.UUID) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.Preload("Student").Preload("Course").First(&enrollment, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}
	return &enrollment, nil
}

func (r *repository) CreateEnrollment(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *repository) UpdateEnrollment(enrollment *models.Enrollment) error {
	return r.db.Save(enrollment).Error
}

func (r *repository) DeleteEnrollment(id uuid.UUID) error {
	return r.db.Delete(&models.Enrollment{}, "id = ?", id).Error
}

func (r *repository) FindStudentByEmail(email string) (*models.Student, error) {
	var student models.Student
	if err := r.db.Where("email = ?", email).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

func (r *repository) CreateStudent(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *repository) UpdateStudent(student *models.Student) error {
	return r.db.Save(student).Error
}

func (r *repository) FindCourseByName(name string) (*models.Course, error) {
	var course models.Course
	if err := r.db.Where("name = ?", name).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *repository) CreateCourse(course *models.Course) error {
	return r.db.Create(course).Error
}
