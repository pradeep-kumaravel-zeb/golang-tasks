package handlers

import (
	"encoding/json"
	"net/http"
	"student-enrollment/internal/dto"
	"student-enrollment/internal/service"

	"github.com/gorilla/mux"
)

type HandlerInterface interface {
	GetAllEnrollments(w http.ResponseWriter, r *http.Request)
	GetEnrollmentByID(w http.ResponseWriter, r *http.Request)
	CreateEnrollment(w http.ResponseWriter, r *http.Request)
	UpdateEnrollment(w http.ResponseWriter, r *http.Request)
	PatchEnrollment(w http.ResponseWriter, r *http.Request)
	DeleteEnrollment(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service service.Service
}

func NewHandler(service service.Service) HandlerInterface {
	return &handler{service: service}
}

func (h *handler) GetAllEnrollments(w http.ResponseWriter, r *http.Request) {
	courseName := r.URL.Query().Get("course_name")
	paymentStatus := r.URL.Query().Get("payment_status")
	studentName := r.URL.Query().Get("student_name")

	enrollments, err := h.service.GetEnrollments(courseName, paymentStatus, studentName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	respondWithSuccess(w, http.StatusOK, "Enrollments retrieved successfully", enrollments)
}

func (h *handler) GetEnrollmentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	enrollment, err := h.service.GetEnrollmentByID(id)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondWithError(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respondWithSuccess(w, http.StatusOK, "Enrollment retrieved successfully", enrollment)
}

func (h *handler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	enrollment, err := h.service.CreateEnrollment(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respondWithSuccess(w, http.StatusCreated, "Enrollment created successfully", enrollment)
}

func (h *handler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.UpdateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	enrollment, err := h.service.UpdateEnrollment(id, req)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondWithError(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respondWithSuccess(w, http.StatusOK, "Enrollment updated successfully", enrollment)
}

func (h *handler) PatchEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.PatchEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	enrollment, err := h.service.PatchEnrollment(id, req)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondWithError(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respondWithSuccess(w, http.StatusOK, "Enrollment patched successfully", enrollment)
}

func (h *handler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteEnrollment(id)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondWithError(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respondWithSuccess(w, http.StatusOK, "Enrollment deleted successfully", nil)
}

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Status        int                 `json:"status"`
	Message       string              `json:"message"`
	Data          []interface{}       `json:"data"`
	MissingFields []map[string]string `json:"missingFields,omitempty"`
}

func respondWithSuccess(w http.ResponseWriter, code int, message string, data interface{}) {
	if data == nil {
		data = []interface{}{}
	}
	response := SuccessResponse{
		Status:  code,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, code int, message string, missingFields []map[string]string) {
	response := ErrorResponse{
		Status:        code,
		Message:       message,
		Data:          []interface{}{},
		MissingFields: missingFields,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
