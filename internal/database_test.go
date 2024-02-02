package internal

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

type mockDB struct {
	connect func(driverName, dataSourceName string) (*sqlx.DB, error)
}

func (m *mockDB) Connect(driverName, dataSourceName string) (*sqlx.DB, error) {
	return m.connect(driverName, dataSourceName)
}

func TestConnectDatabase(t *testing.T) {
	// Define test cases
	cases := []struct {
		name     string
		db       mockDB
		database Database
		wantErr  bool
	}{
		{
			name: "success case",
			db: mockDB{
				connect: func(driverName, dataSourceName string) (*sqlx.DB, error) {
					return nil, nil
				},
			},
			database: Database{
				Connection: &Connection{
					Username: "testUser",
					Password: "testPassword",
					Hostname: "localhost",
					Port:     3306,
				},
				Name: "testDB",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			db: mockDB{
				connect: func(driverName, dataSourceName string) (*sqlx.DB, error) {
					return nil, fmt.Errorf("error connecting to db")
				},
			},
			database: Database{
				Connection: &Connection{
					Username: "testUser",
					Password: "testPassword",
					Hostname: "localhost",
					Port:     3306,
				},
				Name: "testDB",
			},
			wantErr: true,
		},
	}

	// Execute test cases
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.db.Connect("mysql", fmt.Sprintf(
				"%s:%s@(%s:%d)/%s?parseTime=true",
				tc.database.Connection.Username,
				tc.database.Connection.Password,
				tc.database.Connection.Hostname,
				tc.database.Connection.Port,
				tc.database.Name,
			))

			if (err != nil) != tc.wantErr {
				t.Errorf("ConnectDatabase() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
