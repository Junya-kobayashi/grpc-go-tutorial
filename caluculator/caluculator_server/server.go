package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/Junya-kobayashi/grpc-go-course/caluculator/caluculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
}

func (*server) Sum(ctx context.Context, req *caluculatorpb.SumRequest) (*caluculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v", req)
	firstNumber := req.GetSum().GetFirstNumber()
	secondNumber := req.GetSum().GetSecondNumber()

	result := firstNumber + secondNumber
	res := &caluculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumber(req *caluculatorpb.PrimeNumberRequest, stream caluculatorpb.CaluculatorService_PrimeNumberServer) error {
	fmt.Printf("PrimeNumber function was invoked with %v", req)
	number := req.GetPrimeNumber()

	divisor := int32(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(
				&caluculatorpb.PrimeNumberResponse{
					Result: divisor,
				})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v\n", divisor)
		}
	}

	return nil
}

func (*server) Average(stream caluculatorpb.CaluculatorService_AverageServer) error {
	fmt.Printf("average function was invoked with request")
	sum := float32(0)
	count := float32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := sum / count
			return stream.SendAndClose(&caluculatorpb.AverageResponse{
				Average: average,
			})
		}
		fmt.Println(sum)
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		count++
		sum += req.GetNumber()
	}
}

func (*server) Max(stream caluculatorpb.CaluculatorService_MaxServer) error {
	fmt.Printf("maximum function was invoked with request")

	maxNumber := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
			return err
		}

		var sendErr error
		if maxNumber >= int32(req.GetNumber()) {
			sendErr = stream.Send(&caluculatorpb.MaximumResponse{
				Maximum: maxNumber,
			})
		} else {
			maxNumber = req.GetNumber()
			sendErr = stream.Send(&caluculatorpb.MaximumResponse{
				Maximum: req.GetNumber(),
			})
		}

		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *caluculatorpb.SquareRootRequest) (*caluculatorpb.SquareRootResponse, error) {
	fmt.Printf("square root function was invoked with request\n")
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v\n", number),
		)
	}

	return &caluculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	caluculatorpb.RegisterCaluculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
