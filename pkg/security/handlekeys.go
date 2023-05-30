package security

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "os"
    "io/fs"

    "github.com/lex0c/openet/pkg/config"
)

func GenerateKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, config.KeySizeInBits)

    if err != nil {
        return nil, nil, err
    }

    return privateKey, &privateKey.PublicKey, nil
}

func WritePrivateKeyToFile(privateKey *rsa.PrivateKey, filepath string) error {
    privateKeyBlock := &pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    }

    privateKeyPem := pem.EncodeToMemory(privateKeyBlock)

    if err := os.WriteFile(filepath, privateKeyPem, fs.FileMode(config.KeyPermMode)); err != nil {
        return fmt.Errorf("Error writing private key to file: %w", err)
    }

    return nil
}

func ReadPrivateKeyFromFile(filepath string) (*rsa.PrivateKey, error) {
    privateKeyPem, err := os.ReadFile(filepath)

    if err != nil {
        return nil, fmt.Errorf("Error reading private key file: %w", err)
    }

    privateKeyBlock, _ := pem.Decode(privateKeyPem)

    if privateKeyBlock == nil {
        return nil, fmt.Errorf("Error decoding PEM: %w", err)
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)

    if err != nil {
        return nil, fmt.Errorf("Error interpreting the private key: %w", err)
    }

    return privateKey, nil
}

func ReadPublicKeyFromPrivateKey(privateKey *rsa.PrivateKey) *rsa.PublicKey {
    return &privateKey.PublicKey
}

func EncodePublicKey(publicKey *rsa.PublicKey) (string, error) {
    publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)

    if err != nil {
        return "", err
    }

    publicKeyBlock := &pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: publicKeyBytes,
    }

    return string(pem.EncodeToMemory(publicKeyBlock)), nil
}

func DecodePublicKey(publicKeyPem string) (*rsa.PublicKey, error) {
    publicKeyPemDecoded, _ := pem.Decode([]byte(publicKeyPem))

    if publicKeyPemDecoded == nil {
        return nil, fmt.Errorf("Error decoding PEM public key")
    }

    publicKeyDecoded, err := x509.ParsePKIXPublicKey(publicKeyPemDecoded.Bytes)

    if err != nil {
        return nil, err
    }

    rsaPublicKey, ok := publicKeyDecoded.(*rsa.PublicKey)

    if !ok {
        return nil, fmt.Errorf("Not a valid RSA public key")
    }

    return rsaPublicKey, nil
}

func ValidatePublicKeyPem(publicKeyPem string) error {
    _, err := DecodePublicKey(publicKeyPem)
    return err
}

func GetMyAddress(privateKey *rsa.PrivateKey) (string, error) {
    pubKey := ReadPublicKeyFromPrivateKey(privateKey)

    addr, err := EncodePublicKey(pubKey)

    if err != nil {
        return "", fmt.Errorf("An error occurred: %v", err)
    }

    return addr, nil
}

