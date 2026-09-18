package lifecycle

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type lockOwner struct {
	PID       int       `json:"pid"`
	StartedAt time.Time `json:"started_at"`
	Token     string    `json:"token"`
}

type Lock struct {
	path     string
	token    string
	borrowed bool
}

func acquireLock(stateDir, borrowedToken string, now time.Time) (*Lock, error) {
	path := filepath.Join(stateDir, "lifecycle.lock")
	if borrowedToken != "" {
		owner, err := readLockOwner(path)
		if err != nil {
			return nil, fmt.Errorf("borrow lifecycle lock: %w", err)
		}
		if owner.Token != borrowedToken {
			return nil, fmt.Errorf("borrow lifecycle lock: token mismatch")
		}
		return &Lock{path: path, token: borrowedToken, borrowed: true}, nil
	}

	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return nil, fmt.Errorf("create state directory: %w", err)
	}
	if err := os.Mkdir(path, 0o700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("acquire lifecycle lock: %w", err)
		}
		owner, ownerErr := readLockOwner(path)
		if ownerErr == nil && !processAlive(owner.PID) {
			if removeErr := os.RemoveAll(path); removeErr != nil {
				return nil, fmt.Errorf("remove stale lifecycle lock: %w", removeErr)
			}
			return acquireLock(stateDir, "", now)
		}
		if ownerErr != nil {
			return nil, fmt.Errorf("lifecycle lock exists and owner is unreadable: %w", ownerErr)
		}
		return nil, fmt.Errorf("lifecycle operation already running with pid %d since %s", owner.PID, owner.StartedAt.Format(time.RFC3339))
	}

	tokenBytes := make([]byte, 24)
	if _, err := rand.Read(tokenBytes); err != nil {
		os.RemoveAll(path)
		return nil, fmt.Errorf("create lifecycle lock token: %w", err)
	}
	owner := lockOwner{PID: os.Getpid(), StartedAt: now.UTC(), Token: hex.EncodeToString(tokenBytes)}
	data, err := json.Marshal(owner)
	if err != nil {
		os.RemoveAll(path)
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(path, "owner.json"), data, 0o600); err != nil {
		os.RemoveAll(path)
		return nil, fmt.Errorf("write lifecycle lock owner: %w", err)
	}
	return &Lock{path: path, token: owner.Token}, nil
}

func readLockOwner(path string) (lockOwner, error) {
	data, err := os.ReadFile(filepath.Join(path, "owner.json"))
	if err != nil {
		return lockOwner{}, err
	}
	var owner lockOwner
	if err := json.Unmarshal(data, &owner); err != nil {
		return lockOwner{}, err
	}
	if owner.PID <= 0 || owner.Token == "" || owner.StartedAt.IsZero() {
		return lockOwner{}, fmt.Errorf("invalid lock owner")
	}
	return owner, nil
}

func (l *Lock) Token() string { return l.token }

func (l *Lock) Release() error {
	if l == nil || l.borrowed {
		return nil
	}
	owner, err := readLockOwner(l.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("release lifecycle lock: %w", err)
	}
	if owner.Token != l.token {
		return fmt.Errorf("release lifecycle lock: ownership changed")
	}
	if err := os.Remove(filepath.Join(l.path, "owner.json")); err != nil {
		return fmt.Errorf("release lifecycle lock owner: %w", err)
	}
	if err := os.Remove(l.path); err != nil {
		return fmt.Errorf("release lifecycle lock directory: %w", err)
	}
	return nil
}
