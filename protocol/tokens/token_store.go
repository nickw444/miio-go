package tokens

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TokenStore interface {
	LoadFile(inputPath string) error
	WriteFile(outputPath string) error
	GetToken(deviceId uint32) ([]byte, error)
	AddDevice(deviceId uint32, token []byte) error
	RemoveDevice(deviceId uint32)
}

type tokenStore struct {
	tokens map[uint32][]byte
}

func New() TokenStore {
	return &tokenStore{
		tokens: make(map[uint32][]byte),
	}
}

func FromFile(filePath string) (TokenStore, error) {
	store := New()
	err := store.LoadFile(filePath)
	return store, err
}

func (t *tokenStore) LoadFile(inputPath string) error {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		// File doesn't exist, so don't load anything.
		return nil
	}

	f, err := os.Open(inputPath)
	defer f.Close()

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		splitLine := strings.Split(line, "=")
		if len(splitLine) != 2 {
			return fmt.Errorf("Malformed line: %s", line)
		}

		deviceId, err := strconv.ParseUint(splitLine[0], 10, 32)
		if err != nil {
			return fmt.Errorf("Malformed line: %s", line)
		}

		token, err := hex.DecodeString(splitLine[1])
		if err != nil {
			return fmt.Errorf("Malformed line: %s", line)
		}

		t.tokens[uint32(deviceId)] = token
	}

	return nil
}

func (t *tokenStore) WriteFile(outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	for deviceId, tokenBytes := range t.tokens {
		tokenStr := hex.EncodeToString(tokenBytes)
		_, err := f.WriteString(fmt.Sprintf("%d=%s\n", deviceId, tokenStr))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *tokenStore) GetToken(deviceId uint32) ([]byte, error) {
	if val, ok := t.tokens[deviceId]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Device ID %d does not exist in token store", deviceId)
}

func (t *tokenStore) AddDevice(deviceId uint32, token []byte) error {
	t.tokens[deviceId] = token
	return nil
}

func (t *tokenStore) RemoveDevice(deviceId uint32) {
	delete(t.tokens, deviceId)
}
