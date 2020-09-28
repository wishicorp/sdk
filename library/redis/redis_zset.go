package redis

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/redis.v5"
)

func (r *redisView) ZLen(ctx context.Context, key string) (int64, error) {
	return r.cmd.ZCard(r.expandKey(key)).Result()
}

func (r *redisView) ZCount(ctx context.Context, key string, min, max float64) (int64, error) {
	return r.cmd.ZCount(r.expandKey(key),
		fmt.Sprintf("%f", min), fmt.Sprintf("%f", max)).Result()
}

func (r *redisView) ZLexCount(ctx context.Context, key, min, max string) (int64, error) {
	return 0, errors.New("not implemented")
}

func (r *redisView) ZAdd(ctx context.Context, key string, members ...*ZMember) (int64, error) {
	var zS []redis.Z
	for _, member := range members {
		z := redis.Z{
			Score:  member.Score,
			Member: member.Member,
		}
		zS = append(zS, z)
	}
	return r.cmd.ZAdd(r.expandKey(key), zS...).Result()
}

func (r *redisView) ZRem(ctx context.Context, key string, members ...*ZMember) (int64, error) {
	var zS []interface{}
	for _, member := range members {
		zS = append(zS, member.Member)
	}
	return r.cmd.ZRem(r.expandKey(key), zS...).Result()
}

func (r *redisView) ZRemRangeByLex(ctx context.Context, key, min, max string) (int64, error) {
	return r.cmd.ZRemRangeByLex(r.expandKey(key), min, max).Result()
}

func (r *redisView) ZRemRangeByScore(ctx context.Context, key string, min, max float64) (int64, error) {
	return r.cmd.ZRemRangeByScore(r.expandKey(key),
		fmt.Sprintf("%f", min), fmt.Sprintf("%f", max)).Result()
}

func (r *redisView) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return r.cmd.ZRemRangeByRank(r.expandKey(key), start, stop).Result()
}

func (r *redisView) ZRange(ctx context.Context, key string, start, stop int64, reverse, withScores bool) ([]*ZMember, error) {
	var err error
	var members []string
	var zSlice []redis.Z
	if !reverse && !withScores {
		members, err = r.cmd.ZRange(r.expandKey(key), start, stop).Result()
	}
	if reverse && withScores {
		zSlice, err = r.cmd.ZRevRangeWithScores(r.expandKey(key), start, stop).Result()
	}
	if reverse {
		members, err = r.cmd.ZRevRange(r.expandKey(key), start, stop).Result()
	}
	if withScores {
		zSlice, err = r.cmd.ZRangeWithScores(r.expandKey(key), start, stop).Result()
	}
	return toRangeZMembers(err, members, zSlice)
}

func (r *redisView) ZRangeByScore(ctx context.Context, key string, rangeBy ZRangeBy, reverse, withScores bool) ([]*ZMember, error) {
	var err error
	var members []string
	var zSlice []redis.Z
	if !reverse && !withScores {
		members, err = r.cmd.ZRangeByScore(r.expandKey(key), rangeBy.ToRedisRangeBy()).Result()
	}
	if reverse && withScores {
		zSlice, err = r.cmd.ZRevRangeByScoreWithScores(r.expandKey(key), rangeBy.ToRedisRangeBy()).Result()
	}
	if reverse {
		members, err = r.cmd.ZRevRangeByScore(r.expandKey(key), rangeBy.ToRedisRangeBy()).Result()
	}
	if withScores {
		zSlice, err = r.cmd.ZRangeByScoreWithScores(r.expandKey(key), rangeBy.ToRedisRangeBy()).Result()
	}
	return toRangeZMembers(err, members, zSlice)
}

func (r *redisView) ZRangeByLex(ctx context.Context, key string, rangeBy ZRangeBy, reverse bool) ([]*ZMember, error) {
	var err error
	var members []string
	if reverse {
		members, err = r.cmd.ZRevRangeByLex(key, rangeBy.ToRedisRangeBy()).Result()
	}
	members, err = r.cmd.ZRangeByLex(key, rangeBy.ToRedisRangeBy()).Result()
	return toRangeZMembers(err, members, []redis.Z{})
}

func (r *redisView) ZRank(ctx context.Context, key string, member []byte, reverse bool) (int64, error) {
	if reverse {
		return r.cmd.ZRevRank(r.expandKey(key), string(member)).Result()
	}
	return r.cmd.ZRank(r.expandKey(key), string(member)).Result()
}

func (r *redisView) ZIncr(ctx context.Context, key string, member *ZMember) (float64, error) {
	return r.cmd.ZIncr(r.expandKey(key), redis.Z{Score: member.Score, Member: member.Member}).Result()
}

func (r *redisView) ZIncrNX(ctx context.Context, key string, member *ZMember) (float64, error) {
	return r.cmd.ZIncrNX(r.expandKey(key), redis.Z{Score: member.Score, Member: member.Member}).Result()
}

func (r *redisView) ZInterMerge(ctx context.Context, destination string, merge *ZMerge, keys ...string) (int64, error) {
	var inKeys []string
	for _, key := range keys {
		inKeys = append(inKeys, r.expandKey(key))
	}
	return r.cmd.ZInterStore(destination, merge.ToZStore(), inKeys...).Result()
}

func (r *redisView) ZUnionMerge(ctx context.Context, destination string, merge *ZMerge, keys ...string) (int64, error) {
	var inKeys []string
	for _, key := range keys {
		inKeys = append(inKeys, r.expandKey(key))
	}
	return r.cmd.ZUnionStore(destination, merge.ToZStore(), inKeys...).Result()
}
