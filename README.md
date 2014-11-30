# Single Page Application OAuth Proxy

This little proxy allow you to use OAuth v2 [Resource Owner Password Credentials Grant](http://tools.ietf.org/html/draft-ietf-oauth-v2-31#section-4.3) (aka. Password Grant) with a Single Page Application.
It acts as the Client to the Authorization Server. The application only sees an encrypted token and need to send it back to the proxy as an header.

The proxy supports the following endpoints:

 * `/auth`: negociate the Access Token
 * `/auth/refresh`: ask for a new Access Token if a Refresh Token is provided
 * `/`: proxy all requests back to the API
 * `/ping`: ensure the proxy is alive

This is heavily based on Alex Bilbie [thoughts](http://alexbilbie.com/2014/11/oauth-and-javascript/).

## Usage

First build the application for your target plateform, here I am targeting a 64-bit Linux distribution:

    $ GOARCH=amd64 GOOS=linux go build -o spa_oauth_proxy *.go

Then pass required parameters:

    $ ./spa_oauth_proxy -access-token-url="https://example.com/api/v1/oauth/access-token" -client-id="clientid" -client-secret="clientsecret" -http-address="0.0.0.0:3033" -key="a very very very very long keyss"

You can print the command line help:

    $ ./spa_oauth_proxy -h
    Usage of ./spa_oauth_proxy:
      -access-token-url="": The api endpoint used to create access token and refresh the token.
      -base-path="": The base path of the proxy.
      -client-id="": The OAuth client id.
      -client-secret="": The OAuth client secret.
      -header="X-Auth": The encrypted token header name.
      -http-address="127.0.0.1:8080": The <addr>:<port> to listen on for HTTP clients..
      -key="": The 32 bytes encryption key.
      -pong="pong": The ping response.
      -version=false: Print the version and exit.

I am currently building a Single Page Application with Laravel and AngularJS as an example.

## License

(The MIT license)

Copyright (c) 2014 Kevin Le Brun <lebrun.k@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
