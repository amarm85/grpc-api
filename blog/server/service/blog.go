package service

import (
	"context"
	"fmt"
	"github.com/amarm85/grpc-api/blog/blogpb"
	"github.com/amarm85/grpc-api/blog/server/db"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

var blogCollection *mongo.Collection

func init() {
	blogCollection = db.ConnectToMongoDB().Database("confusion").Collection("blog")
}

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

//Service for Create, update, delete and get Blogs
type BlogService struct {
}

func (s *BlogService) CreateBlog(ctx context.Context, req *blogpb.CreateBlogReq) (*blogpb.CreateBlogRes, error) {

	log.Printf("gRPC service CreateBlog is called ")

	blog := req.GetBlog()

	data := BlogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}
	res, err := blogCollection.InsertOne(ctx, data)

	if err != nil {
		log.Printf("Mongo DB Error in Insert %v", err)
		return nil, status.Errorf(codes.Internal, "Internal err %v", err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Cannot convert to OID %v", err)
	}
	return &blogpb.CreateBlogRes{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (s *BlogService) ReadBlog(ctx context.Context, req *blogpb.ReadBlogReq) (*blogpb.ReadBlogRes, error) {

	log.Printf("gRPC service ReadBlog is called ")

	id := req.GetID()
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID %v is not a valid objectID", id)
	}

	//create an empty struct
	blog := &BlogItem{}

	filter := bson.M{"_id": oid}
	res := blogCollection.FindOne(ctx, filter)

	err = res.Decode(blog)

	if err != nil {
		log.Printf("Error in Decoding result %v", err)
		return nil, status.Errorf(codes.NotFound, "Document not found with ID %v", id)
	}

	return &blogpb.ReadBlogRes{
		Blog: &blogpb.Blog{
			AuthorId: blog.AuthorID,
			Content:  blog.Content,
			Id:       blog.ID.Hex(),
			Title:    blog.Title,
		},
	}, nil

}

func (s *BlogService) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogReq) (*blogpb.UpdateBlogRes, error) {
	log.Printf("gRPC service UpdateBlog is called ")

	blog := req.GetBlog()

	data := BlogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	id := blog.GetId()
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID %v is not a valid objectID", id)
	}
	log.Printf("update document key is %v", oid)
	filter := bson.M{"_id": oid}
	_, err = blogCollection.ReplaceOne(ctx, filter, data)

	if err != nil {
		log.Printf("Error in updating document %v", err)
		return nil, status.Errorf(codes.NotFound, "Document not found with ID %v", id)
	}

	return &blogpb.UpdateBlogRes{
		Blog: &blogpb.Blog{
			AuthorId: data.AuthorID,
			Content:  data.Content,
			Id:       blog.GetId(),
			Title:    data.Title,
		},
	}, nil

}

func (s *BlogService) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogReq) (*blogpb.DeleteBlogRes, error) {

	log.Printf("gRPC service DeleteBlog is called ")

	id := req.GetID()
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID %v is not a valid objectID", id)
	}

	filter := bson.M{"_id": oid}

	res, err := blogCollection.DeleteOne(ctx, filter)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Document not found with ID %v", id)
	}

	return &blogpb.DeleteBlogRes{Result: fmt.Sprintf("%v document has been deleted", res.DeletedCount)}, nil
}

func (s *BlogService) ListBlog(req *blogpb.ListBlogReq, stream blogpb.BlogService_ListBlogServer) error {

	log.Printf("gRPC service ListBlog is called ")

	cursor, err := blogCollection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Printf("Unknown error occur:: %v", err)
		return status.Errorf(codes.Internal, "Unknown error occur %v", err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		data := BlogItem{}
		err := cursor.Decode(&data)

		if err != nil {
			return status.Errorf(codes.Internal, "Error in decoding:: %v", err)
		}

		stream.Send(&blogpb.ListBlogRes{
			Blog: &blogpb.Blog{
				AuthorId: data.AuthorID,
				Content:  data.Content,
				Id:       data.ID.Hex(),
				Title:    data.Title,
			},
		})

	}
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, "Error in closing cursor::  %v", err)
	}

	return nil

}
