package event

type Event interface {
}

type DummyEvent struct {
  data int
  isCanceled bool
}

func (self *DummyEvent) IsCanceled() bool {
  return self.isCanceled 
}

func (self *DummyEvent) SetCanceled(is bool) {
  self.isCanceled = is
}
