package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
	mux := http.NewServeMux()

	mux.HandleFunc("/api/pharmacy/prescriptions/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/pharmacy/prescriptions/")
		resp := pharmacyResponse{ID: id, PatientID: "P001", Drug: "Amoxicillin", Dose: "500mg", Status: "active", PharmacyName: "CVS Pharmacy", PharmacyType: "Retail"}
		_ = json.NewEncoder(w).Encode(resp)
	})

	mux.HandleFunc("/api/billing/invoices/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/billing/invoices/")
		resp := billingResponse{ID: "INV-" + id, PrescriptionID: id, Amount: 29.99, Status: "paid"}
		_ = json.NewEncoder(w).Encode(resp)
	})

	addr := ":9090"
	log.Printf("mock iris service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
