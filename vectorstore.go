package nodechain

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
)

type Document struct {
	ID        string
	Text      string
	Metadata  map[string]any
	Embedding []float32
}

type VectorStore interface {
	Add(ctx context.Context, docs []Document) error
	Search(ctx context.Context, query []float32, k int) ([]Document, error)
}

type InMemoryVectorStore struct {
	mu   sync.RWMutex
	docs []Document
}

func NewInMemoryVectorStore() *InMemoryVectorStore {
	return &InMemoryVectorStore{
		docs: make([]Document, 0),
	}
}

func (s *InMemoryVectorStore) Add(ctx context.Context, docs []Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.docs = append(s.docs, docs...)
	return nil
}

func (s *InMemoryVectorStore) Search(ctx context.Context, query []float32, k int) ([]Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(query) == 0 {
		return nil, fmt.Errorf("Search: empty query embedding")
	}

	type scored struct {
		doc   Document
		score float64
	}

	var results []scored
	for _, d := range s.docs {
		if len(d.Embedding) != len(query) {
			continue
		}
		score := cosineSimilarity(query, d.Embedding)
		results = append(results, scored{doc: d, score: score})
	}

	if len(results) == 0 {
		return nil, nil
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if k > len(results) {
		k = len(results)
	}

	out := make([]Document, k)
	for i := 0; i < k; i++ {
		out[i] = results[i].doc
	}
	return out, nil
}

func cosineSimilarity(a, b []float32) float64 {
	var dot float64
	var na, nb float64

	for i := 0; i < len(a); i++ {
		av := float64(a[i])
		bv := float64(b[i])
		dot += av * bv
		na += av * av
		nb += bv * bv
	}

	if na == 0 || nb == 0 {
		return 0
	}

	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}
