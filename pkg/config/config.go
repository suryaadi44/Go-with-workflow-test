package config

import "os"

// var (
// 	DB_USER    = "root"
// 	DB_PASS    = "root"
// 	DB_HOST    = "localhost"
// 	DB_PORT    = "3307"
// 	DB_NAME    = "training_mk2"
// 	PORT       = ":8000"
// 	JWT_SECRET = "secret"
// )

var (
	DB_USER    = os.Getenv("DB_USER")
	DB_PASS    = os.Getenv("DB_PASS")
	DB_HOST    = os.Getenv("DB_HOST")
	DB_PORT    = os.Getenv("DB_PORT")
	DB_NAME    = os.Getenv("DB_NAME")
	PORT       = os.Getenv("PORT")
	JWT_SECRET = os.Getenv("JWT_SECRET")
)
