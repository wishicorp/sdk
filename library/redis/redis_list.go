package redis

import "context"

func (r *redisView) LRem(ctx context.Context, key string, count int64, value []byte) (int64, error) {
	return r.cmd.LRem(r.expandKey(key), count, value).Result()
}

func (r *redisView) LIndex(ctx context.Context, key string, index int64) ([]byte, error) {
	result, err := r.cmd.LIndex(r.expandKey(key), index).Result()
	if nil != err {
		return nil, err
	}
	return []byte(result), nil
}

func (r *redisView) LTrim(ctx context.Context, key string, start, stop int64) error {
	return wrapResult(func() (interface{}, error) {
		return r.cmd.LTrim(r.expandKey(key), start, stop).Result()
	})
}

func (r *redisView) LSet(ctx context.Context, key string, index int64, value []byte) error {
	return wrapResult(func() (interface{}, error) {
		return r.cmd.LSet(r.expandKey(key), index, value).Result()
	})
}

func (r *redisView) LPush(ctx context.Context, key string, values ...[]byte) (int64, error) {
	var vals []interface{}
	for _, value := range values {
		vals = append(vals, value)
	}
	return r.cmd.LPush(r.expandKey(key), vals...).Result()
}

func (r *redisView) LAppend(ctx context.Context, key string, values ...[]byte) (int64, error) {
	var vals []interface{}
	for _, value := range values {
		vals = append(vals, value)
	}
	return r.cmd.RPush(r.expandKey(key), vals...).Result()
}

func (r *redisView) LPop(ctx context.Context, key string) ([]byte, error) {
	result, err := r.cmd.LPop(r.expandKey(key)).Result()
	if nil != err {
		return nil, err
	}
	return []byte(result), nil
}

func (r *redisView) LRPop(ctx context.Context, key string) ([]byte, error) {
	result, err := r.cmd.RPop(r.expandKey(key)).Result()
	if nil != err {
		return nil, err
	}
	return []byte(result), nil
}

func (r *redisView) LRange(ctx context.Context, key string, start, stop int64) ([][]byte, error) {
	return wrapSliceStringToSliceBytes(func() ([]string, error) {
		return r.cmd.LRange(r.expandKey(key), start, stop).Result()
	})
}

func (r *redisView) LLen(ctx context.Context, key string) (int64, error) {
	return r.cmd.LLen(r.expandKey(key)).Result()
}

func (r *redisView) LInsert(ctx context.Context, key string, op InsertOP, pivot, value []byte) (int64, error) {
	return r.cmd.LInsert(r.expandKey(key), string(op), pivot, value).Result()
}
