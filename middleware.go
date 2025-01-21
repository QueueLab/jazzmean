package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type Middleware struct {
	dbPool *pgxpool.Pool
	mu     sync.Mutex
}

func NewMiddleware(dbURL string) (*Middleware, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Middleware{dbPool: dbPool}, nil
}

func (m *Middleware) Close() {
	m.dbPool.Close()
}

func (m *Middleware) HandleQuery(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var query struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Process the query with AI SDK and PostgreSQL
	response, err := m.processQuery(query.Query)
	if err != nil {
		http.Error(w, "Failed to process query", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (m *Middleware) HandleSubtract(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var params struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result := params.A - params.B

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"result": result})
}

func (m *Middleware) processQuery(query string) (interface{}, error) {
	if query == "subtract" {
		return map[string]string{"response": "Subtraction query received"}, nil
	}
	// Placeholder for AI SDK and PostgreSQL processing logic
	return map[string]string{"response": "Processed query: " + query}, nil
}
