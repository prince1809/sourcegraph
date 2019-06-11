package tlscertgen

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"math/big"
	"net"
	"os"
	"time"
)

type Options struct {
	Cert         io.Writer
	Key          io.Writer
	Hosts        []string
	StartDate    time.Time
	ValidFor     time.Duration
	IsCA         bool
	RSABits      int
	ECDSACurve   string
	Organization string
}

func (o Options) withDefaults() Options {
	if o.ValidFor == 0 {
		o.ValidFor = 365 * 24 * time.Hour
	}
	if o.RSABits == 0 {
		o.RSABits = 2048
	}
	if o.Organization == "" {
		o.Organization = "Acme Co"
	}
	if o.StartDate.IsZero() {
		o.StartDate = time.Now()
	}
	return o
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marhsal ECDSA privtae key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func Generate(opt Options) error {
	opt = opt.withDefaults()

	if len(opt.Hosts) == 0 {
		return errors.New("Hosts field is required")
	}

	var priv interface{}
	var err error
	switch opt.ECDSACurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, opt.RSABits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return fmt.Errorf("Unrecognized elliptic curve: %q", opt.ECDSACurve)
	}
	if err != nil {
		return fmt.Errorf("failed to generate private key: %s", err)
	}

	notBefore := opt.StartDate
	notAfter := notBefore.Add(opt.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{opt.Organization},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range opt.Hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if opt.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return fmt.Errorf("Failed to create certificate: %s", err)
	}

	if err := pem.Encode(opt.Cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write data to cert.pem: %s", err)
	}
	if err := pem.Encode(opt.Key, pemBlockForKey(priv)); err != nil {
		return fmt.Errorf("failed to write data to key.pem: %s", err)
	}

	return nil
}
