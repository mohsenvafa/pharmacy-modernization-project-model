package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type pharmacyResponse struct {
	ID           string `json:"id"`
	PatientID    string `json:"patient_id"`
	Drug         string `json:"drug"`
	Dose         string `json:"dose"`
	Status       string `json:"status"`
	PharmacyName string `json:"pharmacy_name"`
	PharmacyType string `json:"pharmacy_type"`
}

type billingResponse struct {
	ID             string  `json:"id"`
	PrescriptionID string  `json:"prescription_id"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
}

func main() {
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/api/pharmacy/prescriptions/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		resp := pharmacyResponse{ID: id, PatientID: "P001", Drug: "Amoxicillin", Dose: "500mg", Status: "active", PharmacyName: "CVS Pharmacy", PharmacyType: "Retail"}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	r.Get("/api/billing/invoices/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		resp := billingResponse{ID: "INV-" + id, PrescriptionID: id, Amount: 29.99, Status: "paid"}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	addr := ":9090"
	log.Printf("mock iris service listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
