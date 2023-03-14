package test

import (
	"2f-authorization/initiator"
	"2f-authorization/internal/constants/dbinstance"
	"2f-authorization/internal/handler/middleware"
	"2f-authorization/platform/logger"
	"2f-authorization/platform/opa"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type TestInstance struct {
	Server *gin.Engine
	DB     dbinstance.DBInstance
	Opa    opa.Opa
}

func Initiate(ctx context.Context, path string) TestInstance {
	log := logger.New(initiator.InitLogger())
	log.Info(context.Background(), "logger initialized")

	log.Info(context.Background(), "initializing config")
	configName := "config"
	if name := os.Getenv("CONFIG_NAME"); name != "" {
		configName = name
		log.Info(context.Background(), fmt.Sprintf("config name is set to %s", configName))
	} else {
		log.Info(context.Background(), "using default config name 'config'")
	}
	initiator.InitConfig(configName, path+"config", log)
	log.Info(context.Background(), "config initialized")

	log.Info(context.Background(), "initializing database")
	Conn := initiator.InitDB(viper.GetString("database.url"), viper.GetDuration("database.idle_conn_timeout"), log)
	log.Info(context.Background(), "database initialized")

	log.Info(context.Background(), "initializing migration")
	m := initiator.InitiateMigration(path+viper.GetString("migration.path"), path+viper.GetString("database.url"), log)
	initiator.UpMigration(m, log)
	log.Info(context.Background(), "migration initialized")

	log.Info(context.Background(), "initializing persistence layer")
	dbConn := dbinstance.New(Conn)
	persistence := initiator.InitPersistence(dbConn, log)
	log.Info(context.Background(), "persistence layer initialized")

	log.Info(context.Background(), "initializing opa")
	rand.Seed(time.Now().Unix())
	port := rand.Intn(1000) + 40000
	opa := initiator.InitOpa(ctx, path+viper.GetString("opa.path"), path+viper.GetString("opa.data_file"), path+viper.GetString("opa.server_exec"), persistence, port, log)
	log.Info(context.Background(), "opa initialized")

	log.Info(context.Background(), "initializing module")
	module := initiator.InitModule(persistence, log, opa)
	log.Info(context.Background(), "module initialized")

	log.Info(context.Background(), "initializing handler")
	handler := initiator.InitHandler(module, log)
	log.Info(context.Background(), "handler initialized")

	log.Info(context.Background(), "initializing server")
	server := gin.New()
	server.Use(middleware.GinLogger(log.Named("gin")))
	server.Use(ginzap.RecoveryWithZap(log.GetZapLogger().Named("gin.recovery"), true))
	server.Use(middleware.ErrorHandler())
	log.Info(context.Background(), "server initialized")

	log.Info(context.Background(), "initializing router")
	v1 := server.Group("/v1")
	initiator.InitRouter(v1, handler, persistence, log, opa)
	log.Info(context.Background(), "router initialized")

	return TestInstance{
		Server: server,
		DB:     dbConn,
		Opa:    opa,
	}
}

func (t *TestInstance) BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
