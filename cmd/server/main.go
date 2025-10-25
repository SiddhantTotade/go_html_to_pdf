package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/SiddhantTotade/go_html_to_pdf/internal/pdf"
	pb "github.com/SiddhantTotade/go_html_to_pdf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server implements the gRPC PDF generator service
type server struct {
	pb.UnimplementedPDFGeneratorServer
}

// GeneratePDF handles incoming requests to generate a PDF
func (s *server) GeneratePDF(ctx context.Context, req *pb.PDFRequest) (*pb.PDFResponse, error) {
	// Validate request
	if req.HtmlContent == "" {
		return &pb.PDFResponse{
			Success: false,
			Message: "HTML content is required",
		}, status.Error(codes.InvalidArgument, "HTML content cannot be empty")
	}

	log.Printf("Generating PDF for file: %s", req.FileName)

	// Generate the PDF
	pdfData, err := pdf.GeneratePDF(req.HtmlContent)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return &pb.PDFResponse{
			Success: false,
			Message: "Failed to generate PDF: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	// Save PDF to disk
	outputFile := "output.pdf"
	err = os.WriteFile(outputFile, pdfData, 0644)
	if err != nil {
		log.Printf("Failed to write PDF to file: %v", err)
	} else {
		log.Printf("✅ PDF saved successfully to %s", outputFile)
	}

	log.Printf("PDF generated successfully, size: %d bytes", len(pdfData))

	return &pb.PDFResponse{
		PdfData: pdfData,
		Message: "PDF generated successfully",
		Success: true,
	}, nil
}

func main() {
	// Start listening on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server with message size limits
	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(10*1024*1024), // 10MB max receive
		grpc.MaxSendMsgSize(10*1024*1024), // 10MB max send
	)

	// Register the PDF Generator service
	pb.RegisterPDFGeneratorServer(grpcServer, &server{})

	log.Println("✅ Server running on :50051")

	// Start serving requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
