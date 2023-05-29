package security

import (
    "testing"
    "os"
    "bytes"
    "fmt"
    "crypto/rand"
)

func TestGenerateKeys(t *testing.T) {
    privateKey, publicKey, err := GenerateKeys()

    if err != nil {
        t.Fatalf("Failed to generate keys: %v", err)
    }

    if privateKey == nil {
        t.Fatalf("Generated nil private key")
    }

    if publicKey == nil {
        t.Fatalf("Generated nil public key")
    }
}

func TestEncodeAndDecodePublicKey(t *testing.T) {
    _, publicKey, err := GenerateKeys()

    if err != nil {
        t.Fatalf("Failed to generate keys: %v", err)
    }

    encodedKey, err := EncodePublicKey(publicKey)

    if err != nil {
        t.Fatalf("Failed to encode public key: %v", err)
    }

    decodedKey, err := DecodePublicKey(encodedKey)

    if err != nil {
        t.Fatalf("Failed to decode public key: %v", err)
    }

    if publicKey.N.Cmp(decodedKey.N) != 0 || publicKey.E != decodedKey.E {
        t.Errorf("Original and decoded keys don't match")
    }
}

func TestWriteAndReadPrivateKeyToFile(t *testing.T) {
    privateKey, _, err := GenerateKeys()

    if err != nil {
        t.Fatalf("Failed to generate keys: %v", err)
    }

    filename := "test_key.pem"

    defer os.Remove(filename) // Clean up file after test

    if err = WritePrivateKeyToFile(privateKey, filename); err != nil {
        t.Fatalf("Failed to write private key to file: %v", err)
    }

    readKey, err := ReadPrivateKeyFromFile(filename)

    if err != nil {
        t.Fatalf("Failed to read private key from file: %v", err)
    }

    if privateKey.D.Cmp(readKey.D) != 0 || privateKey.PublicKey.N.Cmp(readKey.PublicKey.N) != 0 {
        t.Errorf("Original and read keys don't match")
    }
}

func TestWritePrivateKeyToFile_NonexistentDirectory(t *testing.T) {
    privateKey, _, err := GenerateKeys()

    if err != nil {
        t.Fatalf("Failed to generate keys: %v", err)
    }

    filename := "/nonexistent/directory/test_key.pem"

    err = WritePrivateKeyToFile(privateKey, filename)

    if err == nil {
        t.Errorf("Expected error when writing private key to nonexistent directory, got nil")
    }
}

func TestReadPrivateKeyFromFile_NonexistentFile(t *testing.T) {
    if _, err := ReadPrivateKeyFromFile("/nonexistent/file"); err == nil {
        t.Errorf("Expected error when reading from nonexistent file, got nil")
    }
}

func TestReadPrivateKeyFromFile_InvalidFile(t *testing.T) {
    filename := "invalid_key.pem"

    defer os.Remove(filename) // Clean up file after test

    if err := os.WriteFile(filename, []byte("invalid"), 0600); err != nil {
        t.Fatalf("Failed to write to test file: %v", err)
    }

    if _, err := ReadPrivateKeyFromFile(filename); err == nil {
        t.Errorf("Expected error when reading invalid private key file, got nil")
    }
}

func TestEncryptDecrypt(t *testing.T) {
    privateKey, publicKey, err := GenerateKeys()

		if err != nil {
        t.Fatal(fmt.Errorf("Error generating keys: %w", err))
		}

		message := make([]byte, 100)

		rand.Read(message)

		encryptedMessage, err := EncryptMessage(publicKey, message)

		if err != nil {
        t.Fatal(fmt.Errorf("Error encrypting message: %w", err))
		}

		// Tampered encrypted message should not decrypt correctly
		tamperedEncryptedMessage := append([]byte{}, encryptedMessage...)
		tamperedEncryptedMessage[0] ^= 0xff

		_, err = DecryptMessage(privateKey, tamperedEncryptedMessage)

		if err == nil {
        t.Fatal("Tampered encrypted message decrypted without error")
		}

		decryptedMessage, err := DecryptMessage(privateKey, encryptedMessage)

		if err != nil {
        t.Fatal(fmt.Errorf("Error decrypting message: %w", err))
		}

		if !bytes.Equal(decryptedMessage, message) {
        t.Fatalf("Decrypted message '%s' does not match original message '%s'", decryptedMessage, message)
		}
}

func TestSignVerify(t *testing.T) {
    privateKey, publicKey, err := GenerateKeys()

		if err != nil {
        t.Fatal(fmt.Errorf("Error generating keys: %w", err))
		}

		message := make([]byte, 100)
		rand.Read(message)

		signature, err := SignMessage(privateKey, message)

		if err != nil {
        t.Fatal(fmt.Errorf("Error signing message: %w", err))
		}

		// Tampered message should not verify correctly
		tamperedMessage := append([]byte{}, message...)
		tamperedMessage[0] ^= 0xff

		err = VerifySignature(publicKey, tamperedMessage, signature)

		if err == nil {
        t.Fatal("Tampered message verified without error")
		}

		// Tampered signature should not verify correctly
		tamperedSignature := append([]byte{}, signature...)
		tamperedSignature[0] ^= 0xff

		err = VerifySignature(publicKey, message, tamperedSignature)

		if err == nil {
        t.Fatal("Tampered signature verified without error")
		}

		err = VerifySignature(publicKey, message, signature)

		if err != nil {
        t.Fatal(fmt.Errorf("Error verifying signature: %w", err))
		}
}

