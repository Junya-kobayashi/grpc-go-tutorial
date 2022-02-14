package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Junya-kobayashi/grpc-go-course/caluculator/caluculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer conn.Close()

	c := caluculatorpb.NewCaluculatorServiceClient(conn)

	//doSum(c)
	//doFactor(c)
	// doAverage(c)
	// doMaximum(c)
	doErrorUnary(c)
}

func doSum(c caluculatorpb.CaluculatorServiceClient) {
	fmt.Println("Starting to do a sum RPC...")
	req := &caluculatorpb.SumRequest{
		Sum: &caluculatorpb.Sum{
			FirstNumber:  10,
			SecondNumber: 3,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}

func doFactor(c caluculatorpb.CaluculatorServiceClient) {
	fmt.Println("Starting to do a factor RPC...")

	req := &caluculatorpb.PrimeNumberRequest{
		PrimeNumber: 120,
	}

	resStream, err := c.PrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling primenumber response %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("response from Primenumber is %v", msg.GetResult())
	}
}

func doAverage(c caluculatorpb.CaluculatorServiceClient) {
	fmt.Println("Starting to do a client streaming average RPC...")

	request := []*caluculatorpb.AverageRequest{
		&caluculatorpb.AverageRequest{
			Number: 2.0,
		},
		&caluculatorpb.AverageRequest{
			Number: 6.0,
		},
		&caluculatorpb.AverageRequest{
			Number: 3.0,
		},
		&caluculatorpb.AverageRequest{
			Number: 10.0,
		},
	}

	stream, err := c.Average(context.Background())

	if err != nil {
		log.Fatalf("error while calling avarage: %v", err)
	}

	for _, req := range request {
		fmt.Printf("sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from avearage response %v", err)
	}

	fmt.Printf("average response %v\n", res)

}

func doMaximum(c caluculatorpb.CaluculatorServiceClient) {
	fmt.Println("starting to do a maximum streaming rpc...")
	numberList := []int{5, 3, 6, 8, 9, 1}

	stream, err := c.Max(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream : %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		for _, number := range numberList {
			fmt.Printf("sending message: %v\n", number)
			stream.Send(
				&caluculatorpb.MaximumRequest{
					Number: int32(number),
				},
			)
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
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetMaximum())
		}
		close(waitc)
	}()

	<-waitc
}

func doErrorUnary(c caluculatorpb.CaluculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")
	// correct call
	doErrorCall(c, 10)
	doErrorCall(c, -2)
}

func doErrorCall(c caluculatorpb.CaluculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &caluculatorpb.SquareRootRequest{Number: number})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Printf("Error message from server: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big error calling SquareRoot: %v", err)
			return
		}
	}

	fmt.Printf("Result of square root of %v: %v\n", number, res.GetNumberRoot())
	// error call
}
