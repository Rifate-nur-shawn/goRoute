package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"crypto/rand"
	"encoding/hex"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Product struct {
    ID          int     `json:"id"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    ImgUrl      string  `json:"imageUrl"`
}

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"-"` // Don't expose password in JSON
}

type SignupRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Token string          `json:"token"`
    User  UserResponse    `json:"user"`
}

type UserResponse struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}


var (
    productList []Product
    productMu sync.RWMutex 
    userList []User
    userMu sync.RWMutex
    userIDCounter = 1
)

func corsMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close()
		w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions{
			w.WriteHeader(http.StatusNoContent)
            return
		}
		next.ServeHTTP(w, r)
	})	
}


func HandleName(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{
        "message": "My name is khan",
    }
    json.NewEncoder(w).Encode(response)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
     productMu.RLock()
     product:=productList
     productMu.RUnlock()

       if err := json.NewEncoder(w).Encode(product); err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
    }
}

func createProduct(w http.ResponseWriter, r *http.Request){
    defer r.Body.Close()
	 var newProduct Product
    
    if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid JSON format",
        })
        return
    }

    if newProduct.Title ==""{
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error":"Title is required",
        })
        return
    }

    if newProduct.Price <= 0 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Price must be greater than 0",
        })
        return
    }
        if len(newProduct.Description) > 500 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Description must be less than 500 characters",
        })
        return
    }

    productMu.Lock()
    newProduct.ID = len(productList) + 1
    productList = append(productList, newProduct)
    productMu.Unlock()

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newProduct)
}

// Generate a random token
func generateToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return hex.EncodeToString(b)
}

// Signup handler
func handleSignup(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    var req SignupRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid JSON format",
        })
        return
    }

    // Validate input
    if req.Username == "" || req.Email == "" || req.Password == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Username, email, and password are required",
        })
        return
    }

    if len(req.Password) < 6 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Password must be at least 6 characters",
        })
        return
    }

    userMu.Lock()
    // Check if user already exists
    for _, u := range userList {
        if u.Email == req.Email {
            userMu.Unlock()
            w.WriteHeader(http.StatusConflict)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "User with this email already exists",
            })
            return
        }
    }

    // Create new user
    newUser := User{
        ID:       userIDCounter,
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password, // In production, hash the password!
    }
    userIDCounter++
    userList = append(userList, newUser)
    userMu.Unlock()

    // Generate token
    token := generateToken()

    // Create response
    response := AuthResponse{
        Token: token,
        User: UserResponse{
            ID:       newUser.ID,
            Username: newUser.Username,
            Email:    newUser.Email,
        },
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}

// Login handler
func handleLogin(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    var req LoginRequest

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid JSON format",
        })
        return
    }

    // Validate input
    if req.Email == "" || req.Password == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Email and password are required",
        })
        return
    }

    userMu.RLock()
    var foundUser *User
    for _, u := range userList {
        if u.Email == req.Email && u.Password == req.Password {
            foundUser = &u
            break
        }
    }
    userMu.RUnlock()

    if foundUser == nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid email or password",
        })
        return
    }

    // Generate token
    token := generateToken()

    // Create response
    response := AuthResponse{
        Token: token,
        User: UserResponse{
            ID:       foundUser.ID,
            Username: foundUser.Username,
            Email:    foundUser.Email,
        },
    }

    json.NewEncoder(w).Encode(response)
}


func main(){
	r:= chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(corsMiddleware)


	r.Get("/name",HandleName)
	r.Get("/products",getProducts)
	r.Post("/products",createProduct)
	
	// Auth routes
	r.Post("/api/auth/signup", handleSignup)
	r.Post("/api/auth/login", handleLogin)

	fmt.Println(" Chi Router server running on :8080")
    fmt.Println(" Available routes:")
    fmt.Println(" GET /name")
    fmt.Println(" GET /products")
    fmt.Println(" POST /products")
    fmt.Println(" POST /api/auth/signup")
    fmt.Println(" POST /api/auth/login")



	err := http.ListenAndServe(":8080", r)
    if err != nil {
        fmt.Println("Server error:", err)
    }
}

func init() {
    productList = []Product{
        {
            ID:          1,
            Title:       "orange",
            Description: "Orange is red",
            Price:       100.33,
            ImgUrl:      "https://media.istockphoto.com/id/1231559990/photo/orange-fruit-with-drop-shadow-on-white-background-commercial-image-of-citrus-fruits-in.jpg?s=612x612&w=0&k=20&c=zEnUx_53uqE-GBCLar_fK4PwJBG3U2pV0vSu0amRFDE=",
        },
        {
            ID:          2,
            Title:       "mango",
            Description: "Mango is Green",
            Price:       150.33,
            ImgUrl:      "https://www.foodcraft.hk/cdn/shop/files/GreenMango-1pc.jpg?v=1694480506",
        },
        {
            ID:          3,
            Title:       "Banana",
            Description: "Banana is Yellow",
            Price:       15.00,
            ImgUrl:      "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQG7ElBNPs-HbYJJOMHRu7lEmphTn8-52FYKw&s",
        },
    }
}