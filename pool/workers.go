package pool

import "sort"

var _ sort.Interface = (Workers)(nil)

//任务组
type Workers []*Worker

func (w Workers) Len() int {
	return len(w)
}

//按照任务chan size升序排列
func (w Workers) Less(i, j int) bool {
	return w[i].ChanSize() < w[j].ChanSize()
}

func (w Workers) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}
