package main

import (
	"context"
	"github.com/amarm85/grpc-api/blog/blogpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {

	log.Printf("Starting the blog client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error in dailing the gRPC server %v", err)
	}

	bsClient := blogpb.NewBlogServiceClient(cc)

	//insertBlog(bsClient)
	listBlog(bsClient)
}

func insertBlog(bsClient blogpb.BlogServiceClient) {

	blog := &blogpb.CreateBlogReq{
		Blog: &blogpb.Blog{
			AuthorId: "Amar Jeet singh",
			Content:  "This is a sample blog",
			Title:    "Sample blog",
		},
	}

	res, err := bsClient.CreateBlog(context.Background(), blog)

	if err != nil {
		log.Fatalf("Error in creating blog %v", err)
	}

	log.Printf("Blog created successfully %v", res.GetBlog())

	//read blog rpc call
	readBlogReq := blogpb.ReadBlogReq{
		ID: res.GetBlog().GetId(),
	}

	readRes, err := bsClient.ReadBlog(context.Background(), &readBlogReq)

	if err != nil {
		log.Fatalf("Error in Reading blog %v", err)
	}
	log.Printf("Blog read successfully %v", readRes.GetBlog())

	//update the blog item

	updateBlog := &blogpb.UpdateBlogReq{
		Blog: &blogpb.Blog{
			Id:       readRes.GetBlog().GetId(),
			AuthorId: "Amar Jeet singh",
			Content:  "This is a sample updated",
			Title:    "Sample blog updated",
		},
	}

	updateRes, err := bsClient.UpdateBlog(context.Background(), updateBlog)

	if err != nil {
		log.Fatalf("Error in updating blog %v", err)
	}

	log.Printf("Blog updated  successfully %v", updateRes.GetBlog())

	//Delete the blog
	deleteReq := &blogpb.DeleteBlogReq{
		ID: readRes.GetBlog().GetId(),
	}

	deleteRes, err := bsClient.DeleteBlog(context.Background(), deleteReq)

	if err != nil {
		log.Fatalf("Error in deleting blog %v", err)
	}

	log.Printf("Blog Deleted  successfully %v", deleteRes.GetResult())
}

func listBlog(bsClient blogpb.BlogServiceClient) {

	stream, err := bsClient.ListBlog(context.Background(), &blogpb.ListBlogReq{})

	if err != nil {
		log.Fatalf("Error in opening ListBlog Stream %v", err)
	}

	for {

		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error in reading ListBlog Stream %v", err)
			break
		}

		log.Printf("Recv data from server:: %v", res.GetBlog())

	}

}
