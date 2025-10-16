package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Mock data structures matching the new Request/Response naming

// Pharmacy models
type PrescriptionResponse struct {
	ID           string `json:"id"`
	PatientID    string `json:"patient_id"`
	Drug         string `json:"drug"`
	Dose         string `json:"dose"`
	Status       string `json:"status"`
	PharmacyName string `json:"pharmacy_name"`
	PharmacyType string `json:"pharmacy_type"`
}

// Billing models
type InvoiceResponse struct {
	ID             string  `json:"id"`
	PrescriptionID string  `json:"prescription_id"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at,omitempty"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
}

type CreateInvoiceRequest struct {
	PrescriptionID string  `json:"prescription_id"`
	Amount         float64 `json:"amount"`
	Description    string  `json:"description,omitempty"`
}

type AcknowledgeInvoiceRequest struct {
	AcknowledgedBy string `json:"acknowledged_by"`
	Notes          string `json:"notes,omitempty"`
}

type InvoicePaymentResponse struct {
	InvoiceID     string  `json:"invoice_id"`
	PaymentID     string  `json:"payment_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	Status        string  `json:"status"`
	PaidAt        string  `json:"paid_at,omitempty"`
}

type InvoiceListResponse struct {
	PatientID string            `json:"patient_id"`
	Invoices  []InvoiceResponse `json:"invoices"`
	Total     int               `json:"total"`
}

// Stargate auth models
type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(logHeaders) // Custom middleware to log headers

	// Pharmacy API routes
	r.Route("/pharmacy/v1", func(r chi.Router) {
		r.Get("/prescriptions/{prescriptionID}", handleGetPrescription)
	})

	// Billing API routes
	r.Route("/billing/v1", func(r chi.Router) {
		r.Get("/invoices/{prescriptionID}", handleGetInvoice)
		r.Get("/patients/{patientID}/invoices", handleGetInvoicesByPatient)
		r.Post("/invoices", handleCreateInvoice)
		r.Post("/invoices/{invoiceID}/acknowledge", handleAcknowledgeInvoice)
		r.Get("/invoices/{invoiceID}/payment", handleGetInvoicePayment)
	})

	// Stargate OAuth routes
	r.Route("/oauth", func(r chi.Router) {
		r.Post("/token", handleGetToken)
		r.Post("/refresh", handleRefreshToken)
	})

	log.Println("🚀 IRIS Mock Server starting on :8881")
	log.Println("📍 Pharmacy API: http://localhost:8881/pharmacy/v1")
	log.Println("📍 Billing API:  http://localhost:8881/billing/v1")
	log.Println("📍 Stargate Auth: http://localhost:8881/oauth")
	log.Fatal(http.ListenAndServe(":8881", r))
}

// Custom middleware to log headers
func logHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("📥 %s %s", r.Method, r.URL.Path)

		// Log important headers
		if userID := r.Header.Get("X-IRIS-User-ID"); userID != "" {
			log.Printf("   └─ X-IRIS-User-ID: %s", userID)
		}
		if envName := r.Header.Get("X-IRIS-Env-Name"); envName != "" {
			log.Printf("   └─ X-IRIS-Env-Name: %s", envName)
		}
		if idempotency := r.Header.Get("X-Idempotency-Key"); idempotency != "" {
			log.Printf("   └─ X-Idempotency-Key: %s", idempotency)
		}
		if auth := r.Header.Get("Authorization"); auth != "" {
			log.Printf("   └─ Authorization: %s", maskToken(auth))
		}

		next.ServeHTTP(w, r)
	})
}

// maskToken masks the token for logging
func maskToken(auth string) string {
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) == 2 && len(parts[1]) > 10 {
		return parts[0] + " " + parts[1][:10] + "..."
	}
	return auth
}

