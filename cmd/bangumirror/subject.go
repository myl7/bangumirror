package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/myl7/bangumirror/internal/config"
	"github.com/myl7/bangumirror/internal/subject"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.Background())
	coll := client.Database(config.MongoDb).Collection(config.MongoSubjectColl)

	subject.Start(config.JobNum, coll)
}
