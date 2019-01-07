package main

import (
	"fmt"
	"log"
	"google.golang.org/grpc"
	"net"
	"github.com/amarm85/grpc-api/calculator/calculatorpb"
	"context"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error){
	log.Printf("grpc function Sum has been invoked with %v",req)

	firstNum := req.GetFirstNum()
	secondNum := req.GetSecondNum()

	res := &calculatorpb.SumResponse{
		Sum: firstNum + secondNum,
	}

	return res , nil

}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.Calculator_PrimeNumberDecompositionServer) error{
	num := req.GetNum()
	k := int32(2)
	if num < 1 { 
		return fmt.Errorf("Number %v is not valid", num)

	}
	for num > 1 {

		if num % k == 0 {
			//send the stream
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: k,
			}
			stream.Send(res)
			num = num/k
		}else {
			k++
		}
		
	}
	return nil
}

func main(){
	
	//get the listner 
	lis , err :=  net.Listen("tcp","0.0.0.0:50051")
	
	if err != nil { 
		log.Fatalf("Error is getting listner: %v", err)

	}
	log.Println("gRpc server is running at localhost:50051")
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServer(s,&server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("Error is creating grpc server: %v",err)
	}

}