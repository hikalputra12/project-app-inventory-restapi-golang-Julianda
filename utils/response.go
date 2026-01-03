package utils

import (
	"app-inventory/dto"
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func ResponseError(w http.ResponseWriter, code int, message string, errs interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := Response{
		Status:  false,
		Message: message,
		Errors:  errs,
	}

	json.NewEncoder(w).Encode(response)
}

// ResponseJSON digunakan untuk mengirim respon sukses (200, 201)
func ResponseJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := Response{
		Status:  true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}
func ResponseSuccess(w http.ResponseWriter, code int, message string, data any) {
	response := Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func ResponseBadRequest(w http.ResponseWriter, code int, message string, errors any) {
	response := Response{
		Status:  false,
		Message: message,
		Errors:  errors,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func ResponsePagination(w http.ResponseWriter, code int, message string, data any, pagination dto.Pagination) {
	response := map[string]interface{}{
		"status":     true,
		"message":    message,
		"data":       data,
		"pagination": pagination,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
