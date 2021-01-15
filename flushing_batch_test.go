package rupture

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/stretchr/testify/assert"
)

func TestFlushingBatch(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "test")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	index, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	assert.NoError(t, err)

	batch := NewFlushingBatch(index, 4)
	addVehicles(t, batch.Index)
	assert.NoError(t, batch.Flush())

	result, err := index.Search(vehicleSearchRequest())
	assert.NoError(t, err)
	assert.EqualValues(t, 7, result.Total)

	assert.NoError(t, batch.Delete("10"))
	assert.NoError(t, batch.Delete("11"))
	assert.NoError(t, batch.Delete("13"))
	assert.NoError(t, batch.Flush())

	result, err = index.Search(vehicleSearchRequest())
	assert.NoError(t, err)
	assert.EqualValues(t, 5, result.Total)
}
