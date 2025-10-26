package main

import (
	"log"
	"net"

	"github.com/SiddhantTotade/go_html_to_pdf/internal/service"
	pb "github.com/SiddhantTotade/go_html_to_pdf/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024),
		grpc.MaxSendMsgSize(10*1024*1024),
	)

	pb.RegisterPDFGeneratorServer(grpcServer, &service.PDFGeneratorServiceServer{})

	log.Println("✅ PDF gRPC server running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
