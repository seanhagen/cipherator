package translator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslator_Basics(t *testing.T) {
	var tr *Translator
	var err error

	tr, err = NewPigLatin()
	assert.NotNil(t, tr)
	assert.NoError(t, err)
}
