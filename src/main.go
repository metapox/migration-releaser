package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/metapox/migration-releaser/handlers"
	"log/slog"
	"os"
)

type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MyEvent struct {
	S3Bucket string `json:"s3Bucket"`
}

type S3Object struct {
	Key string
}

// S3レスポンスの構造体
type S3Response struct {
	Contents []S3Object
}

func HandleRequest(ctx context.Context, event MyEvent) error {
	dbType := os.Getenv("DB_TYPE")
	handler := handlers.NewDatabaseHandler(dbType)
	if handler == nil {
		return errors.New(dbType + " is not supported")
	}

	//s3Client := s3.New(session.Must(session.NewSession()))
	//resp, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(event.S3Bucket)})
	//if err != nil {
	//	log.Fatalf("Failed to list objects: %v", err)
	//}
	resp := S3Response{
		Contents: []S3Object{
			{Key: "testdb1"},
			{Key: "testdb2"},
			{Key: "testdb3"},
		},
	}

	for _, item := range resp.Contents {
		dbName := item.Key
		dbUrl, err := dataBaseUrl("")
		if err != nil {
			return err
		}
		err = handler.CreateDatabase(dbUrl, item.Key)
		if err != nil {
			return err
		}
		dbUrl, err = dataBaseUrl(dbName)
		err = handler.UpMigrate(dbUrl, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func dataBaseUrl(dbName string) (string, error) {
	dbType, err := os.LookupEnv("DB_TYPE")
	if !err {
		return "", errors.New("DB_TYPE is not set")
	}
	user, err := os.LookupEnv("USER")
	if !err {
		return "", errors.New("USER is not set")
	}
	password, err := os.LookupEnv("PASSWORD")
	if !err {
		return "", errors.New("PASSWORD is not set")
	}
	host, err := os.LookupEnv("HOST")
	if !err {
		return "", errors.New("HOST is not set")
	}
	port, err := os.LookupEnv("PORT")
	if !err {
		return "", errors.New("PORT is not set")
	}
	if dbName == "" {
		dbName = os.Getenv("INIT_DB_NAME")
		if dbName == "" {
			return "", errors.New("INIT_DB_NAME is not set")
		}
	}

	parameterString := os.Getenv("DB_PARAMETER")
	var parameters map[string]interface{}
	if parameterString != "" {
		err := json.Unmarshal([]byte(parameterString), &parameters)
		if err != nil {
			return "", errors.New("DB_PARAMETER is invalid")
		}
	}
	parameter := ""
	for key, value := range parameters {
		if parameter != "" {
			parameter += "&"
		}
		parameter += fmt.Sprintf("%s=%s", key, value)
	}
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s", dbType, user, password, host, port, dbName, parameter), nil
}

//func main() {
//	lambda.Start(HandleRequest)
//}

func main() {
	event := MyEvent{
		S3Bucket: "test",
	}
	err := HandleRequest(context.Background(), event)
	if err != nil {
		slog.Error("Failed to handle request: %v", err)
		os.Exit(1)
	}

	slog.Info("Succeeded to handle request")
}
