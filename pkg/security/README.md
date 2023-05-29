# Security package

This package provides a set of utilities for working with RSA keys and performing common RSA operations.

## Functionality

Here are the main functions provided in the `security` package:

1. **GenerateKeys**: This function generates an RSA key pair.

2. **WritePrivateKeyToFile**: This function writes a given private key to a file.

3. **ReadPrivateKeyFromFile**: This function reads a private key from a file.

4. **ReadPublicKeyFromPrivateKey**: This function derives a public key from a given private key.

5. **EncodePublicKey**: This function encodes a public key to a PEM string.

6. **DecodePublicKey**: This function decodes a PEM string to a public key.

7. **SignMessage**: This function signs a given message with a private key.

8. **VerifySignature**: This function verifies a given message's signature with a public key.

9. **ValidatePublicKeyPem**: This function will return `nil` if the public key string is properly formatted, or an error message if it is not.

## Usage

### Key Generation

Use the `GenerateKeys` function to generate a new RSA private key and its associated public key.

```go
privateKey, publicKey, err := security.GenerateKeys()
```

### Writing and Reading a Private Key to/from a File

The `WritePrivateKeyToFile` function can be used to write a private key to a file.

```go
err := security.WritePrivateKeyToFile(privateKey, "/path/to/file.pem")
```

To read the private key back from the file, use `ReadPrivateKeyFromFile`.

```go
privateKey, err := security.ReadPrivateKeyFromFile("/path/to/file.pem")
```

### Reading a Public Key from Private Key

The `ReadPublicKeyFromPrivateKey` function derives a public key from a given private key.

```go
publicKey := security.ReadPublicKeyFromPrivateKey(privateKey)
```

### Public Key Encoding and Decoding

Use `EncodePublicKey` to encode a public key to a string.

```go
publicKeyString, err := security.EncodePublicKey(publicKey)
```

To decode the public key from the string, use `DecodePublicKey`.

```go
publicKey, err := security.DecodePublicKey(publicKeyString)
```

### Message Encryption and Decryption

The `EncryptMessage` function can be used to encrypt a message using a public key.

```go
encryptedMessage, err := security.EncryptMessage(publicKey, []byte("Hello"))
```

To decrypt the message back to its original form using the private key, use `DecryptMessage`.

```go
message, err := security.DecryptMessage(privateKey, encryptedMessage)
```

### Message Signing and Signature Verification

You can sign a message using the `SignMessage` function. This creates a signature that can be used to verify the authenticity of the message.

```go
signature, err := security.SignMessage(privateKey, []byte("Hello"))
```

The `VerifySignature` function can be used to verify the authenticity of a message and its signature.

```go
err := security.VerifySignature(publicKey, []byte("Hello"), signature)
```

### Validate Public Key

The `ValidatePublicKeyPem` function will return `nil` if the public key string is properly formatted, or an error message if it is not.

```go
err := security.ValidatePublicKeyPem(publicKeyPem)
```

### Testing

```bash
go test -v ./...
```

