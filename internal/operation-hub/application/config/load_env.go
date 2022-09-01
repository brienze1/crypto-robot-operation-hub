package config

import (
	logg "github.com/brienze1/crypto-robot-operation-hub/pkg/log"
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

const alternateParentFolderName = "crypto-robot-operation-hub"
const configDirPath = "/config/"

func LoadEnv() {
	env := os.Getenv("OPERATION_HUB_ENV")

	if "" == env {
		env = "development"
	}
	load(".env." + env)
	load(".env") // The Original .env
}

func load(file string) {
	err := godotenv.Load("." + configDirPath + file)
	if err != nil {
		rootPath := getRootPath(alternateParentFolderName)
		err := godotenv.Load(rootPath + configDirPath + file)
		if err != nil {
			logg.Logger().Error(err, "failed loading env file "+file)
			panic("Error loading file: " + file)
		}
	}
}

func getRootPath(dirName string) string {
	projectName := regexp.MustCompile(`^(.*` + dirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	return string(projectName.Find([]byte(currentWorkDirectory)))
}

type env string

const (
	test env = "test"
)

func LoadTestEnv() {
	_ = os.Setenv("OPERATION_HUB_ENV", string(test))

	LoadEnv()
}
