package agent

import (
	"bytes"
	"compress/gzip"
	"io"

	"resty.dev/v3"
)

func gzipMiddleware(client *resty.Client, r *resty.Request) error {
	req := r.RawRequest

	if req == nil || req.Body == nil {
		return nil
	}

	if req.Header.Get("Content-Type") != "application/json" {
		return nil
	}

	if req.Header.Get("Content-Encoding") != "" {
		return nil
	}

	body, err := io.ReadAll(req.Body)

	if err != nil {
		return err
	}

	req.Body.Close()

	var buf bytes.Buffer

	gz := gzip.NewWriter(&buf)

	_, err2 := gz.Write(body)

	if err2 != nil {
		gz.Close()
		return err2
	}

	err3 := gz.Close()

	if err3 != nil {
		return err3
	}

	compressed := buf.Bytes()

	req.Body = io.NopCloser(bytes.NewReader(compressed))
	req.ContentLength = int64(len(compressed))
	req.Header.Set("Content-Encoding", "gzip")

	// Важно для retry.
	// req.GetBody = func() (io.ReadCloser, error) {
	// 	return io.NopCloser(bytes.NewReader(compressed)), nil
	// }

	return nil
}
