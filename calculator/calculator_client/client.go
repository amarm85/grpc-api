package main

import (
	"io"

	"context"
	"github.com/amarm85/grpc-api/calculator/calculatorpb"
	"log"
	"google.golang.org/grpc"
)


func main(){
	cc, err :=  grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil { 
		log.Fatalf("Error in dailing grpc connection: %v", err)
	}

	c := calculatorpb.NewCalculatorClient(cc)

	//doUnary(c)
	doServerStreaming(c)
}

func doUnary(c calculatorpb.CalculatorClient){

	req :=  &calculatorpb.SumRequest{
		FirstNum: 14553,
		SecondNum : 10233332,
		
	}

	res, err := c.Sum(context.Background(),req)
	
	if err != nil {
		log.Fatalf("Error is calling grpc Sum: %v",err)
	}
	
	log.Printf("Sum response %v", res.GetSum())
	
}

func doServerStreaming(c calculatorpb.CalculatorClient){
	num := 345666
	log.Printf("doing PrimeNumberDecomposition calls for %v", num)

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Num: int32(num),
	}

	res , err := c.PrimeNumberDecomposition(context.Background(),req)

	if err != nil {
		 log.Fatalf("Error in doServerStreaming %v", err)
		 return
	}

	for {
		msg , err := res.Recv()

		if err == io.EOF { 
			break
		}

		if err != nil {
			log.Fatalf("Error in doServerStreaming %v", msg.GetResult())
		}

		log.Printf("Prime Number Decomposition: %v",msg.GetResult() )


	}

	
}