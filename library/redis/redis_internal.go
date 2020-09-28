package redis

import (
	"gopkg.in/redis.v5"
	"strings"
)

// expandKey is used to expand to the full key path with the prefix
func (r *redisView) expandKey(suffix string) string {
	return r.prefix + RedisKeySep + suffix
}

// truncateKey is used to remove the prefix of the key
func (r *redisView) truncateKey(full string) string {
	return strings.Join(strings.Split(full, RedisKeySep)[1:], RedisKeySep)
}
func wrapResult(call func() (interface{}, error)) error {
	result, err := call()
	if nil != err {
		return err
	}
	switch result.(type) {
	case bool:
		if result != true {
			return ErrorResultNotTrue
		}
	case string:
		if result != "OK" {
			return ErrorResultNotOK
		}
	}

	return nil
}
func wrapSliceStringToSliceBytes(call func() ([]string, error)) ([][]byte, error) {
	results, err := call()
	if nil != err {
		return nil, err
	}
	var out [][]byte
	for _, result := range results {
		out = append(out, []byte(result))
	}
	return out, nil
}
func toRangeZMembers(err error, members []string, zSlice []redis.Z) ([]*ZMember, error) {
	if nil != err {
		return nil, err
	}
	var zMembers []*ZMember
	if len(members) > 0 {
		for _, m := range members {
			member := ZMember{
				Member: []byte(m),
			}
			zMembers = append(zMembers, &member)
		}
	}
	if len(members) > 0 {
		for _, m := range zSlice {
			member := ZMember{
				Score:  m.Score,
				Member: m.Member.([]byte),
			}
			zMembers = append(zMembers, &member)
		}
	}
	return zMembers, nil
}
