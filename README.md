# Single Page Application OAuth Proxy

This little proxy allow you to use OAuth v2 [Resource Owner Password Credentials Grant](http://tools.ietf.org/html/draft-ietf-oauth-v2-31#section-4.3) (aka. Password Grant) with a Single Page Application.
It acts as the Client to the Authorization Server. The application only sees an encrypted token and need to send it back to the proxy as an header.

The proxy supports the following endpoints:

 * `/auth`: negociate the Access Token
 * `/auth/refresh`: ask for a new Access Token if there is a Refresh Token
 * `/`: proxy all requests back to the API
 * `/ping`: ensure the proxy is alive

This is heavily based on Alex Bilbie [thoughts](http://alexbilbie.com/2014/11/oauth-and-javascript/).

## Usage

First build the application for your target plateform, here I am targeting a 64-bit Linux:

    $ GOARCH=amd64 GOOS=linux go build -o spa_oauth_proxy *.go

Then pass required parameters:

    $ ./spa_oauth_proxy -access-token-url="https://example.com/api/v1/oauth/access-token" -client-id="clientid" -client-secret="clientsecret" -http-address="0.0.0.0:3033" -key="a very very very very long keyss"

You can print the command line help:

    $ ./spa_oauth_proxy -h
    Usage of ./spa_oauth_proxy:
      -access-token-url="": The api endpoint used to create access token and refresh the token.
      -base-path="/proxy": The base path of the proxy.
      -client-id="": The OAuth client id.
      -client-secret="": The OAuth client secret.
      -header="X-Auth": The encrypted token header name.
      -http-address="127.0.0.1:8080": The <addr>:<port> to listen on for HTTP clients..
      -key="": The 32 bytes encryption key.
      -pong="pong": The ping response.
      -version=false: Print the version and exit.

I am currently building a Single Page Application with Laravel and AngularJS as an example.
