package handlers

import (
"encoding/json"
"net/http"
"strconv"
"sync"

"goRoute/models"

"github.com/go-chi/chi/v5"
)

var (
products   []models.Product
productsMu sync.RWMutex
nextProdID = 1
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
productsMu.RLock()
defer productsMu.RUnlock()
json.NewEncoder(w).Encode(products)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
id, _ := strconv.Atoi(chi.URLParam(r, "id"))

productsMu.RLock()
defer productsMu.RUnlock()

for _, prod := range products {
if prod.ID == id {
json.NewEncoder(w).Encode(prod)
return
}
}

w.WriteHeader(http.StatusNotFound)
json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
defer r.Body.Close()

var req models.Product
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
return
}

productsMu.Lock()
req.ID = nextProdID
nextProdID++
products = append(products, req)
productsMu.Unlock()

w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(req)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
defer r.Body.Close()
id, _ := strconv.Atoi(chi.URLParam(r, "id"))

var req models.Product
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
return
}

productsMu.Lock()
defer productsMu.Unlock()

for i, prod := range products {
if prod.ID == id {
req.ID = id
products[i] = req
json.NewEncoder(w).Encode(req)
return
}
}

w.WriteHeader(http.StatusNotFound)
json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")
id, _ := strconv.Atoi(chi.URLParam(r, "id"))

productsMu.Lock()
defer productsMu.Unlock()

for i, prod := range products {
if prod.ID == id {
products = append(products[:i], products[i+1:]...)
w.WriteHeader(http.StatusNoContent)
return
}
}

w.WriteHeader(http.StatusNotFound)
json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}
