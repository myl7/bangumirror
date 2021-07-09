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

		obj, err := reqBase(id)
		if err != nil {
			log.Println(errMsg)
			continue
		}

		eps, err := reqEp(id)
		if err != nil {
			log.Println(errMsg)
			continue
		}

		obj["eps"] = eps

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

func reqBase(id int) (map[string]interface{}, error) {
	res, err := http.Get(GetUrl(id))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return nil, err
	}

	errInfo, ok := obj["error"]
	if ok {
		log.Println("Req error: " + errInfo.(string))
		return nil, err
	}

	return obj, nil
}

func reqEp(id int) (interface{}, error) {
	res, err := http.Get(GetEpUrl(id))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return nil, err
	}

	errInfo, ok := obj["error"]
	if ok {
		log.Println("Req error: " + errInfo.(string))
		return nil, err
	}

	return obj["eps"], nil
}
