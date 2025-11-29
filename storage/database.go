package storage

import (
	"context"
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

func (d *Database) Connect() error {
	d.log.Info("Connecting to the database", zap.String("url", d.createDataSourceName(false)))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error

	d.DB, err = sqlx.Connect(ctx, "pgx", d.createDataSourceName(true))

	if err != nil {
		return err
	}
	d.log.Debug("Setting connection pool options",
		zap.Int("max open connections", d.maxOpenConnections),
		zap.Int("max idle connections", d.maxIdleConnections),
		zap.Duration("connection max lifetime", d.connectionMaxLifetime),
		zap.Duration("connection max idle time", d.connectionMaxIdleTime))

	d.DB.SetMaxOpenConns(d.maxOpenConnections)
	d.DB.SetMaxIdleConns(d.maxIdleConnections)
	d.DB.SetConnMaxLifetime(d.connectionMaxLifetime)
	d.DB.SetConnMaxIdleTime(d.connectionMaxIdleTime)

	return nil
}
