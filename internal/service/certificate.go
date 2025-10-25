package service

import (
	"context"
	"log"
	"strings"

	"github.com/SiddhantTotade/go_html_to_pdf/internal/pdf"
	pb "github.com/SiddhantTotade/go_html_to_pdf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PDFGeneratorServiceServer struct {
	pb.UnimplementedPDFGeneratorServer
}

func (s *PDFGeneratorServiceServer) GeneratePDF(ctx context.Context, req *pb.PDFRequest) (*pb.PDFResponse, error) {
	if req == nil {
		return &pb.PDFResponse{
			Success: false,
			Message: "Request is nil",
		}, status.Error(codes.InvalidArgument, "Request cannot be nil")
	}

	html := req.HtmlContent
	fileName := req.FileName

	if html == "" {
		return &pb.PDFResponse{
			Success: false,
			Message: "HTML content is required",
		}, status.Error(codes.InvalidArgument, "HTML content cannot be empty")
	}

	if fileName == "" {
		fileName = "document.pdf"
	}

	log.Printf("Generating PDF for file: %s", fileName)

	// Replace event details placeholders
	for k, v := range req.EventDetails {
		placeholder := "{{ " + k + " }}"
		html = strings.ReplaceAll(html, placeholder, v)
	}

	// Replace participant details placeholders
	if len(req.Participants) > 0 {
		firstParticipant := req.Participants[0]
		if firstParticipant != nil {
			for k, v := range firstParticipant.ParticipantDetails {
				placeholder := "{{ " + k + " }}"
				html = strings.ReplaceAll(html, placeholder, v)
			}
		}
	}

	// Generate PDF
	pdfBytes, err := pdf.GeneratePDF(html)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		return &pb.PDFResponse{
			Success: false,
			Message: "Failed to generate PDF: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	log.Printf("PDF generated successfully, size: %d bytes", len(pdfBytes))

	return &pb.PDFResponse{
		PdfData: pdfBytes,
		Message: "PDF generated successfully",
		Success: true,
	}, nil
}
