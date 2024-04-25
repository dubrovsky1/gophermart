package app

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	Host              string
	ConnectionString  string
	AccrualAddress    string
	MigrateVersion    int64
	ConnectionTimeout time.Duration
}

func NewConfig(conn time.Duration, version int64) *Config {
	return &Config{
		ConnectionTimeout: conn,
		MigrateVersion:    version,
	}
}

func (c *Config) ParseFlags() error {
	a := flag.String("a", "localhost:8080", "address and port to run server")
	d := flag.String("d", "host=localhost port=5432 user=sa password=admin dbname=gophermart sslmode=disable", "database connection string")
	f := flag.String("f", "accrual_windows_amd64.exe", "accrual system address")

	flag.Parse()

	runAddr := *a
	if ra := os.Getenv("RUN_ADDRESS"); ra != "" {
		runAddr = ra
	}

	connString := *d
	if cn := os.Getenv("DATABASE_URI"); cn != "" {
		connString = cn
	}

	accrualAddr := *f
	if aa := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); aa != "" {
		accrualAddr = aa
	}

	c.Host = runAddr
	c.ConnectionString = connString
	c.AccrualAddress = accrualAddr

	return nil
}
