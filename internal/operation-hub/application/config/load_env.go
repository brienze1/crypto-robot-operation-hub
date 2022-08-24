package config

import (
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

const projectDirName = "crypto-robot-operation-hub"
const configDirPath = "/config/"

func LoadEnv() {
	env := os.Getenv("OPERATION_HUB_ENV")

	if "" == env {
		env = "local"
	}
	load(".env." + env)
	load(".env") // The Original .env
}

func load(file string) {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + configDirPath + file)
	if err != nil {
		panic("Error loading .env file")
	}
}

type env string

const (
	test env = "test"
)

func LoadTestEnv() {
	_ = os.Setenv("OPERATION_HUB_ENV", string(test))

	LoadEnv()
}
