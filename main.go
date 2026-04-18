package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/huyshop/product/db"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

type Configs struct {
	GRPCPort        string
	DBPath          string
	DBName          string
	RedisAddr       string
	RedisPassword   string
	RedisDb         string
	RedisCartExpire string
	UserHost        string
}

var config *Configs

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: Error loading .env file:", err)
		} else {
			log.Println("Loaded .env file for local development")
		}
	} else {
		log.Println("No .env file found, using system environment variables")
	}

	config = &Configs{
		GRPCPort:        getEnv("GRPC_PORT", "8000"),
		DBPath:          getEnv("DB_PATH", ".root:123456@tcp(localhost:3306)"),
		DBName:          getEnv("DB_NAME", "product"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDb:         getEnv("REDIS_DB", "0"),
		RedisCartExpire: getEnv("REDIS_CART_EXPIRE", "3600"),
		UserHost:        getEnv("USER_HOST", "localhost:6001"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func startApp(ctx *cli.Context) error {
	log.Printf("Starting product service with config:")
	log.Printf("  GRPC Port: %s", config.GRPCPort)
	log.Printf("  DB Path: %s", config.DBPath)
	log.Printf("  DB Name: %s", config.DBName)
	log.Printf("  Redis Addr: %s", config.RedisAddr)
	v, err := NewProduct(config)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err := startGRPCServe(config.GRPCPort, v); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func createTableDb(ctx *cli.Context) error {
	d := &db.DB{}
	if err := d.ConnectDb(config.DBPath, config.DBName); err != nil {
		debug.PrintStack()
		return err
	}
	if err := d.CreateDb(); err != nil {
		return err
	}
	log.Print("Tables created")
	return nil
}

func appRoot() error {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		return errors.New("Wow, ^.^ dumb")
	}

	app.Commands = []*cli.Command{
		{Name: "start", Action: startApp},
		{Name: "createDb", Action: createTableDb},
	}

	return app.Run(os.Args)
}
func main() {
	go freeMemory()
	if err := appRoot(); err != nil {
		panic(err)
	}
}

func freeMemory() {
	for {
		fmt.Println("run gc")
		start := time.Now()
		runtime.GC()
		debug.FreeOSMemory()
		elapsed := time.Since(start)
		fmt.Printf("gc took %s\n", elapsed)
		time.Sleep(15 * time.Minute)
	}
}
