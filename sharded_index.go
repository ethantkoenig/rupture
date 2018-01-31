package rupture

import (
	"fmt"
	"hash/fnv"
	"path/filepath"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
)

// a type alias for bleve.Index, so that the anonymous field of
// ShardedIndex does not conflict with the Index(..) method.
type bleveIndex bleve.Index

type ShardedIndex struct {
	bleveIndex
	indices []bleve.Index
}

func hash(id string, n int) uint64 {
	fnvHash := fnv.New64()
	fnvHash.Write([]byte(id))
	return fnvHash.Sum64() % uint64(n)
}

func childIndexerPath(rootPath string, i int) string {
	return filepath.Join(rootPath, strconv.Itoa(i))
}

// NewShardedIndex createas a sharded index at the specified path, with the
// specified mapping and number of shards.
func NewShardedIndex(path string, mapping mapping.IndexMapping, numShards int) (*ShardedIndex, error) {
	if numShards <= 0 {
		return nil, fmt.Errorf("Invalid number of shards: %d", numShards)
	}
	err := writeJson(shardedIndexMetadataPath(path), &shardedIndexMetadata{NumShards: numShards})
	if err != nil {
		return nil, err
	}

	s := &ShardedIndex{
		indices: make([]bleve.Index, numShards),
	}
	for i := 0; i < numShards; i++ {
		s.indices[i], err = bleve.New(childIndexerPath(path, i), mapping)
		if err != nil {
			return nil, err
		}
	}
	s.bleveIndex = bleve.NewIndexAlias(s.indices...)
	return s, nil
}

// OpenShardedIndex opens a sharded index at the specified path.
func OpenShardedIndex(path string) (*ShardedIndex, error) {
	var meta shardedIndexMetadata
	var err error
	if err = readJson(shardedIndexMetadataPath(path), &meta); err != nil {
		return nil, err
	}

	s := &ShardedIndex{
		indices: make([]bleve.Index, meta.NumShards),
	}
	for i := 0; i < meta.NumShards; i++ {
		s.indices[i], err = bleve.Open(childIndexerPath(path, i))
		if err != nil {
			return nil, err
		}
	}
	s.bleveIndex = bleve.NewIndexAlias(s.indices...)
	return s, nil
}

func (s *ShardedIndex) Index(id string, data interface{}) error {
	return s.indices[hash(id, len(s.indices))].Index(id, data)
}

func (s *ShardedIndex) Delete(id string) error {
	return s.indices[hash(id, len(s.indices))].Delete(id)
}

func (s *ShardedIndex) Document(id string) (*document.Document, error) {
	return s.indices[hash(id, len(s.indices))].Document(id)
}

func (s *ShardedIndex) Close() error {
	if err := s.bleveIndex.Close(); err != nil {
		return err
	}
	for _, index := range s.indices {
		if err := index.Close(); err != nil {
			return err
		}
	}
	return nil
}

type shardedIndexFlushingBatch struct {
	batches []*singleIndexFlushingBatch
}

func NewShardedFlushingBatch(index *ShardedIndex, maxBatchSize int) FlushingBatch {
	b := &shardedIndexFlushingBatch{
		batches: make([]*singleIndexFlushingBatch, len(index.indices)),
	}
	for i, index := range index.indices {
		b.batches[i] = newFlushingBatch(index, maxBatchSize)
	}
	return b
}

func (b *shardedIndexFlushingBatch) Index(id string, data interface{}) error {
	return b.batches[hash(id, len(b.batches))].Index(id, data)
}

func (b *shardedIndexFlushingBatch) Delete(id string) error {
	return b.batches[hash(id, len(b.batches))].Delete(id)
}

func (b *shardedIndexFlushingBatch) Flush() error {
	for _, batch := range b.batches {
		if err := batch.Flush(); err != nil {
			return err
		}
	}
	return nil
}
