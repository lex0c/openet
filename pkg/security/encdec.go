package security

import (
    "crypto/rand"
    "crypto/rsa"
)

func EncryptMessage(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
    return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

func DecryptMessage(privateKey *rsa.PrivateKey, encryptedMessage []byte) ([]byte, error) {
    return rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedMessage)
}

