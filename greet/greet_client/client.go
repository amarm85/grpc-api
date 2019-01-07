package main

import (
	"io"
	"context"
	"log"
	"github.com/amarm85/grpc-api/greet/greetpb"
	"google.golang.org/grpc"
)

func main(){
	cc, err := grpc.Dial("localhost:50051",grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect : %v",err)
	}
	c :=  greetpb.NewGreetServiceClient(cc)
	defer cc.Close()
	//doUnary(c)
	doServerStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient){
	log.Println("Starting unary opertaton")
	req :=  &greetpb.GreetRequest{
		Greeting : &greetpb.Greeting{
			FirstName: "Amar",
			LastName: "Singh",
		},
	}
	res, err := c.Greet(context.Background(),req)

	if err != nil { 
		log.Fatalf("Error whiling calling greet RPC: %v",err)
	}

	log.Printf("Response from Greet GRPC call %v",res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient){
	log.Println("Starting doServerStreaming opertaton")

	req :=  &greetpb.GreetManyTimesRequest{
		Greeting : &greetpb.Greeting{
			FirstName: "Amar",
			LastName: "Singh",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(),req)

	if err != nil { 
		log.Fatalf("Error whiling calling GreetManyTimes RPC: %v",err)
	}

	for {
	
		msg , err := stream.Recv()
	
		if err == io.EOF {
		//we have reached the end of the stream 
			break
		} 
		
		if err != nil {
			log.Fatalf("Error whiling reading messages GreetManyTimes RPC: %v",err)
		}

		log.Printf("GreetManyTimes response %v", msg.GetResult())

	}

}