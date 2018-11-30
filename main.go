package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront/sign"
)

// Creates signed cookies to authorize private CloudFront access.
func main() {
	keyFile := os.Getenv("CLOUDFRONT_PRIVATE_KEY")
	if keyFile == "" {
		log.Fatalln("CLOUDFRONT_PRIVATE_KEY env variable not set")
	}
	keyID := os.Getenv("CLOUDFRONT_KEY_ID")
	if keyID == "" {
		log.Fatalln("CLOUDFRONT_KEY_ID env variable not set")
	}

	// Load the PEM file into memory so it can be used by the signer
	privKey, err := sign.LoadPEMPrivKeyFile(keyFile)
	if err != nil {
		log.Fatalf("Failed to load private key, %s\n", err)
	}

	// Create the new CookieSigner to get signed cookies for CloudFront
	// resource requests
	signer := sign.NewCookieSigner(keyID, privKey)

	// Get the cookies for the resource
	cookies, err := signer.Sign("http*://*", time.Now().Add(10*365*24*time.Hour)) // 10 years
	if err != nil {
		log.Fatalf("Failed to sign cookies, %s\n", err)
	}

	// Create a dummy request to serialize cookies
	req, err := http.NewRequest("GET", "", nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	fmt.Println(req.Header.Get("Cookie"))
}
