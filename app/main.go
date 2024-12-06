package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Simple struct {
	Id   int
	Name string
}

func main() {

	time.Sleep(60 * 1000)

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("MONGODB_URI is not set")
		return
	}

LABEL:
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		fmt.Println(err)
		time.Sleep(1000 * 10)
		fmt.Println("trying to connect again")
		goto LABEL
	}

	coll := client.Database("test").Collection("test")

	_, err = coll.InsertOne(context.TODO(), bson.M{"id": 1, "name": "test"})
	if err != nil {
		fmt.Println(err)
		goto LABEL
	}

	f := func(wg *sync.WaitGroup) {

		defer wg.Done()

		for {
			r := coll.FindOne(context.TODO(), bson.M{"id": 1, "name": "test"})
			v := Simple{}

			if r.Err() != nil {
				fmt.Println(r.Err().Error())
				time.Sleep(1000 * 10)
				continue
			}

			r.Decode(
				&v,
			)
			fmt.Println("Success ", v)

			time.Sleep(1000 * 10)
		}

	}

	wg := new(sync.WaitGroup)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go f(wg)
	}

	wg.Wait()

}
