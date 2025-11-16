package gateway

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	slog_gorm "github.com/orandin/slog-gorm"
	gorm_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/mocoarow/cocotola-1.25/moonbeam/lib/domain"
)

type MySQLConfig struct {
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	Database string `yaml:"database" validate:"required"`
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func openMySQLForTest() (*gorm.DB, error) {
	host := getEnv("TEST_MYSQL_HOST", "127.0.0.1")
	portStr := getEnv("TEST_MYSQL_PORT", "3307")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid TEST_MYSQL_PORT: %w", err)
	}
	username := getEnv("TEST_MYSQL_USERNAME", "username")
	password := getEnv("TEST_MYSQL_PASSWORD", "password")
	database := getEnv("TEST_MYSQL_DATABASE", "test")
	logLevelStr := getEnv("TEST_LOG_LEVEL", "INFO")
	logLevel := slog.LevelInfo
	if logLevelStr == "DEBUG" {
		logLevel = slog.LevelDebug
	}

	db, err := OpenMySQL(&MySQLConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}, logLevel, "test")
	if err != nil {
		return nil, fmt.Errorf("open test MySQL: %w", err)
	}

	return db, nil
}

// func setupMySQL(sqlFS embed.FS, db *gorm.DB) error {
// 	driverName := "mysql"
// 	sourceDriver, err := iofs.New(sqlFS, driverName)
// 	if err != nil {
// 		return err
// 	}

// 	return setupDB(db, driverName, sourceDriver, func(sqlDB *sql.DB) (database.Driver, error) {
// 		return migrate_mysql.WithInstance(sqlDB, &migrate_mysql.Config{})
// 	})
// }

// func InitMySQL(fs embed.FS, dbHost string, dbPort int) (*gorm.DB, error) {
// 	testDBHost = dbHost
// 	testDBPort = dbPort
// 	db, err := openMySQLForTest()
// 	if err != nil {
// 		return nil, mbliberrors.Errorf("openMySQLForTest: %w", err)
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, mbliberrors.Errorf("DB: %w", err)
// 	}

// 	if err := sqlDB.Ping(); err != nil {
// 		return nil, mbliberrors.Errorf("Ping: %w", err)
// 	}

// 	if err := libgateway.MigrateMySQLDB(db, fs); err != nil {
// 		return nil, mbliberrors.Errorf("MigrateMySQLDB: %w", err)
// 	}

// 	return db, nil
// }

// func InitMySQLWithDSN(fs embed.FS, dsn string) (*gorm.DB, error) {
// 	testDSN = dsn
// 	db, err := openMySQLForTest()
// 	if err != nil {
// 		return nil, mbliberrors.Errorf("openMySQLForTest: %w", err)
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, mbliberrors.Errorf("DB: %w", err)
// 	}

// 	if err := sqlDB.Ping(); err != nil {
// 		return nil, mbliberrors.Errorf("Ping: %w", err)
// 	}

// 	if err := libgateway.MigrateMySQLDB(db, fs); err != nil {
// 		return nil, mbliberrors.Errorf("MigrateMySQLDB: %w", err)
// 	}

// 	return db, nil
// }

func OpenMySQLWithDSN(dsn string, logLevel slog.Level, appName string) (*gorm.DB, error) {
	gormDialector := gorm_mysql.Open(dsn)

	options := make([]slog_gorm.Option, 0)
	options = append(options, slog_gorm.WithHandler(slog.Default().With(slog.String(domain.LoggerNameKey, appName+"-gorm")).Handler()))
	if logLevel == slog.LevelDebug {
		options = append(options, slog_gorm.WithTraceAll()) // trace all messages
	}

	gormConfig := gorm.Config{ //nolint:exhaustruct
		Logger: slog_gorm.New(options...),
	}

	db, err := gorm.Open(gormDialector, &gormConfig)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	return db, nil
}

func OpenMySQL(cfg *MySQLConfig, logLevel slog.Level, appName string) (*gorm.DB, error) {
	c := mysql.Config{ //nolint:exhaustruct
		DBName:               cfg.Database,
		User:                 cfg.Username,
		Passwd:               cfg.Password,
		Addr:                 fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Net:                  "tcp",
		ParseTime:            true,
		MultiStatements:      true,
		Params:               map[string]string{"charset": "utf8mb4"},
		Collation:            "utf8mb4_bin",
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
		MaxAllowedPacket:     64 << 20, // 64 MiB.
		Loc:                  time.UTC,
	}

	return OpenMySQLWithDSN(c.FormatDSN(), logLevel, appName)
}
