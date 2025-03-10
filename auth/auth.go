package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"social-network/database"
	// "forum/utils"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

var sessionStore = make(map[string]string)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	email := EscapeString(strings.ToLower(r.FormValue("email")))
	// fmt.Println("Username: ", email)
	password := EscapeString(r.FormValue("password"))

	const maxEmail = 100
	const maxPassword = 100

	if len(email) > maxEmail {
		response := map[string]string{"error": fmt.Sprintf("Email cannot be longer than %d characters.", maxEmail)}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(password) > maxPassword {
		response := map[string]string{"error": fmt.Sprintf("Password cannot be longer than %d characters.", maxPassword)}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var storedPassword, sessionToken, username string
	err := database.DB.QueryRow("SELECT password, session_token, username FROM users WHERE email = ?", email).Scan(&storedPassword, &sessionToken, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]string{"error": "Invalid email or password"}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
		} else {
			log.Printf("Database error: %v", err)
			response := map[string]string{"error": "Internal server error"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)); err != nil {
		response := map[string]string{"error": "Invalid email or password"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Generate a new session token
	newSessionToken, _ := uuid.NewV4()
	sessionToken = newSessionToken.String()

	// Update the session token in the database
	_, err = database.DB.Exec("UPDATE users SET session_token = ? WHERE email = ?", sessionToken, email)
	if err != nil {
		log.Printf("Error updating session token: %v", err)
		response := map[string]string{"error": "Internal server error"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set a cookie with the session token
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(1 * time.Hour),
	})

	response := map[string]string{"message": "Login successful!", "username": username}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)

	if r.Method != http.MethodPost {
		response = map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	username := EscapeString(r.FormValue("username"))
	firstName := EscapeString(r.FormValue("firstName"))
	lastName := EscapeString(r.FormValue("lastName"))
	email := EscapeString(strings.ToLower(r.FormValue("email")))
	gender := EscapeString(r.FormValue("gender"))
	age := EscapeString(r.FormValue("age"))
	password := EscapeString(r.FormValue("password"))

	if gender == "" || (gender != "M" && gender != "F") {
		response["error"] = "Gender must be either 'M' or 'F'"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	ageInt, err := strconv.Atoi(age)
	if err != nil || ageInt < 1 || ageInt > 120 {
		response["error"] = "Age must be a number between 1 and 120"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if firstName == "" {
		response["error"] = "First Name is required"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if lastName == "" {
		response["error"] = "Last Name is required"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	errors, valid := ValidateInput(username, firstName, lastName, email, gender, age, password)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation error",
			"fields": errors,
		})
		return
	}

	var existingUsername, existingEmail string
	err = database.DB.QueryRow("SELECT username, email FROM users WHERE email = ? OR username = ?", email, username).Scan(&existingUsername, &existingEmail)
	if err == nil {
		var conflictField, conflictMessage string
		if existingUsername == username {
			conflictField = "username"
			conflictMessage = "Username already exists"
		} else if existingEmail == email {
			conflictField = "email"
			conflictMessage = "Email already exists"
		}

		response = map[string]string{"error": conflictMessage, "field": conflictField}
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		log.Printf("Database error: %v", err)
		response = map[string]string{"error": "Database error"}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response = map[string]string{"error": "Error hashing password"}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sessionToken, _ := uuid.NewV4()

	_, err = database.DB.Exec("INSERT INTO users (username, firstname, lastname, email, gender, age, password, session_token) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", firstName, lastName, username, email, gender, age, hashedPassword, sessionToken)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		response = map[string]string{"error": "Registration failed"}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response = map[string]string{"message": "Registration successful! Please log in."}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "You are not logged in", http.StatusBadRequest)
		return
	}

	// Remove the session from the session store
	delete(sessionStore, cookie.Value)

	// Expire the cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "guest",
		Expires: time.Now().Add(-1 * time.Hour),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Fprintln(w, "You have been logged out.")
}

func RequireLogin(w http.ResponseWriter, r *http.Request) (string, string, bool, error) {
	cookie, _ := r.Cookie("session_token")
	if cookie == nil {
		return "", "guest", false, nil
	}

	var username, sessionToken string
	err := database.DB.QueryRow("SELECT username, session_token FROM users WHERE session_token = ?", cookie.Value).Scan(&username, &sessionToken)
	if err == sql.ErrNoRows {
		cookies := r.Cookies()
		// Loop through the cookies and expire them
		for _, cookie := range cookies {
			http.SetCookie(w, &http.Cookie{
				Name:    cookie.Name,
				Value:   "",
				Expires: time.Now().Add(-1 * time.Hour),
				// MaxAge:  -1,
			})
		}
		return "", "guest", false, err
	} else if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return "", "guest", false, err
	}

	return username, sessionToken, true, nil
}

func CheckSessionHandler(w http.ResponseWriter, r *http.Request) {
	_, _, loggedIn, err := RequireLogin(w, r)
	// fmt.Println("sessionGuest1:", sessionGuest)
	if err != nil {
		fmt.Println("Error in the RequiredLogin !!! :", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if loggedIn {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"loggedIn": true}`)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"loggedIn": false}`)
	}
}

func ValidateInput(username, firstName, lastName, email, gender, age, password string) (map[string]string, bool) {
	errors := make(map[string]string)

	const maxUsername = 50
	const maxEmail = 100
	const maxPassword = 100

	if len(username) == 0 {
		errors["username"] = "Username cannot be empty"
		return errors, false
	} else if len(username) > maxUsername {
		errors["username"] = fmt.Sprintf("Username cannot be longer than %d characters.", maxUsername)
		return errors, false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if len(email) == 0 {
		errors["email"] = "Email cannot be empty"
		return errors, false
	} else if len(email) > maxEmail {
		errors["email"] = fmt.Sprintf("Email cannot be longer than %d characters.", maxEmail)
		return errors, false
	} else if !emailRegex.MatchString(email) {
		errors["email"] = "Invalid email format"
		return errors, false
	}

	if len(password) < 8 {
		errors["password"] = "Password must be at least 8 characters long"
		return errors, false
	} else if len(password) > maxPassword {
		errors["password"] = fmt.Sprintf("Password cannot be longer than %d characters.", maxPassword)
		return errors, false
	} else {
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)

		if !hasUpper {
			errors["password"] = "Password must include at least one uppercase letter"
			return errors, false
		} else if !hasLower {
			errors["password"] = "Password must include at least one lowercase letter"
			return errors, false
		} else if !hasDigit {
			errors["password"] = "Password must include at least one digit"
			return errors, false
		} else if !hasSpecial {
			errors["password"] = "Password must include at least one special character"
			return errors, false
		}
	}

	if len(errors) > 0 {
		log.Println(errors)
		return errors, false
	}
	return nil, true
}

func EscapeString(s string) string {
	return html.EscapeString(s)
}
