package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"hospital-management/backend/database"
	"hospital-management/backend/utils"
)

// Response represents a generic API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginRequest represents the login request structure
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest represents the registration request structure
type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
}

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"-"` // "-" prevents password from being sent in JSON responses
	Email     string    `bson:"email" json:"email"`
	FirstName string    `bson:"firstName" json:"firstName"`
	LastName  string    `bson:"lastName" json:"lastName"`
	Role      string    `bson:"role" json:"role"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate input
	if loginReq.Username == "" || loginReq.Password == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Username and password are required",
		})
		return
	}

	// Get user from database
	user, err := database.GetUserByUsername(loginReq.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.SendJSONResponse(w, http.StatusUnauthorized, Response{
				Success: false,
				Message: "Invalid credentials",
			})
			return
		}
		utils.SendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error processing request",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		utils.SendJSONResponse(w, http.StatusUnauthorized, Response{
			Success: false,
			Message: "Invalid credentials",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Username, user.Role)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error generating token",
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, Response{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
				"fullName": user.FirstName + " " + user.LastName,
			},
		},
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var registerReq RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate input
	if err := validateRegisterRequest(registerReq); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Check if username already exists
	existingUser, err := database.GetUserByUsername(registerReq.Username)
	if err == nil && existingUser.Username != "" {
		utils.SendJSONResponse(w, http.StatusConflict, Response{
			Success: false,
			Message: "Username already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error processing registration",
		})
		return
	}

	// Create new user
	newUser := database.User{ // Use the User type from database package
		Username:  registerReq.Username,
		Password:  string(hashedPassword),
		Email:     registerReq.Email,
		FirstName: registerReq.FirstName,
		LastName:  registerReq.LastName,
		Role:      registerReq.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// Save to database
	err = database.CreateUser(newUser)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, Response{
			Success: false,
			Message: "Error creating user",
		})
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, Response{
		Success: true,
		Message: "User registered successfully",
	})
}

func validateRegisterRequest(req RegisterRequest) error {
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if len(req.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if req.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if req.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if req.Role == "" {
		return fmt.Errorf("role is required")
	}

	validRoles := map[string]bool{
		"admin":   true,
		"doctor":  true,
		"nurse":   true,
		"patient": true,
	}
	if !validRoles[req.Role] {
		return fmt.Errorf("invalid role specified")
	}
	return nil
}
