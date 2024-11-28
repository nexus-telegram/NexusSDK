package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nexus-telegram/NexusSDK/types"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"time"
)

// HTTPClient wraps an HTTP client and adds support for proxy configuration and custom headers.
// This struct provides methods to perform HTTP requests with optional proxy settings and custom headers.
//
// # Fields:
//   - client: A pointer to a `http.Client` instance used to perform HTTP requests.
//   - proxy: A `types.Proxy` struct containing proxy configuration details.
//   - headers: A `map[string]string` to store custom headers as key-value pairs.
//
// # Example:
//
//	proxyConfig := types.Proxy{
//		Ip:        "192.168.1.100",
//		Port:      1080,
//		Username:  "proxyuser",
//		Password:  "proxypass",
//		SocksType: 5,
//		Timeout:   15,
//	}
//	httpClient, err := NewHTTPClient(proxyConfig)
//	if err != nil {
//		log.Fatalf("Failed to initialize HTTP client: %v", err)
//	}
//
// # Example without proxy:
//
//	httpClient, err := NewHTTPClient(types.Proxy{})
//	if err != nil {
//		log.Fatalf("Failed to initialize HTTP client: %v", err)
//	}
//
// # Notes:
//   - Ensure the proxy server is reachable and properly configured when using a proxy.
//   - Timeout is set to 10 seconds by default but can be adjusted using the `Timeout` field in `proxyConfig`.
//   - SOCKS4 proxies are not supported in this implementation.
//
// # Errors:
//   - Returns an error if an invalid SOCKS type is specified.
//   - Returns an error if a SOCKS dialer cannot be created (e.g., invalid proxy address or credentials).
type HTTPClient struct {
	client  *http.Client
	proxy   types.Proxy
	headers map[string]string
}

// NewHTTPClient initializes and returns a new HTTP client, optionally configured to use a SOCKS proxy.
//
// This function supports SOCKS5 proxies with or without authentication. If proxy settings
// are not provided or incomplete, the client will operate without a proxy.
//
// # Parameters:
//   - proxyConfig: A `types.Proxy` struct containing proxy configuration details, including
//     the proxy server's IP address, port, authentication credentials, SOCKS type, and timeout.
//
// # Returns:
//   - *HTTPClient: A pointer to an `HTTPClient` instance with the configured HTTP client.
//   - error: An error if the proxy configuration is invalid, the SOCKS type is unsupported,
//     or the proxy dialer cannot be initialized.
//
// # Supported SOCKS Types:
//   - SOCKS5: Fully supported, with optional username and password authentication.
//   - SOCKS4: Not supported. Returns an error if specified.
//
// # Example:
//
//	proxyConfig := types.Proxy{
//		Ip:        "192.168.1.100",
//		Port:      1080,
//		Username:  "proxyuser",
//		Password:  "proxypass",
//		SocksType: 5,
//		Timeout:   15,
//	}
//	httpClient, err := NewHTTPClient(proxyConfig)
//	if err != nil {
//		log.Fatalf("Failed to initialize HTTP client: %v", err)
//	}
//
// # Example without proxy:
//
//	httpClient, err := NewHTTPClient(types.Proxy{})
//	if err != nil {
//		log.Fatalf("Failed to initialize HTTP client: %v", err)
//	}
//
// # Notes:
//   - Ensure the proxy server is reachable and properly configured when using a proxy.
//   - Timeout is set to 10 seconds by default but can be adjusted using the `Timeout` field in `proxyConfig`.
//   - SOCKS4 proxies are not supported in this implementation.
//
// # Errors:
//   - Returns an error if an invalid SOCKS type is specified.
//   - Returns an error if a SOCKS dialer cannot be created (e.g., invalid proxy address or credentials).
func NewHTTPClient(proxyConfig types.Proxy) (*HTTPClient, error) {
	var transport *http.Transport
	if proxyConfig.Ip != "" && proxyConfig.Port > 0 {
		proxyAddress := fmt.Sprintf("%s:%d", proxyConfig.Ip, proxyConfig.Port)
		var dialer proxy.Dialer
		var err error
		if proxyConfig.SocksType == 5 {
			if proxyConfig.Username != "" && proxyConfig.Password != "" {
				auth := proxy.Auth{
					User:     proxyConfig.Username,
					Password: proxyConfig.Password,
				}
				dialer, err = proxy.SOCKS5("tcp", proxyAddress, &auth, proxy.Direct)
			} else {
				dialer, err = proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
			}
		} else if proxyConfig.SocksType == 4 {
			return nil, fmt.Errorf("SOCKS4 proxy is not supported in this implementation")
		} else {
			return nil, fmt.Errorf("invalid SOCKS type: %d", proxyConfig.SocksType)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to create SOCKS dialer: %v", err)
		}
		transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			},
		}
	} else {
		transport = &http.Transport{}
	}
	timeout := 10 * time.Second
	if proxyConfig.Timeout > 0 {
		timeout = time.Duration(proxyConfig.Timeout) * time.Second
	}
	return &HTTPClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
	}, nil
}

// TODO: Add random headers to the HTTP client

// DoRequest sends an HTTP request with the specified method, URL, body, and additional headers.
func (httpClient *HTTPClient) DoRequest(method, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for key, value := range httpClient.headers {
		req.Header.Set(key, value)
	}
	resp, err := httpClient.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(responseBody))
	}
	return resp, nil
}

// Get performs a GET request to the specified URL with optional headers.
//
// This method sends an HTTP GET request to the provided URL. Additional headers can be
// specified in the headers map, which will be included in the request.
//
// # Parameters:
//   - url: The URL to which the GET request will be sent.
//
// # Returns:
//   - *http.Response: The HTTP response received from the server.
//   - error: An error if the request fails or the response status code indicates a failure.
func (httpClient *HTTPClient) Get(url string) (*http.Response, error) {
	return httpClient.DoRequest(http.MethodGet, url, nil)
}

// Post performs a POST request to the specified URL with a body and optional headers.
//
// This method sends an HTTP POST request to the provided URL with the specified body.
// Additional headers can be specified in the headers map, which will be included in the request.
//
// # Parameters:
//   - url: The URL to which the POST request will be sent.
//   - body: The body of the POST request, provided as a byte slice.
//   - headers: A map of additional headers to include in the request.
//
// # Returns:
//   - *http.Response: The HTTP response received from the server.
//   - error: An error if the request fails or the response status code indicates a failure.
func (httpClient *HTTPClient) Post(url string, body []byte) (*http.Response, error) {
	return httpClient.DoRequest(http.MethodPost, url, body)
}

// ReadResponseBody reads and returns the response body parsed as a JSON object or as a string if unmarshalling fails.
//
// This function reads the HTTP response body and tries to unmarshal it into the provided
// interface{} parameter. If unmarshalling fails, it returns the body as a string.
//
// # Parameters:
//   - resp: The HTTP response from which the body will be read.
//   - v: A pointer to the variable where the unmarshalled JSON object will be stored.
//
// # Returns:
//   - error: An error if reading the response body fails.
func ReadResponseBody(resp *http.Response, v interface{}) (string, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return string(body), nil
	}
	return "", nil
}
