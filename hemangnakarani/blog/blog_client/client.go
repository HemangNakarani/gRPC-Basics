package main

import (
	"context"
	"fmt"
	"gRPC-Course/hemangnakarani/blog/blogpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello ! I'm Client.")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}

	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)
	//createBlogFunc(c)
	//readBlogFunc(c)
	//updateBlogFunc(c)
	//deleteBlogFunc(c)

	listBlogFunc(c)

}

func createBlogFunc(c blogpb.BlogServiceClient) {

	fmt.Println("ADD BLOG: Unary RPC...")

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "hemmmang-12",
			Title:    "My First Blog-12",
			Content:  "Learning gRPC with GoLang-12",
		},
	}

	fmt.Println("Creating BLOG:...")

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not Catch Response : %v", err)
	}

	fmt.Printf("Response Unary: %v\n", res.GetBlog())

}

func readBlogFunc(c blogpb.BlogServiceClient) {

	fmt.Println("Read BLOG: Unary RPC...")

	req := &blogpb.ReadBlogRequest{
		BlogId: "5fc36d4e9406e6d5f6676a69",
	}

	fmt.Println("Reading BLOG:...")

	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not Catch Response : %v\n", err)
		return
	}

	fmt.Printf("Response Unary: %v\n", res.GetBlog())

}

func updateBlogFunc(c blogpb.BlogServiceClient) {

	fmt.Println("Update BLOG: Unary RPC...")

	//{"_id":{"$oid":"5fc36d4e9406e6d5f6676a69"},"author_id":"hemmmang","title":"My First Blog","content":"Learning gRPC with GoLang"}

	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       "5fc36d4e9406e6d5f6676a69",
			AuthorId: "HemangNakarani",
			Title:    "Kya Leka aaya jagat me",
			Content:  "Jo bhi laya me , terko kya panchat hai bhai  !!",
		},
	}

	fmt.Println("Updateing BLOG:...")

	res, err := c.UpdateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not Catch Response : %v\n", err)
		return
	}

	fmt.Printf("Response Unary: %v\n", res.GetBlog())

}

func deleteBlogFunc(c blogpb.BlogServiceClient) {

	fmt.Println("Delete BLOG: Unary RPC...")

	req := &blogpb.DeleteBlogRequest{
		BlogId: "5fc36deb8c71972cae2d9b09",
	}

	fmt.Println("Deleting BLOG:...")

	res, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not Catch Response : %v\n", err)
		return
	}

	fmt.Printf("Response Unary: %v\n", res.GetBlogId())

}

func listBlogFunc(c blogpb.BlogServiceClient) {

	fmt.Println("List BLOG: Unary RPC...")

	fmt.Println("Listing BLOG:...")

	resStream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Could not Catch the Stream : %v", err)
	}

	for {

		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error While Reading the stream: %v", err)
		}

		log.Printf("Response Streaming: %v\n", msg.GetBlog())
	}

}
