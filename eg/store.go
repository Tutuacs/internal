package eg

import (
	"database/sql"

	"github.com/Tutuacs/pkg/db"
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

		db.NewConnection()

		return &Store{
			db:      con,
			extends: false,
		}, err
	}

	return &Store{
		db:      conn[0],
		extends: true,
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

func (s *Store) Create() (eg Eg, err error) {
	return
}
