package main

import (
	"github.com/amarm85/grpc-api/blog/blogpb"
	"github.com/amarm85/grpc-api/blog/server/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {

	//set the log config to get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	listner, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Error in getting TCP listner %v", err)
	}
	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &service.BlogService{})

	go func() {
		log.Printf("Going to start Blog Service")
		if err = s.Serve(listner); err != nil {
			log.Fatalf("Error in starting gRPC server %v", err)
		}

	}()

	ch := make(chan os.Signal, 1)
	// wait for CNTL + c to exit
	signal.Notify(ch, os.Interrupt)

	//block until a signal is recv
	<-ch
	log.Printf("Going to stop the server")
	s.Stop()
	listner.Close()
	//mClient.Disconnect(ctx)
	log.Printf("gRPC server is stopped")
}
