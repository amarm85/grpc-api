package main

import (
	"context"
	"github.com/amarm85/grpc-api/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Println("Greet function has been invoked")
	firstName := req.GetGreeting().GetFirstName()
	result := "hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {

	log.Println("GreetManyTimes function has been invoked")

	firstName := req.GetGreeting().GetFirstName()

	for ii := 0; ii < 10; ii++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(ii)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil

}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Println("LongGreet function has been invoked")
	result := "hello "

	for {

		msg, err := stream.Recv()

		if err == io.EOF {
			//finished reading stream
			stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
			break
		}

		if err != nil {
			log.Fatalf("Error in LongGreet %v", err)
		}

		firstName := msg.GetGreeting().GetFirstName()

		result += firstName + "!"
	}
	return nil
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Println("GreetEveryone rpc function has been invoked")

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error in GreetEveryone %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()

		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: "Hello " + firstName,
		})
		if err != nil {
			log.Fatalf("Error in GreetEveryone 2 %v", err)
			return err
		}
	}

}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	log.Println("GreetWithDeadline rpc function has been invoked")

	for ii := 0; ii < 3; ii++ {

		if ctx.Err() == context.Canceled {
			//client has cancelled request
			log.Println("GreetWithDeadline rpc function client cancelled request")
			return nil, status.Error(codes.DeadlineExceeded, "client cancelled the request")
		}
		time.Sleep(1 * time.Second)
		log.Println("GreetWithDeadline rpc function slept for 1 second")
	}

	firstName := req.GetGreeting().GetFirstName()

	return &greetpb.GreetWithDeadlineResponse{
		Result: "Hello " + firstName,
	}, nil
}

func main() {

	log.Println("gRpc server is running at localhost:50051")
	//create a listner
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to start TCP listner: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to start the server : %v", err)
		return
	}

}
