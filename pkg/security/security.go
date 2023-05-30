package security

import (
    "fmt"
    "crypto/rsa"

    "github.com/lex0c/openet/pkg/config"
)

type Security interface {
    GenerateKeys() (*rsa.PrivateKey, *rsa.PublicKey, error)
    WritePrivateKeyToFile(privateKey *rsa.PrivateKey, filepath string) error
    ReadPrivateKeyFromFile(filepath string) (*rsa.PrivateKey, error)
    ReadPublicKeyFromPrivateKey(privateKey *rsa.PrivateKey) *rsa.PublicKey
    EncodePublicKey(publicKey *rsa.PublicKey) (string, error)
    DecodePublicKey(publicKeyPem string) (*rsa.PublicKey, error)
    ValidatePublicKeyPem(publicKeyPem string) error
    GetMyAddress() (string, error)
}

func getMyKeyPath() string {
    return fmt.Sprintf("%s/%s", config.DefaultFileDir, config.KeyFileName)
}

func GetMyAddress() (string, error) {
    privKey, err := ReadPrivateKeyFromFile(getMyKeyPath())

    if err != nil {
        return "", fmt.Errorf("An error occurred: %v", err)
    }

    pubKey := ReadPublicKeyFromPrivateKey(privKey)

    addr, err := EncodePublicKey(pubKey)

    if err != nil {
        return "", fmt.Errorf("An error occurred: %v", err)
    }

    return addr, nil
}

