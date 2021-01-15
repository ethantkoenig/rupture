package rupture

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/stretchr/testify/assert"
)

func TestShardedIndex(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "test")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	newIndex, err := NewShardedIndex(dir, bleve.NewIndexMapping(), 3)
	assert.NoError(t, err)

	addVehicles(t, newIndex.Index)
	_, err = newIndex.Document("7")
	assert.NoError(t, err)

	assert.NoError(t, newIndex.Close())

	openedIndex, err := OpenShardedIndex(dir)
	assert.NoError(t, err)

	result, err := openedIndex.Search(vehicleSearchRequest())
	assert.NoError(t, err)
	assert.EqualValues(t, result.Total, 7)
}

func TestShardedIndexFlushingBatch(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "test")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	index, err := NewShardedIndex(dir, bleve.NewIndexMapping(), 3)
	assert.NoError(t, err)

	batch := NewShardedFlushingBatch(index, 3)
	addVehicles(t, batch.Index)
	assert.NoError(t, batch.Flush())

	result, err := index.Search(vehicleSearchRequest())
	assert.NoError(t, err)
	assert.EqualValues(t, result.Total, 7)
}
