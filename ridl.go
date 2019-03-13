package ridl

import "time"

// Riddle represents a record inside the Database
type Riddle struct {
	Question     string
	Crypted      []byte
	IgnoreCase   bool
	IgnoreSpaces bool
	Expiry       time.Time
}

// IsExpired checks if a riddle is expired and it should be removed
func (r *Riddle) IsExpired() bool {
	return time.Now().After(r.Expiry)
}

// Repository defines which methods should be implemented by databases
type Repository interface {
	Put(r *Riddle) (string, error)
	Get(id string) (*Riddle, error)
	Delete(id string) error
	DeleteExpired() error
}

// Service defines the methods of the ridl service
type Service interface {
	GetRiddle(req *GetRequest) (*GetResponse, error)
	Encrypt(req *EncryptRequest) (*EncryptResponse, error)
	Decrypt(req *DecryptRequest) (*DecryptResponse, error)
}

// GetRequest contains the fields needed to retrieve a riddle.
type GetRequest struct {
	RiddleID string `json:"riddleID"`
}

// GetResponse is the response of a successful read request
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

// ExpiryToTime converts a string received in EncryptRequest to a Time stored in the Database
func ExpiryToTime(exp string) time.Time {
	switch exp {
	case "2h":
		return time.Now().Add(time.Hour * 2)
	case "6h":
		return time.Now().Add(time.Hour * 6)
	case "12h":
		return time.Now().Add(time.Hour * 12)
	default:
		return time.Now().Add(time.Hour * 24)
	}
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
