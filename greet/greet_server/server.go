package main

import (
	"time"
	"strconv"
	"context"
	"github.com/amarm85/grpc-api/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error){
	log.Println("Greet function has been invoked")
	firstName :=  req.GetGreeting().GetFirstName()
	result := "hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	log.Println("GreetManyTimes function has been invoked")
	
	firstName := req.GetGreeting().GetFirstName()

	for ii :=0; ii < 10 ; ii++ {
		result :=  "Hello " +  firstName + " number " + strconv.Itoa(ii)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(res)
		time.Sleep(1000 *  time.Millisecond)
	}
	return nil

}

func main() {

	log.Println("gRpc server is running at localhost:50051")
	//create a listner
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to start TCP listner: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s,&server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to start the server : %v", err)
		return
	}

	
}
