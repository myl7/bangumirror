package subject

import (
	"fmt"
	"github.com/myl7/bangumirror/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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

	metNotFound := false
	for id := range ids {
		err := fetch(id, coll)
		if err != nil {
			if _, ok := err.(notFoundError); ok {
				if metNotFound {
					break
				} else {
					metNotFound = true
				}
			}
			metNotFound = false

			log.Println(fmt.Sprintf("Failed to fetch subject %d", id))
			continue
		}
	}
}
