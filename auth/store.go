package auth

import (
	"database/sql"
	"fmt"

	"github.com/Tutuacs/pkg/db"
)

type Store struct {
	db      *sql.DB
	extends bool
	Table   string
}

func NewStore(conn ...*sql.DB) (*Store, error) {
	if len(conn) == 0 {
		con, err := db.NewConnection()
		if err != nil {
			return nil, err
		}
		return &Store{db: con, extends: false, Table: "users"}, nil
	}
	return &Store{db: conn[0], extends: true}, nil
}

func (s *Store) CloseStore() {
	if !s.extends {
		s.db.Close()
	}
}

func (s *Store) GetUserByEmail(email string) (*User, error) {
	var user User

	query := "SELECT * FROM " + s.Table + " WHERE email = ?"
	row := s.db.QueryRow(query, email)
	db.ScanRow(row, user)

	if user.ID == 0 {
		return nil, fmt.Errorf("User not found")
	}

	return &user, nil
}

func (s *Store) GetUserByID(ID int) (*User, error) {

	sql := "SELECT * FROM users WHERE id = ?"

	rows, err := s.db.Query(sql, ID)
	if err != nil {
		return nil, err
	}

	usr := new(User)

	for rows.Next() {
		err = db.ScanRows(rows, usr)
		if err != nil {
			return nil, err
		}
	}

	return usr, err
}

func (s *Store) CreateUser(user User) error {
	query := "INSERT INTO " + s.Table + " (firstName, lastName, email, password) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(query, user.Email, user.Password)
	return err
}

func (s *Store) GetLogin(email string) (int64, string, string, error) {
	var userEmail, password string
	var id int64

	query := fmt.Sprintf("SELECT email, password FROM %s WHERE email = ?", s.Table)
	err := s.db.QueryRow(query, email).Scan(&id, &userEmail, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", "", fmt.Errorf("user not found")
		}
		return 0, "", "", err
	}

	return id, userEmail, password, nil
}
