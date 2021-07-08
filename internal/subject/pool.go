package subject

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/myl7/bangumirror/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func Start(size int, coll *mongo.Collection) {
	var wait sync.WaitGroup
	ids := make(chan int, size)
	for i := 0; i < size; i++ {
		wait.Add(1)
		go worker(&wait, ids, coll)
	}

	for i := config.SubjectStart; i < config.SubjectEnd; i++ {
		ids <- i
	}
	close(ids)

	wait.Wait()
	fmt.Printf("Fetched subjects %d - %d\n", config.SubjectStart, config.SubjectEnd)
}

func worker(wait *sync.WaitGroup, ids <-chan int, coll *mongo.Collection) {
	defer wait.Done()

	for id := range ids {
		errMsg := fmt.Sprintf("Failed to fetch subject %d", id)

		res, err := http.Get(GetUrl(id))
		if err != nil {
			log.Println(errMsg)
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(errMsg)
			continue
		}

		err = res.Body.Close()
		if err != nil {
			log.Println(errMsg)
			continue
		}

		var obj map[string]interface{}
		err = json.Unmarshal(body, &obj)
		if err != nil {
			log.Println(errMsg)
			continue
		}

		errInfo, ok := obj["error"]
		if ok {
			log.Println(errMsg + ": " + errInfo.(string))
			continue
		}

		_, err = coll.FindOneAndReplace(
			context.Background(),
			bson.D{{"id", obj["id"]}},
			obj,
			options.FindOneAndReplace().SetUpsert(true),
		).DecodeBytes()
		if err != nil && err != mongo.ErrNoDocuments {
			log.Println(errMsg)
			continue
		}
	}
}
