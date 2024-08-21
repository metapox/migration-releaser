package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/metapox/migration-releaser/handlers"
	"github.com/pkg/errors"
	"log/slog"
	"os"
)

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
	debugMode := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	var logger *slog.Logger
	if *debugMode {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)

	err := Migrate(ctx, event)
	if err != nil {
		logger.Error(
			err.Error(),
			"StackTrace", err.(interface{ StackTrace() errors.StackTrace }).StackTrace(),
		)
		return err
	}
	return nil
}

func Migrate(ctx context.Context, event MyEvent) error {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		return errors.New("DB_TYPE is not set")
	}
	handler, err := handlers.NewDatabaseHandler(dbType)
	if err != nil {
		return err
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
		initDb := os.Getenv("INIT_DB_NAME")
		if initDb == "" {
			return errors.WithStack(errors.New("INIT_DB_NAME is not set"))
		}

		dbUrl, err := dataBaseUrl(initDb)
		if err != nil {
			return err
		}

		slog.Info("Creating database " + item.Key)
		err = handler.CreateDatabase(dbUrl, item.Key)
		if err != nil {
			return err
		}
		slog.Info("Database " + item.Key + " is created")

		targetDB := item.Key
		dbUrl, err = dataBaseUrl(targetDB)
		if err != nil {
			return err
		}

		slog.Info("Migrating database " + targetDB)
		err = handler.UpMigrate(dbUrl, "")
		if err != nil {
			return err
		}

		slog.Info("Database " + targetDB + " is migrated")
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
		os.Exit(1)
	}

	slog.Info("Succeeded to handle request")
}
