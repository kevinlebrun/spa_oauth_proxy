package main

// TODO handle errors if no Refresh Token is given
// TODO ensure proper error handling
// TODO replace all ReverseProxy with a custom one

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type OAuthAPI struct {
	AccessTokenURL *url.URL
	ClientId       string
	ClientSecret   string
}

type Proxy struct {
	HeaderName string
	BasePath   string
	Key        []byte
	OAuthAPI   OAuthAPI
	ServeMux   *http.ServeMux
}

func (p Proxy) ReverseProxyHandler() *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			start := len(p.BasePath)
			req.URL.Path = req.URL.Path[start:]
			req.URL.Scheme = p.OAuthAPI.AccessTokenURL.Scheme
			req.URL.Host = p.OAuthAPI.AccessTokenURL.Host

			req.Header.Add("Authorization", "Bearer invalid_token")

			token := req.Header.Get(p.HeaderName)
			if token != "" {
				at, err := DecodeAccessToken(p.Key, token)
				if err == nil {
					req.Header.Set("Authorization", "Bearer "+at.AccessToken)
				}
			}
		},
	}
}

func (p Proxy) AuthHandler() *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			var credentials = struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}{}

			err := json.NewDecoder(req.Body).Decode(&credentials)
			if err != nil {
				log.Println("Something is wrong with the body, should return 4xx error.")
			}
			req.Body.Close()

			values := &url.Values{}
			values.Set("client_id", p.OAuthAPI.ClientId)
			values.Set("client_secret", p.OAuthAPI.ClientSecret)
			values.Set("username", credentials.Username)
			values.Set("password", credentials.Password)
			values.Set("grant_type", "password")

			req.URL = p.OAuthAPI.AccessTokenURL
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Body, req.ContentLength = createFormBodyFromValues(values)
		},
		Transport: &TokenTransport{Key: p.Key},
	}
}

func (p Proxy) RefreshAuthHandler() *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			values := &url.Values{}
			values.Set("client_id", p.OAuthAPI.ClientId)
			values.Set("client_secret", p.OAuthAPI.ClientSecret)
			values.Set("grant_type", "refresh_token")
			values.Set("refresh_token", "invalid_token")

			token := req.Header.Get(p.HeaderName)
			if token != "" {
				at, err := DecodeAccessToken(p.Key, token)
				if err == nil {
					values.Set("refresh_token", at.RefreshToken)
				}
			}

			req.URL = p.OAuthAPI.AccessTokenURL
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Body, req.ContentLength = createFormBodyFromValues(values)
		},
		Transport: &TokenTransport{Key: p.Key},
	}
}

func createFormBodyFromValues(values *url.Values) (io.ReadCloser, int64) {
	data := values.Encode()
	return ioutil.NopCloser(bytes.NewBufferString(data)), int64(len(data))
}

func (p Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.ServeMux.ServeHTTP(w, req)
}

func (p Proxy) Handle(pattern string, handler http.Handler) {
	p.ServeMux.Handle(pattern, handler)
}

func (p Proxy) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	p.ServeMux.HandleFunc(pattern, handler)
}

type TokenTransport struct {
	Key []byte
}

func (t *TokenTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)

	if response.StatusCode >= 400 {
		response.Body = ioutil.NopCloser(bytes.NewBufferString(""))
		response.ContentLength = int64(0)
		response.StatusCode = 403
		response.Status = "403 Forbidden"
		return response, err
	}

	body, _ := ioutil.ReadAll(response.Body)

	var at AccessToken

	json.Unmarshal(body, &at)

	token, err := at.Encode(t.Key)
	if err != nil {
		return response, err
	}

	as := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	newBody, _ := json.Marshal(as)

	response.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))
	response.ContentLength = int64(len(body))

	return response, err
}
