package block

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterSameIdTwiceShouldFail(t *testing.T) {
	var err error

	err = AddBlock(&Block{ID: "mockblock"})
	assert.NoError(t, err)

	err = AddBlock(&Block{ID: "mockblock"})
	assert.Error(t, err)
}

func TestRegisterBlockWithEmptyIDShouldFail(t *testing.T) {
	var err error

	err = AddBlock(&Block{})
	assert.Error(t, err)
}