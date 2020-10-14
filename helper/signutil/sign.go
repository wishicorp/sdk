package signutil

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/wishicorp/sdk/helper/jsonutil"
	"sort"
	"strconv"
	"strings"
)

var ErrInvalidSign = errors.New("invalid sign")

type Params map[string]string
type SignType string

const (
	SHA1       SignType = "SHA1"
	SHA256     SignType = "SHA256"
	HMACSHA1   SignType = "HMAC-SHA1"
	HMACSHA256 SignType = "HMAC-SHA256"
	MD5        SignType = "MD5"
	SignKey             = "sign"
)

func KVSignVerify(signType SignType, key string, in interface{}, sign string) error {
	nSign, err := KVSign(signType, key, in)
	if nil != err {
		return err
	}
	if nSign != strings.ToLower(sign) {
		return ErrInvalidSign
	}
	return nil
}

func ValueSignVerify(signType SignType, key string, in interface{}, sign string) error {
	nSign, err := ValueSign(signType, key, in)
	if nil != err {
		return err
	}
	if nSign != strings.ToLower(sign) {
		return ErrInvalidSign
	}
	return nil
}

func ValueSign(signType SignType, key string, in interface{}) (string, error) {
	var data map[string]string
	if err := jsonutil.Swap(in, &data); err != nil {
		return "", err
	}
	return Sign(signType, key, data, false), nil
}

func KVSign(signType SignType, key string, in interface{}) (string, error) {
	var data map[string]string
	if err := jsonutil.Swap(in, &data); err != nil {
		return "", err
	}
	return Sign(signType, key, data, true), nil
}

func Sign(signType SignType, key string, params Params, withKey bool) string {
	// 创建切片
	var keys = make([]string, 0, len(params))
	// 遍历签名参数
	for k := range params {
		if k != SignKey { // 排除sign字段
			keys = append(keys, k)
		}
	}
	// 由于切片的元素顺序是不固定，所以这里强制给切片元素加个顺序
	sort.Strings(keys)
	//创建字符缓冲
	var buf bytes.Buffer
	for _, k := range keys {
		if len(params.GetString(k)) > 0 {
			if withKey {
				buf.WriteString(k)
				buf.WriteString(`=`)
			}
			buf.WriteString(params.GetString(k))
			if withKey {
				buf.WriteString(`&`)
			}
		}
	}
	// 加入apiKey作加密密钥
	if withKey {
		buf.WriteString(`key=`)
	}
	buf.WriteString(key)
	fmt.Println(buf.String())
	return sign(signType, key, buf)

}

func sign(signType SignType, key string, buf bytes.Buffer) string {
	var (
		dataMd5 [16]byte
		dataSha []byte
		str     string
	)
	switch signType {
	case MD5:
		dataMd5 = md5.Sum(buf.Bytes())
		str = hex.EncodeToString(dataMd5[:]) //需转换成切片
	case HMACSHA256:
		h := hmac.New(sha256.New, []byte(key))
		h.Write(buf.Bytes())
		dataSha = h.Sum(nil)
		str = hex.EncodeToString(dataSha[:])
	case HMACSHA1:
		h := hmac.New(sha1.New, []byte(key))
		h.Write(buf.Bytes())
		dataSha = h.Sum(nil)
		str = hex.EncodeToString(dataSha[:])
	case SHA1:
		h := sha1.New()
		h.Write(buf.Bytes())
		dataSha = h.Sum(nil)
		str = hex.EncodeToString(dataSha[:])
	case SHA256:
		h := sha256.New()
		h.Write(buf.Bytes())
		dataSha = h.Sum(nil)
		str = hex.EncodeToString(dataSha[:])
	}

	return str
}

// map本来已经是引用类型了，所以不需要 *Params
func (p Params) SetString(k, s string) Params {
	p[k] = s
	return p
}

func (p Params) GetString(k string) string {
	s, _ := p[k]
	return s
}

func (p Params) SetInt64(k string, i int64) Params {
	p[k] = strconv.FormatInt(i, 10)
	return p
}

func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}
func (p Params) SetInt(k string, i int) Params {
	p[k] = fmt.Sprintf("%d", i)
	return p
}
func (p Params) GetInt(k string) int {
	i, _ := strconv.Atoi(p.GetString(k))
	return i
}

// 判断key是否存在
func (p Params) ContainsKey(key string) bool {
	_, ok := p[key]
	return ok
}
