package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Simple struct {
	Id   int `bson:"_id"`
	Name string
}

func main() {

	time.Sleep(time.Second * 10)

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("MONGODB_URI is not set")
		return
	}

LABEL:

	// poolMonitor := &event.PoolMonitor{
	// 	Event: func(evt *event.PoolEvent) {

	// 		log.Println(*evt)

	// 	},
	// }

	// cmdMonitor := &event.CommandMonitor{
	// 	Started: func(_ context.Context, evt *event.CommandStartedEvent) {
	// 		log.Println(evt)
	// 	},
	// }

	// svrMonitor := &event.ServerMonitor{
	// 	ServerHeartbeatStarted: func(e *event.ServerHeartbeatStartedEvent) {
	// 		log.Println(e)
	// 	},
	// 	ServerDescriptionChanged: func(e *event.ServerDescriptionChangedEvent) {
	// 		log.Println(e)
	// 	},
	// }

	//clientOptions := options.Client().SetPoolMonitor(poolMonitor).ApplyURI(uri)

	// clientOptions := options.Client().ApplyURI(uri).SetServerMonitor(&event.ServerMonitor{
	// 	ServerDescriptionChanged: func(sdce *event.ServerDescriptionChangedEvent) {
	// 		log.Println("SERVER DESCRIPTION CHANGED", sdce)
	// 	},
	// 	TopologyDescriptionChanged: func(tdce *event.TopologyDescriptionChangedEvent) {
	// 		log.Println("TOPOLOGY DESCRIPTION CHANGED", tdce)
	// 	},
	// 	ServerHeartbeatSucceeded: func(shse *event.ServerHeartbeatSucceededEvent) {
	// 		log.Print("HEARTBEAT SUCCEEDED", shse)
	// 	},
	// })

	clientOptions := options.Client().ApplyURI(uri).SetPoolMonitor(&event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			log.Println(*evt)
		},
	})

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 10)
		fmt.Println("trying to connect again")
		goto LABEL
	}

	coll := client.Database("test").Collection("test")

	_, err = coll.InsertOne(context.TODO(), bson.M{"_id": 1, "name": "test"})
	if err != nil {
		fmt.Println(err)
		goto LABEL
	}

	f := func(wg *sync.WaitGroup) {

		defer wg.Done()

		for {
			r := coll.FindOne(context.TODO(), bson.M{"_id": 1, "name": "test"})
			v := Simple{}

			if r.Err() != nil {
				fmt.Println(r.Err().Error())
				time.Sleep(time.Second * 10)
				continue
			}

			r.Decode(
				&v,
			)
			
			log.Println("Success ", v)

			time.Sleep(time.Second * 10)
		}

	}

	wg := new(sync.WaitGroup)

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go f(wg)
	}

	wg.Wait()

}
