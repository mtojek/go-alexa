package alexa

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidCertURL    = errors.New("invalid certificate URL")
	ErrDecodingCert      = errors.New("cannot parse certificate PEM")
	ErrAmazonCertExpired = errors.New("Amazon certificate expired")
	ErrAmazonCertInvalid = errors.New("Amazon certificate invalid")
	ErrSignatureMismatch = errors.New("signature match failed")
)

// ValidateAmazonRequest runs all the mandatory Amazon security checks on the request.
func ValidateAmazonRequest(r *http.Request) error {
	// Check for development mode flag
	devMode := r.URL.Query().Get("_dev") != ""

	// Get the certificate URL and verify it
	certURL := r.Header.Get("SignatureCertChainUrl")
	if !verifyCertURL(certURL) && !devMode {
		return ErrInvalidCertURL
	}

	// Fetch certificate data
	certContents, err := readCert(certURL)
	if err != nil && !devMode {
		return err
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil && !devMode {
		return ErrDecodingCert
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil && !devMode {
		return err
	}

	// Check the certificate date
	certExpired := time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix()
	if certExpired && !devMode {
		return ErrAmazonCertExpired
	}

	// Check the certificate alternate names
	foundName := false
	for _, altName := range cert.Subject.Names {
		if altName.Value == "echo-api.amazon.com" {
			foundName = true
		}
	}

	if !foundName && !devMode {
		return ErrAmazonCertInvalid
	}

	// Verify the key
	publicKey := cert.PublicKey
	encryptedSig, _ := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))

	// Make the request body SHA1 and verify the request with the public key
	bodyBuf := new(bytes.Buffer)
	hash := sha1.New()
	_, err = io.Copy(hash, io.TeeReader(r.Body, bodyBuf))
	if err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(bodyBuf)

	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSig)
	if err != nil && !devMode {
		return ErrSignatureMismatch
	}

	return nil
}

func readCert(certURL string) ([]byte, error) {
	cert, err := http.Get(certURL)
	if err != nil {
		return nil, errors.New("Could not download Amazon cert file.")
	}
	defer cert.Body.Close()
	certContents, err := ioutil.ReadAll(cert.Body)
	if err != nil {
		return nil, errors.New("Could not read Amazon cert file.")
	}

	return certContents, nil
}

func verifyCertURL(path string) bool {
	if !strings.HasSuffix(path, "/echo.api/echo-api-cert.pem") {
		return false
	}

	if !strings.HasPrefix(path, "https://s3.amazonaws.com/echo.api/") && !strings.HasPrefix(path, "https://s3.amazonaws.com:443/echo.api/") {
		return false
	}

	return true
}
