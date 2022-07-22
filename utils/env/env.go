package env

import (
	"log"
	"os"
)

type Env int

const (
	DEV = iota
	DOCKER_DEV
	TEST
	STAGE
	PROD
)

var envString = [...]string{
	"dev",
	"docker_dev",
	"test",
	"stage",
	"prod",
}

func (e Env) String() string {
	return envString[int(e)]
}

var env Env

func SetEnv(k string) {
	found := false
	for i, v := range envString {
		if k == v {
			found = true
			env = Env(i)
		}
	}

	if !found {
		log.Fatal("env error")
	}
}

func GetEnv() Env {
	return env
}

func IsDev() bool {
	return GetEnv() == DEV
}

func IsTest() bool {
	return GetEnv() == TEST
}

func IsStage() bool {
	return GetEnv() == STAGE
}

func IsProd() bool {
	return GetEnv() == PROD
}

func InitEnv() {
	e := os.Getenv("RUN_ENV")
	if e != "" {
		SetEnv(e)
	} else {
		SetEnv(envString[0])
	}
}
