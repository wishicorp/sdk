package pool

//in 输入数据
//out输出数据，err 错误消息
type Factory func() func(in interface{}) (out interface{}, err error)

var _ Subject = (*subject)(nil)

type Observer interface {
	Update(interface{}, error)
}

type Subject interface {
	Observer(o Observer)
}

type subject struct {
	data      interface{}
	observers []Observer
	result    interface{}
	err       error
}

func NewSubject(data interface{}) *subject {
	return &subject{
		data:      data,
		observers: make([]Observer, 0),
	}
}

func (s *subject) Observer(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *subject) notify() {

	for _, o := range s.observers {
		o.Update(s.result, s.err)
	}
}

func (s *subject) updateContext(result interface{}, err error) {
	s.result = result
	s.err = err
	s.notify()
}
