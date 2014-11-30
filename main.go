package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {

	var (
		showVersion    = flag.Bool("version", false, "Print the version and exit.")
		clientId       = flag.String("client-id", "", "The OAuth client id.")
		clientSecret   = flag.String("client-secret", "", "The OAuth client secret.")
		httpAddress    = flag.String("http-address", "127.0.0.1:8080", "The <addr>:<port> to listen on for HTTP clients..")
		pong           = flag.String("pong", "pong", "The ping response.")
		headerName     = flag.String("header", "X-Auth", "The encrypted token header name.")
		basePath       = flag.String("base-path", "", "The base path of the proxy.")
		accessTokenURL = flag.String("access-token-url", "", "The api endpoint used to create access token and refresh the token.")
		key            = flag.String("key", "", "The 32 bytes encryption key.")
	)

	flag.Parse()

	if *showVersion {
		fmt.Printf("spa_oauth_proxy v%s\n", VERSION)
		return
	}

	if len(*key) != 32 {
		log.Fatal("The key length should be exactly 32 bytes.")
	}

	apiURL, err := url.Parse(*accessTokenURL)
	if err != nil {
		log.Fatal("Failed to parse access token URL")
	}

	if *clientId == "" {
		log.Fatal("Invalid client id.")
	}

	if *clientSecret == "" {
		log.Fatal("Invalid client secret.")
	}

	if *pong == "" {
		log.Fatal("Invalid ping handler response.")
	}

	if *headerName == "" {
		log.Fatal("Invalid authentication header name.")
	}

	proxy := Proxy{
		HeaderName: *headerName,
		BasePath:   *basePath,
		Key:        []byte(*key),
		OAuthAPI: OAuthAPI{
			AccessTokenURL: apiURL,
			ClientId:       *clientId,
			ClientSecret:   *clientSecret,
		},
		ServeMux: http.NewServeMux(),
	}

	proxy.Handle(proxy.BasePath+"/", proxy.ReverseProxyHandler())
	proxy.Handle(proxy.BasePath+"/auth", proxy.AuthHandler())
	proxy.Handle(proxy.BasePath+"/auth/refresh", proxy.RefreshAuthHandler())

	proxy.Handle(proxy.BasePath+"/ping", &PingHandler{Response: *pong})

	log.Println("Listening on", *httpAddress)
	err = http.ListenAndServe(*httpAddress, proxy)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
