package main

import (
	"io"
	"time"

	"context"
	"github.com/amarm85/grpc-api/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error in dailing grpc connection: %v", err)
	}

	c := calculatorpb.NewCalculatorClient(cc)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c calculatorpb.CalculatorClient) {

	req := &calculatorpb.SumRequest{
		FirstNum:  14553,
		SecondNum: 10233332,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Error is calling grpc Sum: %v", err)
	}

	log.Printf("Sum response %v", res.GetSum())

}

func doServerStreaming(c calculatorpb.CalculatorClient) {
	num := 345666
	log.Printf("doing PrimeNumberDecomposition calls for %v", num)

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Num: int32(num),
	}

	res, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error in doServerStreaming %v", err)
		return
	}

	for {
		msg, err := res.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error in doServerStreaming %v", msg.GetResult())
		}

		log.Printf("Prime Number Decomposition: %v", msg.GetResult())

	}

}

func doClientStreaming(c calculatorpb.CalculatorClient) {

	log.Printf("doing doClientStreaming calls")
	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error in doClientStreaming %v", err)
	}

	numbers := []int64{3, 5, 9, 54, 23}

	for _, number := range numbers {

		stream.Send(&calculatorpb.ComputeAverageRequest{
			InNum: number,
		})
		log.Printf("Sending  %v", number)
	}

	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error in doClientStreaming 2 %v", err)
	}

	log.Printf("The average is %.2f", msg.GetAverage())

}

func doBiDiStreaming(c calculatorpb.CalculatorClient) {

	log.Printf("doing doBiDiStreaming calls")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error in doBiDiStreaming %v", err)
	}

	var numbers  = []int32{3,4,7,8,2,4,19,3,4}

	waitc := make(chan struct{})

	go func() {
		for _, number := range numbers {
			log.Printf("sending %v", number)
			stream.Send(&calculatorpb.FindMaximumReq{
				Number : number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error in doBiDiStreaming 2 %v", err)
				break
			}

			log.Printf("Max number till now is %v", res.GetMaxNum())
		}

		close(waitc)

	}()

	<-waitc
}
