package database

import (
	"byte-go-mall/constant/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL(dbConfig *config.MySQLConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dbConfig.DSN()), &gorm.Config{
		Logger: newLogger(dbConfig),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 应用连接池配置
	if dbConfig.Pool.Enable {
		sqlDB.SetMaxIdleConns(dbConfig.Pool.MaxIdle)
		sqlDB.SetMaxOpenConns(dbConfig.Pool.MaxOpen)
		sqlDB.SetConnMaxIdleTime(dbConfig.Pool.MaxIdleTime)
		sqlDB.SetConnMaxLifetime(dbConfig.Pool.MaxLifeTime)
	}

	if dbConfig.Trace {
		err = SetupTracing(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
