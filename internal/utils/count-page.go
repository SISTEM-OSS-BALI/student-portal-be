package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFPageCountOptions struct {
	Timeout   time.Duration
	MaxBytes  int64
	UserAgent string
}

func CountPDFPagesFromURL(url string, opt *PDFPageCountOptions) (int, error) {
	if url == "" {
		return 0, errors.New("url is empty")
	}

	// defaults
	o := PDFPageCountOptions{
		Timeout:   20 * time.Second,
		MaxBytes:  25 * 1024 * 1024, // 25MB
		UserAgent: "Mozilla/5.0",
	}
	if opt != nil {
		if opt.Timeout > 0 {
			o.Timeout = opt.Timeout
		}
		if opt.MaxBytes > 0 {
			o.MaxBytes = opt.MaxBytes
		}
		if opt.UserAgent != "" {
			o.UserAgent = opt.UserAgent
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), o.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", o.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return 0, fmt.Errorf("fetch failed: %s", resp.Status)
	}

	// limit size
	lr := io.LimitReader(resp.Body, o.MaxBytes+1)
	data, err := io.ReadAll(lr)
	if err != nil {
		return 0, err
	}
	if int64(len(data)) > o.MaxBytes {
		return 0, fmt.Errorf("pdf too large: %d bytes (limit %d)", len(data), o.MaxBytes)
	}

	// IMPORTANT: PageCount butuh io.ReadSeeker
	rs := bytes.NewReader(data)

	conf := model.NewDefaultConfiguration()
	return api.PageCount(rs, conf)
}