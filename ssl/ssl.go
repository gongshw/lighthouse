package ssl

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

const (
	DEFAULT_CRT_NAME = "temp.crt"
	DEFAULT_KEY_NAME = "temp.key"
	DEFAULT_KEY_BITS = 1024
)

func CreateCrtAndKey(hostName string, crtFile, keyFile *os.File) error {
	if crtFile == nil || keyFile == nil {
		return errors.New("crt or ley file not exist")
	}
	if privateKey, err := rsa.GenerateKey(rand.Reader, DEFAULT_KEY_BITS); err == nil {
		now := time.Now()
		serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
		crtTemplate := &x509.Certificate{
			SerialNumber: serialNumber,
			Subject: pkix.Name{
				Organization: []string{hostName},
			},
			NotBefore:             now,
			NotAfter:              now.Add(365 * 24 * time.Hour),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			BasicConstraintsValid: true,
		}
		publicKey := privateKey.Public()
		cert, err := x509.CreateCertificate(rand.Reader, crtTemplate, crtTemplate, publicKey, privateKey)
		if err != nil {
			return err
		}
		if writePem("CERTIFICATE", crtFile, cert) != nil {
			return err
		}
		if writePem("RSA PRIVATE KEY", keyFile, x509.MarshalPKCS1PrivateKey(privateKey)) != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

func writePem(pemType string, file *os.File, bytes []byte) error {
	writer := bufio.NewWriter(file)
	err := pem.Encode(writer, &pem.Block{
		Type:  pemType,
		Bytes: bytes,
	})
	if err != nil {
		return err
	} else {
		writer.Flush()
		return nil
	}
}

func CreateTempCrtAndKey(hostName string) (string, string, error) {
	var err error
	var tempCrtFile, tempKeyFile *os.File
	if tempCrtFile, err = ioutil.TempFile("", DEFAULT_CRT_NAME); err != nil {
		return "", "", err
	}
	if tempKeyFile, err = ioutil.TempFile("", DEFAULT_KEY_NAME); err != nil {
		return "", "", err
	}

	if err = CreateCrtAndKey(hostName, tempCrtFile, tempKeyFile); err == nil {
		return tempCrtFile.Name(), tempKeyFile.Name(), nil
	} else {
		os.Remove(tempCrtFile.Name())
		os.Remove(tempKeyFile.Name())
		return "", "", err
	}
}
