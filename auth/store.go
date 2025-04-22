package auth

import (
	"database/sql"
	"fmt"

	"github.com/Tutuacs/pkg/db"
	"github.com/Tutuacs/pkg/types"
)

type Store struct {
	db    *sql.DB
	Table string
}

func NewStore(conn *sql.DB) (s *Store, err error) {
	con, err := db.NewConnection()
	if err != nil {
		return
	}

	s = &Store{db: con, Table: "users"}

	return
}

func (s *Store) GetUserByEmail(email string) (usr *types.User, err error) {
	err = nil
	usr = &types.User{}

	query := "SELECT * FROM " + s.Table + " WHERE email = $1"
	row := s.db.QueryRow(query, email)

	db.ScanRow(row, usr)

	if usr.ID == 0 {
		err = fmt.Errorf("user not found")
		return
	}

	return
}

func (s *Store) GetUserByID(ID int) (*types.User, error) {

	sql := "SELECT * FROM users WHERE id = $1"

	rows, err := s.db.Query(sql, ID)
	if err != nil {
		return nil, err
	}

	usr := new(types.User)

	for rows.Next() {
		err = db.ScanRows(rows, usr)
		if err != nil {
			return nil, err
		}
	}

	return usr, err
}

func (s *Store) CreateUser(user types.User) error {
	query := "INSERT INTO " + s.Table + " (name, email, password) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(query, user.Name, user.Email, user.Password)
	return err
}

func (s *Store) GetLogin(email string) (usr *types.User, err error) {

	sql := "SELECT * FROM " + s.Table + " WHERE email = $1"

	rows, err := s.db.Query(sql, email)
	if err != nil {
		return nil, err
	}

	usr = new(types.User)

	for rows.Next() {
		err = db.ScanRows(rows, usr)
		if err != nil {
			return
		}
	}

	return
}

func (s *Store) UpdatePassword(email, password string) error {
	query := "UPDATE " + s.Table + " SET password = $1 WHERE email = $2"
	_, err := s.db.Exec(query, password, email)
	return err
}
