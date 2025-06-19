package samlsp

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
)

// NewClassLinkMiddleware returns a samlsp.Middleware for one ClassLink tenant.
//
//	subdomain   – e.g. "mydistrict"
//	publicURL   – root URL of *your* service (scheme + host, no trailing slash)
//	key         – RSA private key (nil ⇒ generated dev key)
//	certPEM     – matching cert in PEM (nil ⇒ self-signed dev cert)
func NewClassLinkMiddleware(
	subdomain string,
	publicURL *url.URL,
	key *rsa.PrivateKey,
	certPEM []byte,
) (*samlsp.Middleware, error) {

	/* ---------- 1. Dev key / cert if not provided ---------- */

	if key == nil {
		var err error
		key, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
	}

	var cert *x509.Certificate
	if len(certPEM) != 0 {
		block, _ := pem.Decode(certPEM)
		if block == nil {
			return nil, fmt.Errorf("invalid cert PEM")
		}
		cert, _ = x509.ParseCertificate(block.Bytes)
	}

	/* ---------- 2. Fetch IdP metadata ---------- */

	idpMetaURL := fmt.Sprintf(
		"https://%s.classlink.com/samlsso/saml/metadata",
		subdomain,
	)

	idpMeta, err := samlsp.FetchMetadata(
		context.Background(), http.DefaultClient, *mustURL(idpMetaURL),
	)
	if err != nil {
		return nil, err
	}

	/* ---------- 3. Build samlsp.Options ---------- */

	opts := samlsp.Options{
		URL:               *publicURL, // e.g. https://svc.com
		Key:               key,
		Certificate:       cert,
		IDPMetadata:       idpMeta,
		AllowIDPInitiated: false,
	}

	return samlsp.New(opts)
}

func mustURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}
