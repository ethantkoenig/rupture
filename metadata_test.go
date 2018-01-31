package rupture

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "test")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	meta, err := ReadIndexMetadata(dir)
	assert.NoError(t, err)
	assert.Equal(t, &IndexMetadata{}, meta)

	meta.Version = 24
	assert.NoError(t, WriteIndexMetadata(dir, meta))

	meta, err = ReadIndexMetadata(dir)
	assert.NoError(t, err)
	assert.EqualValues(t, 24, meta.Version)
}
