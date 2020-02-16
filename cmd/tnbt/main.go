package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/LuD1161/restructuring-tnbt/pkg/database/postgres"
	"github.com/LuD1161/restructuring-tnbt/pkg/middlewares/auth"
	"github.com/LuD1161/restructuring-tnbt/pkg/user"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	// _ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/LuD1161/restructuring-tnbt/cmd/tnbt/docs"
)

// @title TNBT Swagger API
// @version 1.0
// @description Swagger API for TNBT
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email aseemshrey@gmail.com

// @license.name MIT
// @license.url https://github.com/github.com/LuD1161/restructuring-tnbt

// @BasePath /
func main() {
	var server bool
	var dbType, dbURL, serverPort, logLevel string
	var err error

	flag.StringVar(&dbType, "database", "postgres", "database type [redis, postgres]")
	flag.StringVar(&serverPort, "port", "3000", "server port to run on")
	flag.StringVar(&logLevel, "loglevel", "debug", "log level [trace, debug, info, warn, error, fatal, panic], default : debug")
	flag.BoolVar(&server, "server", false, "run in server mode")
	flag.Parse()

	log := NewLogger(logLevel)
	err = godotenv.Load()
	// Todo : Maybe integrate vault here
	dbURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	if err != nil {
		logrus.Fatalf("Error getting env, not comint through %v", err)
	} else {
		logrus.Info("We are getting the env values")
	}

	var userRepo user.Repository

	switch dbType {
	case "postgres":
		pconn := postgresConnection(dbURL)
		defer pconn.Close()
		userRepo = postgres.NewPostgresUserRepository(pconn)
		seedData(pconn)
	// case "redis":
	// 	dbURL = env.EnvString("DATABASE_URL", DefaultRedisUrl)
	// 	redisPassword = env.EnvString("REDIS_PASSWORD", DefaultRedisPassword)
	// 	rconn := redisConnect(dbURL, redisPassword)
	// 	defer rconn.Close()
	// 	userRepo = redisdb.NewRedisTicketRepository(rconn)
	default:
		panic("Unknown database")
	}

	userService := user.NewService(userRepo, log)
	userHandler := user.NewHandler(userService, log)

	router := gin.Default()

	url := ginSwagger.URL("http://localhost:" + serverPort + "/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.POST("/login", userHandler.Login)
	router.POST("/user", userHandler.CreateUser)

	authorized := router.Group("/")
	authorized.Use(auth.SetMiddleWareAuthentication())
	authorized.GET("/user/:id", userHandler.GetUserByID)
	authorized.PUT("/user", userHandler.UpdateUser)

	// http.Handle("/", accessControl(middleware.Authenticate(router)))

	errs := make(chan error, 2)
	var docker = "false"
	go func() {
		if server || docker == "true" {
			logrus.Info("Listening server mode on port :3001")
			errs <- http.ListenAndServe(":3001", nil)
		} else {
			logrus.Info("Listening locally on port :" + serverPort)
			errs <- router.Run(":" + serverPort)
		}
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logrus.Errorf("terminated %s", <-errs)
}

func postgresConnection(database string) *gorm.DB {
	logrus.Info("Connecting to PostgreSQL DB")
	db, err := gorm.Open("postgres", database)
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
	return db
}

// NewLogger : Logger for injecting into all the dependencies
func NewLogger(logLevel string) *logrus.Logger {
	var level logrus.Level
	switch logLevel {
	case "trace":
		level = logrus.TraceLevel
	case "info":
		level = logrus.InfoLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		level = logrus.DebugLevel
	}
	return &logrus.Logger{
		Out: os.Stderr,
		Formatter: &runtime.Formatter{
			ChildFormatter: &logrus.JSONFormatter{},
			Package:        true,
			Line:           true,
			File:           true,
		},
		Level: level,
	}
}

func seedData(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&user.User{}).Error
	if err != nil {
		logrus.Fatalf("cannot drop table: %v", err)
	}
	var users = []user.User{
		user.User{
			UserInfoPayload: user.UserInfoPayload{
				Username: "user1",
				Email:    "user1@gmail.com",
			},
			Password: "$2a$10$q3wpVNrAkMDDB9QYt1Ym7.E6CCIscJL5kdsgyNkrKqcxDV4VAvaYS",
		},
		user.User{
			UserInfoPayload: user.UserInfoPayload{
				Username: "user2",
				Email:    "user2@gmail.com",
			},
			Password: "$2a$10$LRazVeYPrx6MOYCZAGrSZ.CCD7d3qYUn4T5CbqZSZ37Pe28pvaBkq",
		},
	}
	err = db.Debug().AutoMigrate(&user.User{}).Error
	if err != nil {
		logrus.Fatalf("cannot migrate table: %v", err)
	}
	for i := range users {
		err = db.Debug().Model(&user.User{}).Create(&users[i]).Error
		if err != nil {
			logrus.Fatalf("cannot seed users table: %v", err)
		}
	}
}
