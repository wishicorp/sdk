package httputil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// http客户端请求分页结构
type Pageable struct {
	page  int
	size  int
	index int
	sort  []Sort
}

// http客户端请求排序结构
type Sort struct {
	Key       string
	Direction string
}

func NewPageable(page, size int) *Pageable {

	return &Pageable{page: page, size: size}
}

func (p *Pageable) Skip() int {
	s := p.size * p.page
	if s < 0 {
		s = 0
	}
	return s
}

func (p *Pageable) Page() int {
	return p.page
}
func (p *Pageable) Size() int {
	return p.size
}
func (p *Pageable) Sort() string {
	if nil == p.sort || 0 == len(p.sort) {
		return ""
	}

	sorts := make([]string, len(p.sort))
	for idx := range p.sort {
		sorts[idx] = fmt.Sprintf("%s %s", p.sort[idx].Key, p.sort[idx].Direction)
	}
	return strings.Join(sorts, ",")
}
func ParsePageable(c *gin.Context) *Pageable {

	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	sorts := ParseSort(c)
	if size < 1 {
		size = 1
	}
	if page == 0 {
		page = 50
	}
	return &Pageable{page: page, size: size, sort: sorts}
}

func ParseSort(c *gin.Context) []Sort {
	sort := c.QueryArray("sort")
	sorts := make([]Sort, len(sort))
	for k, v := range sort {
		vs := strings.Split(v, ",")
		key := vs[0]
		direction := "ASC"
		if len(vs) == 2 {
			direction = vs[1]
		}
		sorts[k] = Sort{Key: key, Direction: strings.ToUpper(direction)}
	}
	return sorts
}
