package ridl

import (
	"encoding/base64"

	"github.com/giornetta/ridl/crypto"
)

// Service defines the methods of the ridl service
type Service interface {
	Encrypt(req *EncryptRequest) (*EncryptResponse, error)
	Decrypt(req *DecryptRequest) (*DecryptResponse, error)
}

// EncryptRequest contains the fields required to encrypt a riddle message
type EncryptRequest struct {
	Riddle  string `json:"riddle"`
	Answer  string `json:"answer"`
	Message string `json:"message"`
}

// EncryptResponse is the response of a successful encryption
type EncryptResponse struct {
	Riddle  string `json:"riddle"`
	Message string `json:"message"`
}

// DecryptRequest contains the fields needed to decrypt a riddle message
type DecryptRequest struct {
	Message string `json:"message"`
	Answer  string `json:"answer"`
}

// DecryptResponse is the response to a successful decryption
type DecryptResponse struct {
	Message string `json:"message"`
}

type service struct{}

// New returns an implementation of Service
func New() Service {
	return &service{}
}

func (s *service) Encrypt(req *EncryptRequest) (*EncryptResponse, error) {
	// Encrypt the message using the Riddle's answer as the key
	crypted, err := crypto.Encrypt([]byte(req.Message), []byte(req.Answer))
	if err != nil {
		return nil, err
	}

	// Encode the result as Base64
	message := base64.StdEncoding.EncodeToString(crypted)

	// Send back the response
	return &EncryptResponse{
		Riddle:  req.Riddle,
		Message: message,
	}, nil
}

func (s *service) Decrypt(req *DecryptRequest) (*DecryptResponse, error) {
	// Decode the given encrypted message from Base64
	message, err := base64.StdEncoding.DecodeString(req.Message)
	if err != nil {
		return nil, err
	}

	// Decrypt the message using the given answer as the key
	decrypted, err := crypto.Decrypt(message, []byte(req.Answer))
	if err != nil {
		return nil, err
	}

	// Send back the response
	return &DecryptResponse{
		Message: string(decrypted),
	}, nil
}
