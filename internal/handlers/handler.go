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
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, enrollments)
}

func (h *handler) GetEnrollmentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	enrollment, err := h.service.GetEnrollmentByID(id)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, enrollment)
}

func (h *handler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	enrollment, err := h.service.CreateEnrollment(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, enrollment)
}

func (h *handler) UpdateEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.UpdateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	enrollment, err := h.service.UpdateEnrollment(id, req)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, enrollment)
}

func (h *handler) PatchEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req dto.PatchEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	enrollment, err := h.service.PatchEnrollment(id, req)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, enrollment)
}

func (h *handler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteEnrollment(id)
	if err != nil {
		if err.Error() == "enrollment not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "enrollment deleted successfully"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
