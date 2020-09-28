package logical

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	log "github.com/hashicorp/go-hclog"

	radix "github.com/armon/go-radix"
)

// InMemBackend is an in-memory only physical backend. It is useful
// for testing and development situations where the data is not
// expected to be durable.
type InMemBackend struct {
	sync.RWMutex
	root         *radix.Tree
	logger       log.Logger
	failGet      *uint32
	failPut      *uint32
	failDelete   *uint32
	failList     *uint32
	logOps       bool
	maxValueSize int
}

// NewInmem constructs a new in-memory backend
func NewInmem(conf map[string]string, logger log.Logger) (Storage, error) {
	maxValueSize := 0
	maxValueSizeStr, ok := conf["max_value_size"]
	if ok {
		var err error
		maxValueSize, err = strconv.Atoi(maxValueSizeStr)
		if err != nil {
			return nil, err
		}
	}

	return &InMemBackend{
		root:         radix.New(),
		logger:       logger,
		failGet:      new(uint32),
		failPut:      new(uint32),
		failDelete:   new(uint32),
		failList:     new(uint32),
		logOps:       logger.IsTrace(),
		maxValueSize: maxValueSize,
	}, nil
}

// Put is used to insert or update an entry
func (i *InMemBackend) Put(ctx context.Context, entry *StorageEntry) error {
	i.Lock()
	defer i.Unlock()

	return i.PutInternal(ctx, entry)
}

func (i *InMemBackend) PutInternal(ctx context.Context, entry *StorageEntry) error {
	if i.logOps {
		i.logger.Trace("put", "key", entry.Key)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if i.maxValueSize > 0 && len(entry.Value) > i.maxValueSize {
		return fmt.Errorf("%s", "put failed due to value being too large")
	}

	i.root.Insert(entry.Key, entry.Value)
	return nil
}

func (i *InMemBackend) FailPut(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(i.failPut, val)
}

// Backend is used to fetch an entry
func (i *InMemBackend) Get(ctx context.Context, key string) (*StorageEntry, error) {

	i.RLock()
	defer i.RUnlock()

	return i.GetInternal(ctx, key)
}

func (i *InMemBackend) GetInternal(ctx context.Context, key string) (*StorageEntry, error) {
	if i.logOps {
		i.logger.Trace("get", "key", key)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if raw, ok := i.root.Get(key); ok {
		return &StorageEntry{
			Key:   key,
			Value: raw.([]byte),
		}, nil
	}
	return nil, nil
}

func (i *InMemBackend) FailGet(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(i.failGet, val)
}

// Delete is used to permanently delete an entry
func (i *InMemBackend) Delete(ctx context.Context, key string) error {

	i.Lock()
	defer i.Unlock()

	return i.DeleteInternal(ctx, key)
}

func (i *InMemBackend) DeleteInternal(ctx context.Context, key string) error {
	if i.logOps {
		i.logger.Trace("delete", "key", key)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	i.root.Delete(key)
	return nil
}

func (i *InMemBackend) FailDelete(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(i.failDelete, val)
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (i *InMemBackend) List(ctx context.Context, prefix string) ([]*StorageEntry, error) {

	i.RLock()
	defer i.RUnlock()

	return i.ListInternal(ctx, prefix)
}

func (i *InMemBackend) ListInternal(ctx context.Context, prefix string) ([]*StorageEntry, error) {
	if i.logOps {
		i.logger.Trace("list", "prefix", prefix)
	}

	var out []*StorageEntry
	seen := make(map[string]interface{})
	walkFn := func(s string, v interface{}) bool {
		trimmed := strings.TrimPrefix(s, prefix)
		sep := strings.Index(trimmed, "/")
		if sep == -1 {
			out = append(out, &StorageEntry{
				Key:   s,
				Value: v.([]byte),
			})
		} else {
			trimmed = trimmed[:sep+1]
			if _, ok := seen[trimmed]; !ok {
				out = append(out, &StorageEntry{
					Key:   s,
					Value: v.([]byte),
				})
				seen[trimmed] = struct{}{}
			}
		}
		return false
	}
	i.root.WalkPrefix(prefix, walkFn)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return out, nil
}

func (i *InMemBackend) FailList(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(i.failList, val)
}
