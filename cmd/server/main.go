package main

import (
	"context"
	"log"
	"net"

	"github.com/SiddhantTotade/go_html_to_pdf/internal/pdf"
	pb "github.com/SiddhantTotade/go_html_to_pdf/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedPDFGeneratorServer
}

func (s *server) GeneratePDF(ctx context.Context, req *pb.PDFRequest) (*pb.PDFResponse, error) {
	// Validate request
	if req.HtmlContent == "" {
		return &pb.PDFResponse{
			Success: false,
			Message: "HTML content is required",
		}, status.Error(codes.InvalidArgument, "HTML content cannot be empty")
	}

	log.Printf("Generating PDF for file: %s", req.FileName)

	// Generate PDF
	pdfData, err := pdf.GeneratePDF(req.HtmlContent)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return &pb.PDFResponse{
			Success: false,
			Message: "Failed to generate PDF: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	log.Printf("PDF generated successfully, size: %d bytes", len(pdfData))

	return &pb.PDFResponse{
		PdfData: pdfData,
		Message: "PDF generated successfully",
		Success: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024), // 10MB max receive
		grpc.MaxSendMsgSize(10*1024*1024), // 10MB max send
	)

	pb.RegisterPDFGeneratorServer(grpcServer, &server{})

	log.Println("gRPC PDF Generator running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
