package service

import (
	"fmt"
	"strings"

	"github.com/giornetta/ridl"

	"github.com/giornetta/ridl/cipher"
)

type service struct {
	c    cipher.Cipher
	repo ridl.Repository
}

// New returns an implementation of ridl.Service.
// This implementation will work with every other implementation of the Cipher and Repository interfaces.
func New(c cipher.Cipher, repo ridl.Repository) ridl.Service {
	return &service{
		c:    c,
		repo: repo,
	}
}

func (s *service) GetRiddle(req *ridl.GetRequest) (*ridl.GetResponse, error) {
	r, err := s.repo.Get(req.RiddleID)
	if err != nil {
		return nil, err
	}

	fmt.Println(r.Question)
	return &ridl.GetResponse{
		Question: r.Question,
	}, nil
}

func (s *service) Encrypt(req *ridl.EncryptRequest) (*ridl.EncryptResponse, error) {
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

	exp := ridl.ExpiryToTime(req.Expiry)

	id, err := s.repo.Put(&ridl.Riddle{
		Question:     req.Question,
		Crypted:      crypted,
		IgnoreCase:   req.IgnoreCase,
		IgnoreSpaces: req.IgnoreSpaces,
		Expiry:       exp,
	})
	if err != nil {
		return nil, err
	}

	// Send back the response
	return &ridl.EncryptResponse{
		RiddleID: id,
	}, nil
}

func (s *service) Decrypt(req *ridl.DecryptRequest) (*ridl.DecryptResponse, error) {
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
	return &ridl.DecryptResponse{
		Message: string(decrypted),
	}, nil
}
