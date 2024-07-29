package database

import (
	"context"
	"os"
	"sync"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Database struct {
	pool *pgxpool.Pool
}

var (
	instance *Database
	once     sync.Once
	logger    *zap.Logger
)

func SetLogger(l *zap.Logger) {
    logger = l
}

func NewDatabase() *Database {
	once.Do(func() {
		// Print the current working directory
		dir, err := os.Getwd()
		if err != nil {
			logger.Fatal("Error getting current working directory", zap.Error(err))
		}
		logger.Info("Current working directory", zap.String("directory", dir))

		// Print the absolute path to the .env file
		envPath, err := filepath.Abs(".env")
		if err != nil {
			logger.Fatal("Error getting absolute path to .env file", zap.Error(err))
		}
		logger.Info("Absolute path to .env file", zap.String("envPath", envPath))

		err = godotenv.Load()
		if err != nil {
			logger.Fatal("Error loading .env file", zap.Error(err))
		}

		instance = &Database{}
		instance.InitPool()
	})
	return instance
}

func (d *Database) InitPool() {
	connStr := "user=" + os.Getenv("PGUSER") + " password=" + os.Getenv("PGPASSWORD") + " dbname=" + os.Getenv("PGDATABASE") + " host=" + os.Getenv("PGHOST") + " sslmode=disable"

	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		logger.Error("Error creating database connection pool", zap.Error(err))
		return
	}

	d.pool = pool
	logger.Info("Connection to DB is created")
}

func (d *Database) CloseConnection() {
	if d.pool != nil {
		d.pool.Close()
		logger.Info("Connection to DB is closed")
		d.pool = nil
	} else {
		logger.Warn("Trying to close a non-open connection")
	}
}

func (d *Database) GetData(ctx context.Context, query string) ([]map[string]interface{}, error) {
	rows, err := d.pool.Query(ctx, query)
	if err != nil {
		logger.Error("Error while fetching list", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	fieldDescriptions := rows.FieldDescriptions()
	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		values := make([]interface{}, len(fieldDescriptions))
		valuePointers := make([]interface{}, len(fieldDescriptions))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			logger.Error("Error while scanning row", zap.Error(err))
			return nil, err
		}

		row := make(map[string]interface{})
		for i, fd := range fieldDescriptions {
			row[string(fd.Name)] = values[i]
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error while iterating over rows", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (d *Database) SetData(ctx context.Context, query string) error {
	_, err := d.pool.Exec(ctx, query)
	if err != nil {
		logger.Error("Error while executing query", zap.Error(err))
		return err
	}

	return nil
}

func (d *Database) GetList(ctx context.Context, query string) ([]map[string]interface{}, error) {
	return d.GetData(ctx, query)
}

func (d *Database) GetFetchval(ctx context.Context, query string) (interface{}, error) {
	var result interface{}
	err := d.pool.QueryRow(ctx, query).Scan(&result)
	if err != nil {
		logger.Error("Error while fetching value", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (d *Database) getRowCount(ctx context.Context, tableName string) (int64, error) {
    var count int64
    err := d.pool.QueryRow(ctx,"SELECT COUNT(*) FROM " + tableName).Scan(&count)
    if err != nil {
        return 0, err
    }
    return count, nil
}
