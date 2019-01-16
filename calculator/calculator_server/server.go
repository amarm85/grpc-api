package main

import (
	"math"
	"google.golang.org/grpc/codes"
	"context"
	"fmt"
	"github.com/amarm85/grpc-api/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("grpc function Sum has been invoked with %v", req)

	firstNum := req.GetFirstNum()
	secondNum := req.GetSecondNum()

	res := &calculatorpb.SumResponse{
		Sum: firstNum + secondNum,
	}

	return res, nil

}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.Calculator_PrimeNumberDecompositionServer) error {
	num := req.GetNum()
	k := int32(2)
	if num < 1 {
		return fmt.Errorf("Number %v is not valid", num)

	}
	for num > 1 {

		if num%k == 0 {
			//send the stream
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: k,
			}
			stream.Send(res)
			num = num / k
		} else {
			k++
		}

	}
	return nil
}

func (s *server) ComputeAverage(stream calculatorpb.Calculator_ComputeAverageServer) error {

	log.Printf("grpc function ComputeAverage has been invoked ")

	var sum, count int64

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// finished getting
			// log.Printf("grpc function ComputeAverage count %v ", count)
			// log.Printf("grpc function ComputeAverage sum %v ", sum)
			average := float64(sum) / float64(count)
			// log.Printf("grpc function ComputeAverage average %v ", sum)

			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Error in ComputeAverage %v", err)
		}

		sum += msg.GetInNum()
		count++

	}
	return nil

}

func (s *server) FindMaximum(stream calculatorpb.Calculator_FindMaximumServer) error {

	log.Printf("grpc function FindMaximum has been invoked ")

	var maxNum int32
	for {

		req, err := stream.Recv()

		if err == io.EOF {
			return nil

		}

		if err != nil {
			log.Fatalf("Error in rpc FindMaximum %v", err)
			return err
		}

		num := req.GetNumber()

		if num > maxNum {
			maxNum = num

			err := stream.Send(&calculatorpb.FindMaximumRes{
				MaxNum: maxNum,
			})

			if err != nil {
				log.Fatalf("Error in rpc FindMaximum %v", err)
				return err
			}
		}


	}

}

func (s *server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	
	log.Printf("grpc function SquareRoot has been invoked ")

	number := req.GetNumber()

	if number < 0 {
		return nil, 
		status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Recieved negative number: %v",number),
		)
	}

	return &calculatorpb.SquareRootResponse{
		Sqroot : math.Sqrt(float64 (number)),
	},nil

}


func main() {

	//get the listner
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Error is getting listner: %v", err)

	}
	log.Println("gRpc server is running at localhost:50051")
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Error is creating grpc server: %v", err)
	}

}
