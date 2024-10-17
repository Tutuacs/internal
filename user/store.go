package user

import (
	"database/sql"

	"github.com/Tutuacs/pkg/db"
	"github.com/Tutuacs/pkg/types"
)

type Store struct {
	db.Store
	db      *sql.DB
	extends bool
	Table   string
}

func NewStore(conn ...*sql.DB) (*Store, error) {
	if len(conn) == 0 {

		con, err := db.NewConnection()

		return &Store{
			db:      con,
			extends: false,
			Table:   "users",
		}, err
	}

	return &Store{
		db:      conn[0],
		extends: true,
		Table:   "users",
	}, nil
}

func (s *Store) CloseStore() {
	if !s.extends {
		s.db.Close()
	}

	// db.ScanRow()
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

func (s *Store) Update(id int64, newUser types.UpdateUserDto) (usr *types.User, err error) {

	sql := "UPDATE " + s.Table + " SET name = $1, role = $2, email = $3 WHERE id = $4 RETURNING id, name, role, email, createdAt"

	rows, err := s.db.Query(sql, newUser.Name, newUser.Role, newUser.Email, id)
	if err != nil {
		return
	}

	db.ScanRows(rows, usr)

	return
}

func (s *Store) Delete(id int64) (usr *types.User, err error) {

	sql := "DELETE FROM " + s.Table + " WHERE id = $1 RETURNING id, name, role, email, createdAt"

	row, err := s.db.Query(sql, id)
	if err != nil {
		return
	}

	db.ScanRows(row, usr)

	return
}
