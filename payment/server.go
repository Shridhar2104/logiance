package payment

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	service Service
}

func NewServer(service Service) *Server {
	return &Server{service: service}
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/wallet/recharge", s.RechargeWalletHandler)
	mux.HandleFunc("/wallet/deduct", s.DeductBalanceHandler)
	mux.HandleFunc("/wallet/remittance", s.ProcessRemittanceHandler)
	mux.HandleFunc("/wallet/details", s.GetWalletDetailsHandler)
}

func (s *Server) RechargeWalletHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balance, err := s.service.RechargeWallet(r.Context(), req.AccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"account_id": req.AccountID,
		"balance":    balance,
	})
}

func (s *Server) DeductBalanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountID string  `json:"account_id"`
		Amount    float64 `json:"amount"`
		OrderID   string  `json:"order_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balance, err := s.service.DeductBalance(r.Context(), req.AccountID, req.Amount, req.OrderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"account_id": req.AccountID,
		"balance":    balance,
	})
}

func (s *Server) ProcessRemittanceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountID string   `json:"account_id"`
		OrderIDs  []string `json:"order_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	details, err := s.service.ProcessRemittance(r.Context(), req.AccountID, req.OrderIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(details)
}

func (s *Server) GetWalletDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	accountID := r.URL.Query().Get("account_id")
	if accountID == "" {
		http.Error(w, "Missing account_id query parameter", http.StatusBadRequest)
		return
	}

	walletDetails, err := s.service.GetWalletDetails(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(walletDetails)
}
