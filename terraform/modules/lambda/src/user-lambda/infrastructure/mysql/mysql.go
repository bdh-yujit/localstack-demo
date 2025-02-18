package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/bootsdigitalhealth/go-db/sql"
	"github.com/go-sql-driver/mysql"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	readerConn *gorm.DB
	writerConn *gorm.DB
)

func InitDBConn() error {

	c := mysql.Config{
		DBName:                  os.Getenv("DB_NAME"),
		User:                    os.Getenv("DB_USER"),
		Passwd:                  os.Getenv("DB_PASSWORD"),
		Net:                     "tcp",
		ParseTime:               true,
		Collation:               "utf8_general_ci",
		Loc:                     time.UTC,
		AllowNativePasswords:    true,
		AllowCleartextPasswords: true,
		TLSConfig:               "false",
	}

	readerConf := c
	readerConf.Addr = fmt.Sprintf("%s:%s", os.Getenv("DB_READER_HOST"), os.Getenv("DB_PORT"))
	writerConf := c
	writerConf.Addr = fmt.Sprintf("%s:%s", os.Getenv("DB_WRITER_HOST"), os.Getenv("DB_PORT"))

	var logLevel logger.LogLevel
	if os.Getenv("STAGE") == "production" {
		logLevel = logger.Warn
	} else {
		logLevel = logger.Info
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,  // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false, // Don't include params in the SQL log
			Colorful:                  false, // Disable color
		},
	)

	readerDb, err := sql.Open("mysql", readerConf.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to open reader db: %w", err)
	}
	readerDb.SetMaxOpenConns(20)
	readerDb.SetMaxIdleConns(30)
	readerDb.SetConnMaxLifetime(20 * time.Second)
	readerGormDB, err := gorm.Open(gorm_mysql.New(gorm_mysql.Config{
		Conn: readerDb,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to open reader gorm db: %w", err)
	}

	writerDb, err := sql.Open("mysql", writerConf.FormatDSN())
	if err != nil {
		return fmt.Errorf("failed to open writer db: %w", err)
	}
	writerDb.SetMaxOpenConns(20)
	writerDb.SetMaxIdleConns(30)
	writerDb.SetConnMaxLifetime(20 * time.Second)
	writerGormDB, err := gorm.Open(gorm_mysql.New(gorm_mysql.Config{
		Conn: writerDb,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to open writer gorm db: %w", err)
	}

	readerConn = readerGormDB
	writerConn = writerGormDB

	return nil
}
