package pdf

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PDFOptions struct {
	PrintBackground bool
	PaperWidth      float64
	PaperHeight     float64
	MarginTop       float64
	MarginBottom    float64
	MarginLeft      float64
	MarginRight     float64
}

func DefaultPDFOptions() *PDFOptions {
	return &PDFOptions{
		PrintBackground: true,
		PaperWidth:      8.5,
		PaperHeight:     11.0,
		MarginTop:       0.4,
		MarginBottom:    0.4,
		MarginLeft:      0.4,
		MarginRight:     0.4,
	}
}

func GeneratePDF(htmlContent string) ([]byte, error) {
	return GeneratePDFWithOptions(htmlContent, DefaultPDFOptions())
}

func GeneratePDFWithOptions(htmlContent string, opts *PDFOptions) ([]byte, error) {
	var pdfBuf []byte

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(func(format string, v ...interface{}) {}), // Disable logs
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	dataURI := "data:text/html;charset=utf-8," + htmlContent

	err := chromedp.Run(ctx,
		chromedp.Navigate(dataURI),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond), // Give time for fonts/images to load
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Configure PDF printing
			printParams := page.PrintToPDF().
				WithPrintBackground(opts.PrintBackground).
				WithPaperWidth(opts.PaperWidth).
				WithPaperHeight(opts.PaperHeight).
				WithMarginTop(opts.MarginTop).
				WithMarginBottom(opts.MarginBottom).
				WithMarginLeft(opts.MarginLeft).
				WithMarginRight(opts.MarginRight).
				WithPreferCSSPageSize(false)

			buf, _, err := printParams.Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to print to PDF: %w", err)
			}
			pdfBuf = buf
			return nil
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("chromedp error: %w", err)
	}

	if len(pdfBuf) == 0 {
		return nil, fmt.Errorf("generated PDF is empty")
	}

	return pdfBuf, nil
}
