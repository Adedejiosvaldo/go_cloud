package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Database Relational Storage abstraction
type Database struct {
	DB                    *sqlx.DB
	host                  string
	port                  int
	user                  string
	password              string
	name                  string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
	log                   *zap.Logger
}

// for newDatabaseOptions for NewDatabase
type NewDatabaseOptions struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	MaxIdleConnections    int
	MaxOpenConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
	Log                   *zap.Logger
}

func NewDatabase(opt NewDatabaseOptions) *Database {
	if opt.Log == nil {
		opt.Log = zap.NewNop()
	}

	return &Database{
		host:                  opt.Host,
		port:                  opt.Port,
		user:                  opt.User,
		password:              opt.Password,
		name:                  opt.Name,
		maxOpenConnections:    opt.MaxOpenConnections,
		maxIdleConnections:    opt.MaxIdleConnections,
		connectionMaxIdleTime: opt.ConnectionMaxIdleTime,
		log:                   opt.Log,
	}
}
