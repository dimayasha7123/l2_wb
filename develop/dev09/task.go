package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	var output = "index.html"
	flag.StringVar(&output, "o", output, "path to output file")
	flag.Parse()

	reqURL := flag.Arg(0)
	if reqURL == "" {
		fmt.Fprintf(os.Stderr, "url must be non empty, give it with first arg\n")
		os.Exit(2)
	}

	bytes, err := wget(reqURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't make req: %v\n", err)
		os.Exit(1)
	}

	file, err := os.Create(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't create file %s: %v", output, err)
		os.Exit(1)
	}

	_, err = file.Write(bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't write data to file: %v\n", err)
		os.Exit(1)
	}
}

func wget(reqURL string) ([]byte, error) {
	client := http.Client{Timeout: 5 * time.Second}

	_, err := url.Parse(reqURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create http req obj: %v", err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:105.0) Gecko/20100101 Firefox/105.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get data: %v", err)
	}

	bytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("can't read body: %v", err)
	}

	return bytes, nil
}
