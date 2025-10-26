package service

import (
	"context"
	"fmt"
	"log"

	"github.com/SiddhantTotade/go_html_to_pdf/internal/pdf"
	pb "github.com/SiddhantTotade/go_html_to_pdf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PDFGeneratorServiceServer struct {
	pb.UnimplementedPDFGeneratorServer
}

func (s *PDFGeneratorServiceServer) GeneratePDF(ctx context.Context, req *pb.PDFRequest) (*pb.PDFResponse, error) {
	if req == nil || req.HtmlContent == "" {
		return &pb.PDFResponse{
			Success: false,
			Message: "HTML content is required",
		}, status.Error(codes.InvalidArgument, "HTML content cannot be empty")
	}

	orientation := req.Orientation
	if orientation != "landscape" && orientation != "portrait" {
		orientation = "portrait"
	}

	pdfBytes, err := pdf.GeneratePDF(req.HtmlContent, orientation)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return &pb.PDFResponse{
			Success: false,
			Message: "Failed to generate PDF: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.PDFResponse{
		PdfData: pdfBytes,
		Success: true,
		Message: fmt.Sprintf("PDF generated successfully in %s mode", orientation),
	}, nil
}