// Pharmacy handlers
func handleGetPrescription(w http.ResponseWriter, r *http.Request) {
	prescriptionID := chi.URLParam(r, "prescriptionID")

	response := PrescriptionResponse{
		ID:           prescriptionID,
		PatientID:    "PAT-001",
		Drug:         "Lisinopril",
		Dose:         "10mg",
		Status:       "active",
		PharmacyName: "CVS Pharmacy",
		PharmacyType: "Retail",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("✅ Returned prescription: %s", prescriptionID)
}

// Billing handlers
func handleGetInvoice(w http.ResponseWriter, r *http.Request) {
	prescriptionID := chi.URLParam(r, "prescriptionID")

	response := InvoiceResponse{
		ID:             "INV-" + prescriptionID,
		PrescriptionID: prescriptionID,
		Amount:         125.50,
		Status:         "pending",
		CreatedAt:      "2025-10-14T10:00:00Z",
		UpdatedAt:      "2025-10-14T10:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("✅ Returned invoice for prescription: %s", prescriptionID)
}

func handleGetInvoicesByPatient(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")

	// Mock data - return a list of invoices for the patient
	invoices := []InvoiceResponse{
		{
			ID:             "INV-001",
			PrescriptionID: "RX-" + patientID + "-001",
			Amount:         125.50,
			Status:         "paid",
			CreatedAt:      "2025-10-01T10:00:00Z",
			UpdatedAt:      "2025-10-02T14:30:00Z",
		},
		{
			ID:             "INV-002",
			PrescriptionID: "RX-" + patientID + "-002",
			Amount:         89.99,
			Status:         "pending",
			CreatedAt:      "2025-10-10T09:15:00Z",
			UpdatedAt:      "2025-10-10T09:15:00Z",
		},
		{
			ID:             "INV-003",
			PrescriptionID: "RX-" + patientID + "-003",
			Amount:         250.00,
			Status:         "overdue",
			CreatedAt:      "2025-09-15T11:20:00Z",
			UpdatedAt:      "2025-09-15T11:20:00Z",
		},
	}

	response := InvoiceListResponse{
		PatientID: patientID,
		Invoices:  invoices,
		Total:     len(invoices),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("✅ Returned %d invoices for patient: %s", len(invoices), patientID)
}

func handleCreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for idempotency key
	idempotencyKey := r.Header.Get("X-Idempotency-Key")
	log.Printf("💡 Idempotency key: %s", idempotencyKey)

	response := InvoiceResponse{
		ID:             "INV-NEW-" + req.PrescriptionID,
		PrescriptionID: req.PrescriptionID,
		Amount:         req.Amount,
		Status:         "pending",
		CreatedAt:      "2025-10-14T10:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
	log.Printf("✅ Created invoice: %s (Amount: %.2f)", response.ID, response.Amount)
}

func handleAcknowledgeInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "invoiceID")

	var req AcknowledgeInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := InvoiceResponse{
		ID:             invoiceID,
		PrescriptionID: "RX-123",
		Amount:         125.50,
		Status:         "acknowledged",
		UpdatedAt:      "2025-10-14T10:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("✅ Acknowledged invoice: %s by %s", invoiceID, req.AcknowledgedBy)
}

func handleGetInvoicePayment(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "invoiceID")

	response := InvoicePaymentResponse{
		InvoiceID:     invoiceID,
		PaymentID:     "PAY-" + invoiceID,
		Amount:        125.50,
		PaymentMethod: "credit_card",
		Status:        "completed",
		PaidAt:        "2025-10-14T10:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("✅ Returned payment for invoice: %s", invoiceID)
}

// Stargate OAuth handlers
func handleGetToken(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("🔐 Token request from client: %s", req.ClientID)

	// Mock token response
	response := TokenResponse{
		AccessToken:  "mock-access-token-" + req.ClientID,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // 1 hour
		RefreshToken: "mock-refresh-token-" + req.ClientID,
		Scope:        req.Scope,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("✅ Issued token for: %s (expires in %d seconds)", req.ClientID, response.ExpiresIn)
}

func handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("🔄 Token refresh request")

	response := TokenResponse{
		AccessToken:  "mock-refreshed-token-" + req["client_id"],
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: req["refresh_token"],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("✅ Token refreshed")
}
