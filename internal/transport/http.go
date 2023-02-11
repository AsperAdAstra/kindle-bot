package transport

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, dest string) error {
	out, err := os.Create(dest)
	defer out.Close()
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}
