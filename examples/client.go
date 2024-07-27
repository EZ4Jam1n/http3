/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Load the TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"h3", "h3-29"},
	}

	// Create a new HTTP3 RoundTripper
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: tlsConfig,
		QUICConfig:      &quic.Config{},
	}

	client := &http.Client{
		Transport: roundTripper,
	}

	// Define the URLs to test
	urls := []string{
		"https://127.0.0.1:8080/demo/tile",
		"https://127.0.0.1:8080/ping",
		"https://127.0.0.1:8080/struct",
		"https://127.0.0.1:8080/v1/hello/world",
	}

	for _, url := range urls {
		fmt.Printf("Requesting %s\n", url)
		req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
		if err != nil {
			log.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to perform request: %v", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}
		defer resp.Body.Close()

		fmt.Printf("Response status: %s\n", resp.Status)
		fmt.Printf("Response body: %s\n", body)
	}

	// Close the HTTP3 RoundTripper
	roundTripper.Close()
}
