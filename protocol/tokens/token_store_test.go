package tokens

import (
	"testing"

	"encoding/hex"

	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	store, err := FromFile("tokens.example.txt")
	assert.NoError(t, err)
	token, err := store.GetToken(123456)
	assert.NoError(t, err)
	assert.Equal(t, "ffffffffffffffffffffffffffffffff", hex.EncodeToString(token))
	token, err = store.GetToken(111222)
	assert.NoError(t, err)
	assert.Equal(t, "badcafefffffffffffffffffffffffff", hex.EncodeToString(token))
}
