package license

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/run-bigpig/jb-active/internal/cert"
	"github.com/run-bigpig/jb-active/internal/model"
	"math/big"
)

func generateLicenseID() string {
	const allowedCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const licenseLength = 10
	b := make([]byte, licenseLength)
	for i := range b {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allowedCharacters))))
		b[i] = allowedCharacters[index.Int64()]
	}
	return string(b)
}

func GenerateLicense(license *model.License) (string, error) {
	crt, err := cert.ParseCertPem()
	if err != nil {
		return "", err
	}
	privateKey, err := cert.ParseCertKey()
	if err != nil {
		return "", err
	}
	license.LicenseID = generateLicenseID()
	licenseStr, _ := json.Marshal(license)
	hashed := sha1.Sum(licenseStr)
	signature, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashed[:])

	licensePartBase64 := base64.StdEncoding.EncodeToString(licenseStr)
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	crtBase64 := base64.StdEncoding.EncodeToString(crt.Raw)
	return fmt.Sprintf("%s-%s-%s-%s", license.LicenseID, licensePartBase64, signatureBase64, crtBase64), nil
}
