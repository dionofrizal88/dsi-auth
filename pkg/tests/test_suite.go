package tests

import (
	"github.com/dionofrizal88/dsi/auth/config"
	"github.com/dionofrizal88/dsi/auth/infrastructure/dao"
	"github.com/dionofrizal88/dsi/auth/pkg/db"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// TestSuite is a struct represent is self.
type TestSuite struct {
	Config      config.Configuration
	DBService   *dao.Repositories
	RedisClient *redis.Client
}

// InitTestSuite knows how to initialize test suite.
func InitTestSuite() *TestSuite {
	conf := config.GetConfig("../../config/config.json")
	if conf.AppName == "" {
		conf = config.GetConfig("../config/config.json")
	}

	dbConnection, errDBConnection := db.NewPostgresDBTestConnection(conf)
	if errDBConnection != nil {
		log.Fatalf("Failed to connect into postgres db %v", errDBConnection)
	}

	redisConnection, errRedisConnection := db.NewRedisConnection(conf)
	if errRedisConnection != nil {
		log.Fatalf("Failed to connect into redis db %v", errRedisConnection)
	}

	dbService := dao.NewDBService(dbConnection)
	entityRegistry := dao.NewRegistry()
	errAutoReset := entityRegistry.ResetDatabase(dbConnection)
	if errAutoReset != nil {
		panic(errAutoReset)
	}

	time.Sleep(4 * time.Second)

	// Auto Migrate
	errAutoMigrate := entityRegistry.AutoMigrate(dbConnection)
	if errAutoMigrate != nil {
		log.Fatalf("Failed to run auto migration: %v", errAutoMigrate)
	}

	return &TestSuite{
		Config:      conf,
		DBService:   dbService,
		RedisClient: redisConnection,
	}
}

// Clean is a method uses to close DB connection.
func (t *TestSuite) Clean() {
	dbClean, err := t.DBService.DB.DB()
	if err != nil {
		return
	}

	err = dbClean.Close()
	if err != nil {
		return
	}
}
