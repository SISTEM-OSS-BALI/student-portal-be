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

type PDFMergeOptions struct {
	Timeout       time.Duration
	MaxBytes      int64
	UserAgent     string
	AddDividerPage bool
}

func MergePDFsFromURLs(urls []string, opt *PDFMergeOptions) ([]byte, error) {
	if len(urls) < 2 {
		return nil, errors.New("need at least 2 urls to merge")
	}

	// defaults
	o := PDFMergeOptions{
		Timeout:   30 * time.Second,
		MaxBytes:  25 * 1024 * 1024, // 25MB per file
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
		o.AddDividerPage = opt.AddDividerPage
	}

	ctx, cancel := context.WithTimeout(context.Background(), o.Timeout)
	defer cancel()

	readSeekers := make([]io.ReadSeeker, 0, len(urls))
	for _, u := range urls {
		if u == "" {
			return nil, errors.New("url is empty")
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", o.UserAgent)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("fetch failed: %s", resp.Status)
		}

		lr := io.LimitReader(resp.Body, o.MaxBytes+1)
		data, err := io.ReadAll(lr)
		_ = resp.Body.Close()
		if err != nil {
			return nil, err
		}
		if int64(len(data)) > o.MaxBytes {
			return nil, fmt.Errorf("pdf too large: %d bytes (limit %d)", len(data), o.MaxBytes)
		}

		readSeekers = append(readSeekers, bytes.NewReader(data))
	}

	var buf bytes.Buffer
	conf := model.NewDefaultConfiguration()
	if err := api.MergeRaw(readSeekers, &buf, o.AddDividerPage, conf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

