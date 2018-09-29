package util_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"util"
)

func TestEncodeAndDecode(t *testing.T) {
	s := "I love PingCAP"

	encoded := util.Encode([]byte(s))
	decoded, err := util.Decode(encoded)

	assert.Nil(t, err)
	assert.Equal(t, string(decoded), s)

}

func TestEncodeAndDecodeWithSalt(t *testing.T) {
	s := "I love PingCAP"

	encoded := util.EncodeWithRandomSalt([]byte(s))
	decoded, err := util.DecodeFromRandomSalt(encoded)

	assert.Nil(t, err)
	assert.Equal(t, string(decoded), s)

}
