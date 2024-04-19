package power

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/run-bigpig/jb-active/internal/cert"
	"github.com/run-bigpig/jb-active/internal/utils"
	"math/big"
	"os"
	"path/filepath"
)

func GenerateEqualResult() error {
	crt, err := cert.ParseCertPem()
	if err != nil {
		return err
	}
	x := new(big.Int).SetBytes(crt.Signature)
	y := 65537
	block, _ := pem.Decode([]byte(RootCert))
	rootCertificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	p, _ := rootCertificate.PublicKey.(*rsa.PublicKey)
	z := p.N
	zp, _ := crt.PublicKey.(*rsa.PublicKey)
	r := new(big.Int)
	r.Exp(x, big.NewInt(int64(y)), zp.N)
	equal := fmt.Sprintf("EQUAL,%d,%d,%d->%d", x, y, z, r)
	powerConf := filepath.Join(utils.GetStaticPath(), "conf", "power.conf")
	var wirterStr = `[Result]
;active code 
` + equal
	return os.WriteFile(powerConf, []byte(wirterStr), 0644)
}
