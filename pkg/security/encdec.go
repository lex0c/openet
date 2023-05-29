package security

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "fmt"
)

func EncryptMessage(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
    return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

func DecryptMessage(privateKey *rsa.PrivateKey, encryptedMessage []byte) ([]byte, error) {
    return rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedMessage)
}

func SignMessage(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
    hasher := sha256.New()

    if _, err := hasher.Write(message); err != nil {
        return nil, fmt.Errorf("Error hashing the message: %w", err)
    }

    hash := hasher.Sum(nil)

    signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)

    if err != nil {
        return nil, fmt.Errorf("Error signing the message: %w", err)
    }

    return signature, nil
}

func VerifySignature(publicKey *rsa.PublicKey, message []byte, signature []byte) error {
    hasher := sha256.New()

    if _, err := hasher.Write(message); err != nil {
        return fmt.Errorf("Error hashing the message: %w", err)
    }

    hash := hasher.Sum(nil)

    if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, signature); err != nil {
        return fmt.Errorf("Error verifying the message: %w", err)
    }

    return nil
}

