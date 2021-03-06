package middleware

import (
	"github.com/awakesecurity/csp"
	"github.com/rs/cors"
	"github.com/unrolled/secure"

	"net/http"

	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"
)

func Init(h http.Handler) http.Handler {
	n := negroni.New()
	initLogger(n)
	initRecovery(n)
	initCORS(n)
	initSecure(n)
	// Only In Debug Mode
	// initDelay(n)
	initCSP(n)
	n.UseHandler(h)
	return n
}

func initCORS(n *negroni.Negroni) {
	c := cors.New(
		cors.Options{
			AllowedOrigins: []string{"*"},
		},
	)
	n.Use(c)
}

func initLogger(n *negroni.Negroni) {
	l := negronilogrus.NewMiddleware()
	n.Use(l)
}

func initRecovery(n *negroni.Negroni) {
	r := negroni.NewRecovery()
	n.Use(r)
}

func initSecure(n *negroni.Negroni) {
	secure := secure.New(secure.Options{
		// AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
		AllowedHosts: []string{},
		// HostsProxyHeaders is a set of header keys that may hold a proxied hostname value for the request.
		HostsProxyHeaders: []string{"X-Forwarded-Hosts"},
		// If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
		SSLRedirect: false,
		// If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
		SSLTemporaryRedirect: false,
		// SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
		SSLHost: "ssl.example.com",
		// SSLHostFunc is a function pointer, the return value of the function is the host name that has same functionality as `SSHost`. Default is nil. If SSLHostFunc is nil, the `SSLHost` option will be used.
		SSLHostFunc: nil,
		// SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		// STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSSeconds: 315360000,
		// If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSIncludeSubdomains: true,
		// If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload: true,
		// STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
		ForceSTSHeader: false,
		// If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		FrameDeny: true,
		// CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option. Default is "".
		CustomFrameOptionsValue: "SAMEORIGIN",
		// If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		ContentTypeNosniff: true,
		// If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		BrowserXssFilter: true,
		// CustomBrowserXssValue allows the X-XSS-Protection header value to be set with a custom value. This overrides the BrowserXssFilter option. Default is "".
		CustomBrowserXssValue: "1; report=https://example.com/xss-report",
		// ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "". Passing a template string will replace `$NONCE` with a dynamic nonce value of 16 bytes for each request which can be later retrieved using the Nonce function.
		ContentSecurityPolicy: "default-src 'self'",
		// PublicKey implements HPKP to prevent MITM attacks with forged certificates. Default is "".
		PublicKey: `pin-sha256="base64+primary=="; pin-sha256="base64+backup=="; max-age=5184000; includeSubdomains; report-uri="https://www.example.com/hpkp-report"`,
		// ReferrerPolicy allows the Referrer-Policy header with the value to be set with a custom value. Default is "".
		ReferrerPolicy: "same-origin",
		// FeaturePolicy allows the Feature-Policy header with the value to be set with a custom value. Default is "".
		FeaturePolicy:  "vibrate 'none';",
		ExpectCTHeader: `enforce, max-age=30, report-uri="https://www.example.com/ct-report"`,
		IsDevelopment:  true, // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	})

	s := negroni.HandlerFunc(secure.HandlerFuncWithNext)
	n.Use(s)
}

func initCSP(n *negroni.Negroni) {
	c := csp.New(
		csp.Config{
			Default: csp.None,
			Script:  csp.Self,
			Connect: csp.Self,
			Img:     csp.Self,
			Style:   csp.Self,
		},
	)
	n.Use(c.NegroniHandlerFunc())
}

// func InitJWT() negroni.Handler {
// 	options := jwtmiddleware.Options{
// 		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
// 			return []byte("JWT Secret"), nil
// 		},
// 		SigningMethod: jwt.SigningMethodHS256,
// 	}

// 	j := jwtmiddleware.New(options).HandlerWithNext
// 	return negroni.HandlerFunc(j)
// }
