package service

import (
	"errors"
	"student-enrollment/internal/dto"
	"student-enrollment/internal/models"
	"student-enrollment/internal/repository"
	"student-enrollment/internal/utils"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetEnrollments(courseName, paymentStatus, studentName string) ([]dto.EnrollmentResponse, error)
	GetEnrollmentByID(id string) (*dto.EnrollmentResponse, error)
	CreateEnrollment(req dto.CreateEnrollmentRequest) (*dto.EnrollmentResponse, error)
	UpdateEnrollment(id string, req dto.UpdateEnrollmentRequest) (*dto.EnrollmentResponse, error)
	PatchEnrollment(id string, req dto.PatchEnrollmentRequest) (*dto.EnrollmentResponse, error)
	DeleteEnrollment(id string) error
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetEnrollments(courseName, paymentStatus, studentName string) ([]dto.EnrollmentResponse, error) {
	enrollments, err := s.repo.FetchAllEnrollments(courseName, paymentStatus, studentName)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.EnrollmentResponse, len(enrollments))
	for i, e := range enrollments {
		responses[i] = dto.EnrollmentResponse{
			ID:             e.ID.String(),
			StudentName:    e.Student.Name,
			StudentEmail:   e.Student.Email,
			CourseName:     e.Course.Name,
			CourseFee:      e.Course.Fee,
			EnrollmentDate: e.EnrollmentDate,
			FeePaid:        e.FeePaid,
			PaymentStatus:  e.PaymentStatus,
		}
	}

	return responses, nil
}

func (s *service) GetEnrollmentByID(id string) (*dto.EnrollmentResponse, error) {
	enrollmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid enrollment ID")
	}

	enrollment, err := s.repo.FetchEnrollmentByID(enrollmentID)
	if err != nil {
		return nil, err
	}

	return &dto.EnrollmentResponse{
		ID:             enrollment.ID.String(),
		StudentName:    enrollment.Student.Name,
		StudentEmail:   enrollment.Student.Email,
		CourseName:     enrollment.Course.Name,
		CourseFee:      enrollment.Course.Fee,
		EnrollmentDate: enrollment.EnrollmentDate,
		FeePaid:        enrollment.FeePaid,
		PaymentStatus:  enrollment.PaymentStatus,
	}, nil
}

func (s *service) CreateEnrollment(req dto.CreateEnrollmentRequest) (*dto.EnrollmentResponse, error) {
	if !utils.ValidateEmail(req.StudentEmail) {
		return nil, errors.New("invalid email format")
	}

	if !utils.ValidateName(req.StudentName) {
		return nil, errors.New("student name is required")
	}

	if !utils.ValidateText(req.CourseName) {
		return nil, errors.New("course name is required")
	}

	// Find or create student
	student, err := s.repo.FindStudentByEmail(req.StudentEmail)
	if err != nil {
		return nil, err
	}

	if student == nil {
		student = &models.Student{
			Name:  req.StudentName,
			Email: req.StudentEmail,
		}
		if err := s.repo.CreateStudent(student); err != nil {
			return nil, err
		}
	} else {
		// Update student name if different
		if student.Name != req.StudentName {
			student.Name = req.StudentName
			if err := s.repo.UpdateStudent(student); err != nil {
				return nil, err
			}
		}
	}

	// Find or create course
	course, err := s.repo.FindCourseByName(req.CourseName)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found - please create the course first")
	}

	// Calculate payment status
	paymentStatus := utils.CalculatePaymentStatus(req.FeePaid, course.Fee)

	// Create enrollment
	enrollment := &models.Enrollment{
		StudentID:      student.ID,
		CourseID:       course.ID,
		EnrollmentDate: time.Now(),
		FeePaid:        req.FeePaid,
		PaymentStatus:  paymentStatus,
	}

	if err := s.repo.CreateEnrollment(enrollment); err != nil {
		return nil, err
	}

	return &dto.EnrollmentResponse{
		ID:             enrollment.ID.String(),
		StudentName:    student.Name,
		StudentEmail:   student.Email,
		CourseName:     course.Name,
		CourseFee:      course.Fee,
		EnrollmentDate: enrollment.EnrollmentDate,
		FeePaid:        enrollment.FeePaid,
		PaymentStatus:  enrollment.PaymentStatus,
	}, nil
}

