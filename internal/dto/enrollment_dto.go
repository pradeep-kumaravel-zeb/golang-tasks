package dto

import "time"

type CreateEnrollmentRequest struct {
	StudentName  string  `json:"student_name"`
	StudentEmail string  `json:"student_email"`
	CourseName   string  `json:"course_name"`
	FeePaid      float64 `json:"fee_paid"`
}

type UpdateEnrollmentRequest struct {
	StudentName  string  `json:"student_name"`
	StudentEmail string  `json:"student_email"`
	CourseName   string  `json:"course_name"`
	FeePaid      float64 `json:"fee_paid"`
}

type PatchEnrollmentRequest struct {
	StudentName  *string  `json:"student_name,omitempty"`
	StudentEmail *string  `json:"student_email,omitempty"`
	CourseName   *string  `json:"course_name,omitempty"`
	FeePaid      *float64 `json:"fee_paid,omitempty"`
}

type EnrollmentResponse struct {
	ID             string    `json:"id"`
	StudentName    string    `json:"student_name"`
	StudentEmail   string    `json:"student_email"`
	CourseName     string    `json:"course_name"`
	CourseFee      float64   `json:"course_fee"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	FeePaid        float64   `json:"fee_paid"`
	PaymentStatus  string    `json:"payment_status"`
}
