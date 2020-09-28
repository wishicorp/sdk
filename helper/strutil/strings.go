package strutil

import (
	"bytes"
	"html/template"
	"math/rand"
	"strings"
	"time"
)

const (
	Digits      = "0123456789"
	Hexdigits   = "0123456789abcdefABCDEF"
	Letters     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	Lowercase   = "abcdefghijklmnopqrstuvwxyz"
	Uppercase   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Octdigits   = "01234567"
	Printable   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[]^_`{|}~"
	Punctuation = "!\"#$%&'()*+,-./:;<=>?@[]^_`{|}~"
	Whitespace  = "\t\n\x0b\x0c\r"
)

func IsBlank(s *string) bool {
	return nil == s || *s == ""
}

func NotBlank(s *string) bool {
	return nil != s && *s != ""
}

func RandAlphaString(n int) string {

	return randString(n, Letters)

}
func RandString(n int) string {
	return randString(n, Letters+Digits)
}

func RandNumeric(n int) string {
	return randString(n, Digits)
}

func randString(n int, letterBytes string) string {
	//const letterBytes = Letters + Digits

	var src = rand.NewSource(time.Now().UnixNano())

	const (
		letterIdxBits = 6
		letterIdxMask = 1<<letterIdxBits - 1
		letterIdxMax  = 63 / letterIdxBits
	)
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

type StringArray []string

func (a *StringArray) Get() interface{} { return []string(*a) }

func (a *StringArray) Set(s string) error {
	*a = append(*a, s)
	return nil
}

func (a *StringArray) String() string {
	return strings.Join(*a, ",")
}

func FormatDictionary(s string, dictionary map[string]interface{}) string {
	tmpl, err := template.New("").Parse(s)
	if err != nil {
		return ""
	}
	result := new(bytes.Buffer)

	err = tmpl.Execute(result, dictionary)
	if err != nil {
		return ""
	}
	return result.String()
}
