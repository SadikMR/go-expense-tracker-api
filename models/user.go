package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SadikMR/go-expense-tracker-api/utils"
)

var userCSVHeader = []string{"id", "name", "email", "password", "created_at"}

func usersCSVPath() string {
	if utils.AppConfig.UsersCSVPath != "" {
		return utils.AppConfig.UsersCSVPath
	}
	return filepath.Join(dataDir(), "users.csv")
}

func dataDir() string {
	if utils.AppConfig.DataDir != "" {
		return utils.AppConfig.DataDir
	}
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		return dir
	}
	return "data"
}

// User represents a registered user in the system.
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"` // never serialised to JSON
	CreatedAt string `json:"created_at"`
}

// EnsureUsersCSV creates the users CSV file with its header if it does not exist.
func EnsureUsersCSV() error {
	return utils.EnsureCSVExists(usersCSVPath(), userCSVHeader)
}

// GetAllUsers returns every user record from the CSV file.
func GetAllUsers() ([]User, error) {
	rows, err := utils.ReadCSV(usersCSVPath())
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: %w", err)
	}

	users := make([]User, 0, len(rows))
	for _, row := range rows {
		u, err := rowToUser(row)
		if err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUserByEmail finds a user by email address.
// Returns nil, nil when the user does not exist.
func GetUserByEmail(email string) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("GetUserByEmail: %w", err)
	}
	for _, u := range users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, nil
}

// GetUserByID finds a user by their integer ID.
// Returns nil, nil when the user does not exist.
func GetUserByID(id int) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}
	for _, u := range users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, nil
}

// CreateUser appends a new user record to the CSV file.
func CreateUser(u *User) error {
	u.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := utils.AppendCSV(usersCSVPath(), userToRow(u)); err != nil {
		return fmt.Errorf("CreateUser: %w", err)
	}
	return nil
}

// NextUserID returns the next available user ID.
func NextUserID() (int, error) {
	users, err := GetAllUsers()
	if err != nil {
		return 0, fmt.Errorf("NextUserID: %w", err)
	}
	max := 0
	for _, u := range users {
		if u.ID > max {
			max = u.ID
		}
	}
	return max + 1, nil
}

func rowToUser(row []string) (User, error) {
	if len(row) < 5 {
		return User{}, fmt.Errorf("rowToUser: expected 5 fields, got %d", len(row))
	}
	id, err := strconv.Atoi(row[0])
	if err != nil {
		return User{}, fmt.Errorf("rowToUser: invalid id %q: %w", row[0], err)
	}
	return User{
		ID:        id,
		Name:      row[1],
		Email:     row[2],
		Password:  row[3],
		CreatedAt: row[4],
	}, nil
}

func userToRow(u *User) []string {
	return []string{
		strconv.Itoa(u.ID),
		u.Name,
		u.Email,
		u.Password,
		u.CreatedAt,
	}
}