func (s *service) UpdateEnrollment(id string, req dto.UpdateEnrollmentRequest) (*dto.EnrollmentResponse, error) {
	enrollmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid enrollment ID")
	}

	enrollment, err := s.repo.FetchEnrollmentByID(enrollmentID)
	if err != nil {
		return nil, err
	}

	if !utils.ValidateEmail(req.StudentEmail) {
		return nil, errors.New("invalid email format")
	}

	if !utils.ValidateName(req.StudentName) {
		return nil, errors.New("student name is required")
	}

	if !utils.ValidateText(req.CourseName) {
		return nil, errors.New("course name is required")
	}

	// Find or create student
	student, err := s.repo.FindStudentByEmail(req.StudentEmail)
	if err != nil {
		return nil, err
	}

	if student == nil {
		student = &models.Student{
			Name:  req.StudentName,
			Email: req.StudentEmail,
		}
		if err := s.repo.CreateStudent(student); err != nil {
			return nil, err
		}
	} else if student.Name != req.StudentName {
		student.Name = req.StudentName
		if err := s.repo.UpdateStudent(student); err != nil {
			return nil, err
		}
	}

	// Find course
	course, err := s.repo.FindCourseByName(req.CourseName)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, errors.New("course not found")
	}

	// Update enrollment
	enrollment.StudentID = student.ID
	enrollment.CourseID = course.ID
	enrollment.FeePaid = req.FeePaid
	enrollment.PaymentStatus = utils.CalculatePaymentStatus(req.FeePaid, course.Fee)

	if err := s.repo.UpdateEnrollment(enrollment); err != nil {
		return nil, err
	}

	return &dto.EnrollmentResponse{
		ID:             enrollment.ID.String(),
		StudentName:    student.Name,
		StudentEmail:   student.Email,
		CourseName:     course.Name,
		CourseFee:      course.Fee,
		EnrollmentDate: enrollment.EnrollmentDate,
		FeePaid:        enrollment.FeePaid,
		PaymentStatus:  enrollment.PaymentStatus,
	}, nil
}

func (s *service) PatchEnrollment(id string, req dto.PatchEnrollmentRequest) (*dto.EnrollmentResponse, error) {
	enrollmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid enrollment ID")
	}

	enrollment, err := s.repo.FetchEnrollmentByID(enrollmentID)
	if err != nil {
		return nil, err
	}

	student := &enrollment.Student
	course := &enrollment.Course

	// Update student email if provided
	if req.StudentEmail != nil {
		if !utils.ValidateEmail(*req.StudentEmail) {
			return nil, errors.New("invalid email format")
		}

		if student.Email != *req.StudentEmail {
			existingStudent, err := s.repo.FindStudentByEmail(*req.StudentEmail)
			if err != nil {
				return nil, err
			}

			if existingStudent != nil && existingStudent.ID != student.ID {
				return nil, errors.New("email already exists for another student")
			}

			student.Email = *req.StudentEmail
			if err := s.repo.UpdateStudent(student); err != nil {
				return nil, err
			}
		}
	}

	// Update student name if provided
	if req.StudentName != nil {
		if !utils.ValidateName(*req.StudentName) {
			return nil, errors.New("student name cannot be empty")
		}
		student.Name = *req.StudentName
		if err := s.repo.UpdateStudent(student); err != nil {
			return nil, err
		}
	}

	// Update course if provided
	if req.CourseName != nil {
		if !utils.ValidateText(*req.CourseName) {
			return nil, errors.New("course name cannot be empty")
		}

		newCourse, err := s.repo.FindCourseByName(*req.CourseName)
		if err != nil {
			return nil, err
		}

		if newCourse == nil {
			return nil, errors.New("course not found")
		}

		enrollment.CourseID = newCourse.ID
		course = newCourse
	}

	// Update fee paid if provided
	if req.FeePaid != nil {
		enrollment.FeePaid = *req.FeePaid
	}

	// Recalculate payment status
	enrollment.PaymentStatus = utils.CalculatePaymentStatus(enrollment.FeePaid, course.Fee)

	if err := s.repo.UpdateEnrollment(enrollment); err != nil {
		return nil, err
	}

	return &dto.EnrollmentResponse{
		ID:             enrollment.ID.String(),
		StudentName:    student.Name,
		StudentEmail:   student.Email,
		CourseName:     course.Name,
		CourseFee:      course.Fee,
		EnrollmentDate: enrollment.EnrollmentDate,
		FeePaid:        enrollment.FeePaid,
		PaymentStatus:  enrollment.PaymentStatus,
	}, nil
}

func (s *service) DeleteEnrollment(id string) error {
	enrollmentID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid enrollment ID")
	}

	enrollment, err := s.repo.FetchEnrollmentByID(enrollmentID)
	if err != nil {
		return err
	}

	// Check if payment status is paid
	if enrollment.PaymentStatus == "paid" {
		return errors.New("cannot delete enrollment with paid status - refund required first")
	}

	// Check if course has started
	if time.Now().After(enrollment.Course.StartDate) {
		return errors.New("cannot delete enrollment - course has already started")
	}

	return s.repo.DeleteEnrollment(enrollmentID)
}
