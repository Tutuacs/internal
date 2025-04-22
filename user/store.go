package user

import (
	"database/sql"

	"github.com/Tutuacs/pkg/db"
	"github.com/Tutuacs/pkg/types"
)

type Store struct {
	db.Store
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

func (s *Store) GetConn() *sql.DB {

	return s.db
}

// TODO: Implement the store consults

func (s *Store) Create(newUser types.NewUserDto) (usr *types.User, err error) {
	sql := "INSERT INTO " + s.Table + " (name, role, email, password) VALUES ($1, $2, $3, $4) RETURNING id, name, role, email, password, createdAt"

	row := s.db.QueryRow(sql, newUser.Name, newUser.Role, newUser.Email, newUser.Password)

	usr = &types.User{}
	err = db.ScanRow(row, usr)

	return
}

func (s *Store) GetByID(ID int64) (*types.User, error) {

	sql := "SELECT * FROM " + s.Table + " WHERE id = $1"

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

func (s *Store) GetByEmail(Email string) (*types.User, error) {

	sql := "SELECT * FROM " + s.Table + " WHERE email = $1"

	rows, err := s.db.Query(sql, Email)
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

func (s *Store) List() (usrs []*types.User, err error) {

	sql := "SELECT * FROM " + s.Table

	rows, err := s.db.Query(sql)
	if err != nil {
		return
	}

	for rows.Next() {
		usr := new(types.User)
		err = db.ScanRows(rows, usr)

		if err != nil {
			continue
		}

		if usr.ID != 0 {
			usrs = append(usrs, usr)
		}
	}

	return
}

func (s *Store) Update(id int64, newUser types.UpdateUserDto) (usr *types.UpdateUserDto, err error) {

	sql := "UPDATE " + s.Table + " SET name = $1, role = $2, email = $3 WHERE id = $6 RETURNING name, email, role"

	rows, err := s.db.Query(sql, newUser.Name, newUser.Role, newUser.Email, id)
	if err != nil {
		return
	}

	usr = &types.UpdateUserDto{}
	for rows.Next() {
		err = db.ScanRows(rows, usr)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (s *Store) Delete(id int64) (usr *types.UpdateUserDto, err error) {

	sql := "DELETE FROM " + s.Table + " WHERE id = $1 RETURNING id, name, role, email, createdAt"

	rows, err := s.db.Query(sql, id)
	if err != nil {
		return
	}

	usr = &types.UpdateUserDto{}
	for rows.Next() {
		err = db.ScanRows(rows, usr)
		if err != nil {
			return nil, err
		}
	}

	return
}
