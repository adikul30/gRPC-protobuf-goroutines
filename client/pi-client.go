package main

import (
	"context"
	"log"
	// "os"
	"time"

	"google.golang.org/grpc"
	pb "pi/pi"
)

const (
	address     = "localhost:50051"
	defaultName = "pi"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPiCounterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	ch := make (chan []int64)

	for i := 0; i < 10; i++ {
		go func (skip int64) {
			r, err := c.CountPiDigits(ctx, &pb.CountRequest{Skip: skip})
			if err != nil {
				log.Fatalf("could not count: %v", err)
			}
			count := r.GetCount()
			// log.Printf("Count: %s", r.GetCount())
			// var temp int64
			// for _, val := range count {
				// log.Printf("Inside go routine\n")
				// log.Printf("%d: %d\n", i, val)
				// temp += val
			// }
			// log.Printf("%d", temp)
			ch <- count
		} (int64(i * 20))
	}
	var count [10]int64
	var total int64
	for i := 0; i < 10; i++ {
		val := <- ch
		for index, v := range val {
			count[index] += v
			total += v
		}
	}

	for index, v := range count {
		log.Printf("%d: %d\n", index, v)
	}

	log.Printf("Total: %d", total)
}