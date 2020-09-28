package redis

import "context"

func (r *redisView) SLen(ctx context.Context, key string) (int64, error) {
	return r.cmd.SCard(r.expandKey(key)).Result()
}

func (r *redisView) SAdd(ctx context.Context, key string, values ...[]byte) (int64, error) {
	var in []interface{}
	for _, value := range values {
		in = append(in, value)
	}
	return r.cmd.SAdd(r.expandKey(key), in...).Result()
}

func (r *redisView) SRem(ctx context.Context, key string, values ...[]byte) (int64, error) {
	var in []interface{}
	for _, value := range values {
		in = append(in, value)
	}
	return r.cmd.SRem(r.expandKey(key), in...).Result()
}

func (r *redisView) SPop(ctx context.Context, key string) ([]byte, error) {
	result, err := r.cmd.SPop(r.expandKey(key)).Result()
	if nil != err {
		return nil, err
	}
	return []byte(result), nil
}

func (r *redisView) SPopN(ctx context.Context, key string, count int64) ([][]byte, error) {
	return wrapSliceStringToSliceBytes(func() ([]string, error) {
		return r.cmd.SPopN(r.expandKey(key), count).Result()
	})
}

func (r *redisView) SDiff(ctx context.Context, keys ...string) ([][]byte, error) {
	var inkeys []string
	for _, key := range keys {
		inkeys = append(inkeys, r.expandKey(key))
	}
	return wrapSliceStringToSliceBytes(func() ([]string, error) {
		return r.cmd.SDiff(inkeys...).Result()
	})
}

func (r *redisView) SDiffMerge(ctx context.Context, destination string, keys ...string) (int64, error) {
	var inkeys []string
	for _, key := range keys {
		inkeys = append(inkeys, r.expandKey(key))
	}
	return r.cmd.SDiffStore(destination, inkeys...).Result()
}

func (r *redisView) SInter(ctx context.Context, keys ...string) ([][]byte, error) {
	var inkeys []string
	for _, key := range keys {
		inkeys = append(inkeys, r.expandKey(key))
	}
	return wrapSliceStringToSliceBytes(func() ([]string, error) {
		return r.cmd.SInter(inkeys...).Result()
	})
}

func (r *redisView) SInterMerge(ctx context.Context, destination string, keys ...string) (int64, error) {
	var inkeys []string
	for _, key := range keys {
		inkeys = append(inkeys, r.expandKey(key))
	}
	return r.cmd.SInterStore(destination, keys...).Result()
}

func (r *redisView) SUnion(ctx context.Context, keys ...string) ([][]byte, error) {
	var inkeys []string
	for _, key := range keys {
		inkeys = append(inkeys, r.expandKey(key))
	}
	return wrapSliceStringToSliceBytes(func() ([]string, error) {
		return r.cmd.SUnion(inkeys...).Result()
	})
}

func (r *redisView) SUnionMerge(ctx context.Context, destination string, keys ...string) (int64, error) {
	var inkeys []string
	for _, key := range keys {
		inkeys = append(inkeys, r.expandKey(key))
	}
	return r.cmd.SUnionStore(destination, inkeys...).Result()
}
