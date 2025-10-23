package handlers

import (
"encoding/json"
"net/http"
"sync"

"goRoute/middleware"
"goRoute/models"
"goRoute/utils"

"golang.org/x/crypto/bcrypt"
)

var (
users   []models.User
usersMu sync.RWMutex
nextID  = 1
)

func Signup(w http.ResponseWriter, r *http.Request) {
defer r.Body.Close()
w.Header().Set("Content-Type", "application/json")

var req models.SignupRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
return
}

if req.Username == "" || req.Email == "" || req.Password == "" {
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{"error": "All fields required"})
return
}

usersMu.Lock()
for _, user := range users {
if user.Email == req.Email {
usersMu.Unlock()
w.WriteHeader(http.StatusConflict)
json.NewEncoder(w).Encode(map[string]string{"error": "Email exists"})
return
}
}
usersMu.Unlock()

hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
if err != nil {
w.WriteHeader(http.StatusInternalServerError)
json.NewEncoder(w).Encode(map[string]string{"error": "Password hash failed"})
return
}

user := models.User{
ID:       nextID,
Username: req.Username,
Email:    req.Email,
Password: string(hashedPassword),
}
nextID++

usersMu.Lock()
users = append(users, user)
usersMu.Unlock()

token, err := utils.GenerateToken(user.ID, user.Email)
if err != nil {
w.WriteHeader(http.StatusInternalServerError)
json.NewEncoder(w).Encode(map[string]string{"error": "Token generation failed"})
return
}

w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(models.AuthResponse{Token: token, User: user})
}

func Login(w http.ResponseWriter, r *http.Request) {
defer r.Body.Close()
w.Header().Set("Content-Type", "application/json")

var req models.LoginRequest
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
w.WriteHeader(http.StatusBadRequest)
json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
return
}

usersMu.RLock()
var foundUser *models.User
for i := range users {
if users[i].Email == req.Email {
foundUser = &users[i]
break
}
}
usersMu.RUnlock()

if foundUser == nil {
w.WriteHeader(http.StatusUnauthorized)
json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
return
}

if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password)); err != nil {
w.WriteHeader(http.StatusUnauthorized)
json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
return
}

token, err := utils.GenerateToken(foundUser.ID, foundUser.Email)
if err != nil {
w.WriteHeader(http.StatusInternalServerError)
json.NewEncoder(w).Encode(map[string]string{"error": "Token generation failed"})
return
}

w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(models.AuthResponse{Token: token, User: *foundUser})
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Content-Type", "application/json")

claims := r.Context().Value(middleware.UserContextKey).(*utils.Claims)

usersMu.RLock()
defer usersMu.RUnlock()

for _, user := range users {
if user.ID == claims.UserID {
json.NewEncoder(w).Encode(user)
return
}
}

w.WriteHeader(http.StatusNotFound)
json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
}
