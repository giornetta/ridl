package ridl

import (
	"fmt"
	"strings"
	"time"

	"github.com/giornetta/ridl/repository"

	"github.com/giornetta/ridl/cipher"
)

// Service defines the methods of the ridl service
type Service interface {
	GetRiddle(req *GetRequest) (*GetResponse, error)
	Encrypt(req *EncryptRequest) (*EncryptResponse, error)
	Decrypt(req *DecryptRequest) (*DecryptResponse, error)
}

type GetRequest struct {
	RiddleID string `json:"riddleID"`
}

type GetResponse struct {
	Question string `json:"question"`
}

// EncryptRequest contains the fields required to encrypt a riddle message
type EncryptRequest struct {
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	Message      string `json:"message"`
	IgnoreCase   bool   `json:"ignoreCase"`
	IgnoreSpaces bool   `json:"ignoreSpaces"`
	Expiry       string `json:"expiry"`
}

// EncryptResponse is the response of a successful encryption
type EncryptResponse struct {
	RiddleID string `json:"riddleID"`
}

// DecryptRequest contains the fields needed to decrypt a riddle message
type DecryptRequest struct {
	RiddleID string `json:"riddleID"`
	Answer   string `json:"answer"`
}

// DecryptResponse is the response to a successful decryption
type DecryptResponse struct {
	Message string `json:"message"`
}

type service struct {
	c    cipher.Cipher
	repo repository.Repository
}

// NewService returns an implementation of Service
func NewService(c cipher.Cipher, repo repository.Repository) Service {
	return &service{
		c:    c,
		repo: repo,
	}
}

func (s *service) GetRiddle(req *GetRequest) (*GetResponse, error) {
	r, err := s.repo.Get(req.RiddleID)
	if err != nil {
		return nil, err
	}

	fmt.Println(r.Question)
	return &GetResponse{
		Question: r.Question,
	}, nil
}

func (s *service) Encrypt(req *EncryptRequest) (*EncryptResponse, error) {
	if req.IgnoreCase {
		req.Answer = strings.ToLower(req.Answer)
	}

	if req.IgnoreSpaces {
		req.Answer = strings.Replace(req.Answer, " ", "", -1)
	}

	// Encrypt the message using the Riddle's answer as the key
	crypted, err := s.c.Encrypt([]byte(req.Message), []byte(req.Answer))
	if err != nil {
		return nil, err
	}

	id, err := s.repo.Put(&repository.Riddle{
		Question:     req.Question,
		Crypted:      crypted,
		IgnoreCase:   req.IgnoreCase,
		IgnoreSpaces: req.IgnoreSpaces,
		Expiry:       time.Now().Add(time.Hour * 2),
	})
	if err != nil {
		return nil, err
	}

	// Send back the response
	return &EncryptResponse{
		RiddleID: id,
	}, nil
}

func (s *service) Decrypt(req *DecryptRequest) (*DecryptResponse, error) {
	riddle, err := s.repo.Get(req.RiddleID)
	if err != nil {
		return nil, err
	}

	if riddle.IgnoreCase {
		req.Answer = strings.ToLower(req.Answer)
	}

	if riddle.IgnoreSpaces {
		req.Answer = strings.Replace(req.Answer, " ", "", -1)
	}

	// Decrypt the message using the given answer as the key
	decrypted, err := s.c.Decrypt(riddle.Crypted, []byte(req.Answer))
	if err != nil {
		return nil, err
	}

	// Send back the response
	return &DecryptResponse{
		Message: string(decrypted),
	}, nil
}
