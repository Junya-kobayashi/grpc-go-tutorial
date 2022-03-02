package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/Junya-kobayashi/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Blog client")

	tls := true
	opts := grpc.WithInsecure()

	if tls {
		certFile := "ssl/ca.crt"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)
	fmt.Printf("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Junya",
		Title:    "first blog",
		Content:  "content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
	blogId := createBlogRes.GetBlog().GetId()

	//read Blog
	fmt.Println("Reading the blog")

	_, readErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "1dsojfs"})
	if readErr != nil {
		fmt.Printf("Error happened while reading: %v\n", readErr)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogId}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v\n", readBlogErr)
	}

	fmt.Printf("Blog was read %v", readBlogRes)

	// updateBlog
	newBlog := &blogpb.Blog{
		Id:       blogId,
		AuthorId: "Changed Author",
		Title:    "My First Blog (editted)",
		Content:  "Content of the first blog, with some awesome additions!",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v\n", updateErr)
	}
	fmt.Printf("Blog was read: %v", updateRes)

	// deleteBlog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogId})

	if deleteErr != nil {
		fmt.Printf("Error happend while deleting: %v\n", deleteErr)
	}

	fmt.Printf("Blog was deleted: %v \n", deleteRes)

	// list blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling primenumber response %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Println(res.GetBlog())
	}
}
