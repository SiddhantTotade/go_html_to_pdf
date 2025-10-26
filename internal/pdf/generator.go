package pdf

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func GeneratePDF(htmlContent string, orientation string) ([]byte, error) {
	var pdfBuf []byte

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	dataURI := "data:text/html;charset=utf-8," + url.PathEscape(htmlContent)

	isLandscape := false
	if orientation == "landscape" {
		isLandscape = true
	}

	err := chromedp.Run(ctx,
		chromedp.Navigate(dataURI),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithLandscape(isLandscape).
				WithPrintBackground(true).
				WithPaperWidth(8.27).
				WithPaperHeight(11.7).
				Do(ctx)
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

	return pdfBuf, nil
}
