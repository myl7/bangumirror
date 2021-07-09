package subject

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
)

func fetch(id int, coll *mongo.Collection) error {
	obj, err := reqBase(id)
	if err != nil {
		return err
	}

	eps, err := reqEp(id)
	if err != nil {
		return err
	}

	obj["eps"] = eps

	_, err = coll.FindOneAndReplace(
		context.Background(),
		bson.D{{"id", obj["id"]}},
		obj,
		options.FindOneAndReplace().SetUpsert(true),
	).DecodeBytes()
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	return nil
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
