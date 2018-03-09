package packet

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO (NW): Write more tests

func TestPkcs5Pad1(t *testing.T) {
	data := bytes.Repeat([]byte{0xff}, 16)
	result := pkcs5Pad(data, 16)
	assert.Len(t, result, 32)
	assert.Equal(t, append(data, bytes.Repeat([]byte{16}, 16)...), result)
}

func TestPkcs5Pad2(t *testing.T) {
	data := bytes.Repeat([]byte{0xff}, 15)
	result := pkcs5Pad(data, 16)
	assert.Len(t, result, 16)
	assert.Equal(t, append(data, 0x01), result)
}

func TestPkcs5Unpad1(t *testing.T) {
	data := bytes.Repeat([]byte{0xff}, 16)
	padded := append(data, bytes.Repeat([]byte{16}, 16)...)

	result, err := pkcs5Unpad(padded, 16)
	assert.NoError(t, err)
	assert.Len(t, result, 16)
	assert.Equal(t, data, result)
}

func TestPkcs5Unpad2(t *testing.T) {
	data := bytes.Repeat([]byte{0xff}, 15)
	padded := append(data, 0x01)

	result, err := pkcs5Unpad(padded, 16)
	assert.NoError(t, err)
	assert.Len(t, result, 15)
	assert.Equal(t, data, result)
}
