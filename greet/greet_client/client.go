package main

import (
	"context"
	"github.com/amarm85/grpc-api/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}
	c := greetpb.NewGreetServiceClient(cc)
	defer cc.Close()
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doUnaryWithDeadline(c, 1*time.Second)
	doUnaryWithDeadline(c, 5*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	log.Println("Starting unary opertaton")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Amar",
			LastName:  "Singh",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Error whiling calling greet RPC: %v", err)
	}

	log.Printf("Response from Greet GRPC call %v", res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting doServerStreaming opertaton")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Amar",
			LastName:  "Singh",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error whiling calling GreetManyTimes RPC: %v", err)
	}

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			//we have reached the end of the stream
			break
		}

		if err != nil {
			log.Fatalf("Error whiling reading messages GreetManyTimes RPC: %v", err)
		}

		log.Printf("GreetManyTimes response %v", msg.GetResult())

	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting doClientStreaming opertaton")

	req := &greetpb.LongGreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "amar",
			LastName:  "Singh",
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error in doClientStreaming: %v", err)
	}

	for ii := 0; ii < 10; ii++ {
		if err = stream.Send(req); err != nil {
			log.Fatalf("Error in doClientStreaming 2: %v", err)
		}
	}

	msg, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error in doClientStreaming 3: %v", err)
	}

	log.Printf("Server response %v", msg.GetResult())

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {

	log.Println("Starting doBiDiStreaming opertaton")
	reqs := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "amar",
				LastName:  "singh",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "samar",
				LastName:  "singh",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ramar",
				LastName:  "singh",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tamar",
				LastName:  "singh",
			},
		},
	}
	//create a stream
	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error in doBiDiStreaming 1: %v", err)
		return
	}

	waitc := make(chan struct{})
	// send bunch of messages in own go rutine
	go func() {

		for _, req := range reqs {
			log.Printf("sending %v", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()

	}()

	// recv bunch of message
	go func() {

		for {

			res, err := stream.Recv()

			if err == io.EOF {

				break
			}

			if err != nil {
				log.Fatalf("Error in doBiDiStreaming 2: %v", err)
				return
			}

			msg := res.GetResult()
			log.Printf("recv %v", msg)

		}

		close(waitc)

	}()

	//block until everything is done
	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	log.Printf("Starting doUnaryWithDeadline opertaton with time %v",timeout )

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Amar",
			LastName:  "Singh",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)

	if err != nil {
		status, ok := status.FromError(err)
		if ok {
			log.Printf("GRPC error for time %v : %v , message %v", timeout, status.Code(), status.Message())
		} else {
			log.Fatalf("Error in doUnaryWithDeadline: %v", err)

		}
		return
	}

	log.Printf("Response %v", res.GetResult())

}
