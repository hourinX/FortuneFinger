package gormc

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	oracledriver "github.com/godoes/gorm-oracle"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	mysqlSessionTimeZone    = "'+00:00'"
	postgresSessionTimeZone = "UTC"
	sqlDriverOracle         = "oracle"
	sqlDriverPostgres       = "pgx"
	sqlDriverSQLServer      = "sqlserver"
)

func openConnection(
	ctx context.Context,
	sourceName string,
	conf Config,
) (*sql.DB, gorm.Dialector, error) {
	var (
		pool      *sql.DB
		dialector gorm.Dialector
		err       error
	)

	switch conf.Driver {
	case DriverMySQL:
		pool, err = sql.Open(DriverMySQL, mysqlDSN(conf))
		if err == nil {
			dialector = mysql.New(mysql.Config{Conn: pool})
		}
	case DriverOracle:
		pool, err = sql.Open(sqlDriverOracle, oracleDSN(conf))
		if err == nil {
			dialector = oracledriver.New(oracledriver.Config{Conn: pool})
		}
	case DriverPostgres:
		pool, err = sql.Open(sqlDriverPostgres, postgresDSN(conf))
		if err == nil {
			dialector = postgres.New(postgres.Config{Conn: pool})
		}
	case DriverSQLServer:
		pool, err = sql.Open(sqlDriverSQLServer, sqlServerDSN(conf))
		if err == nil {
			dialector = sqlserver.New(sqlserver.Config{Conn: pool})
		}
	default:
		return nil, nil, fmt.Errorf(
			"gormc: 暂不支持数据源「%s」的驱动「%s」",
			sourceName,
			conf.Driver,
		)
	}
	if err != nil {
		return nil, nil, fmt.Errorf(
			"gormc: 创建数据源「%s」的「%s」数据库连接池失败: %w",
			sourceName,
			conf.Driver,
			err,
		)
	}

	pool.SetConnMaxIdleTime(conf.ConnMaxIdleTime)
	pool.SetConnMaxLifetime(conf.ConnMaxLifetime)
	pool.SetMaxIdleConns(conf.MaxIdleConns)
	pool.SetMaxOpenConns(conf.MaxOpenConns)
	if err := pool.PingContext(ctx); err != nil {
		_ = pool.Close()
		return nil, nil, fmt.Errorf(
			"gormc: 数据源「%s」的「%s」数据库ping失败: %w",
			sourceName,
			conf.Driver,
			err,
		)
	}

	return pool, dialector, nil
}

func mysqlDSN(conf Config) string {
	driverConfig := mysqldriver.NewConfig()
	driverConfig.User = conf.User
	driverConfig.Passwd = conf.Password
	driverConfig.Net = "tcp"
	driverConfig.Addr = net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))
	driverConfig.DBName = conf.DBName
	driverConfig.ParseTime = true
	driverConfig.Loc = time.UTC
	driverConfig.Params = map[string]string{
		"charset":   "utf8mb4",
		"time_zone": mysqlSessionTimeZone,
	}
	if conf.SSLMode != "" && conf.SSLMode != "disable" {
		driverConfig.TLSConfig = conf.SSLMode
	}
	return driverConfig.FormatDSN()
}

func oracleDSN(conf Config) string {
	return oracledriver.BuildUrl(
		conf.Host,
		conf.Port,
		conf.ServiceName,
		conf.User,
		conf.Password,
		nil,
	)
}

func postgresDSN(conf Config) string {
	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(conf.User, conf.Password),
		Host:   net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)),
		Path:   conf.DBName,
	}
	query := dsn.Query()
	query.Set("sslmode", sslModeOrDefault(conf.SSLMode))
	query.Set("TimeZone", postgresSessionTimeZone)
	dsn.RawQuery = query.Encode()
	return dsn.String()
}

func sqlServerDSN(conf Config) string {
	dsn := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(conf.User, conf.Password),
		Host:   net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port)),
	}
	query := dsn.Query()
	query.Set("database", conf.DBName)
	if conf.SSLMode != "" {
		query.Set("encrypt", conf.SSLMode)
	}
	dsn.RawQuery = query.Encode()
	return dsn.String()
}

func sslModeOrDefault(sslMode string) string {
	if sslMode == "" {
		return "disable"
	}
	return sslMode
}
