package main

import (
	"context"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "pi/pi"
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	port = ":50051"
)

func ReadFile(skip int64) string {
	// prompt for name of file
	// fmt.Println("Enter the file name\n")
	// var fileName string
	// fmt.Scan(&fileName)

	// open file
	// file, err := os.Open(fileName)
	file, err := os.Open("pi.txt")
	if err != nil {
		log.Printf("Error: %s", err)
	}
	reader := bufio.NewReader(file)
	_, _ = reader.Discard(int(skip))

	lines := []string{"."}
	for {
		// read line from file
		line, fileErr := reader.ReadString('\n')
		lines = append(lines, line[:20])
		if fileErr != nil {
			break
		}
	}

	file.Close()
	return strings.Join(lines[:], "")
}

func CountDigits(line string) []int64 {
	numberCount := make([]int64, 10)
	line = strings.Replace(line, "$", "", -1)
	line = strings.Replace(line, ".", "", -1)
	for i := 0; i < len(line); i++ {
		if s, err := strconv.Atoi(string(line[i])); err == nil {
			numberCount[s] += 1
		}
	}
	return numberCount
}

type server struct {
	pb.UnimplementedPiCounterServer
}

func (s *server) Test(ctx context.Context, in *pb.TestReceive) (*pb.TestSend, error) {
	log.Printf("Received: %v", in.GetMessage())
	return &pb.TestSend{Message: "Hello " + in.GetMessage()}, nil
}

func (s *server) CountPiDigits(ctx context.Context, in *pb.CountRequest) (*pb.CountResponse, error) {
	log.Printf("To skip: %v", in.GetSkip())
	line := ReadFile(in.GetSkip())
	count := CountDigits(line)
	return &pb.CountResponse{Count: count}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPiCounterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}