package event

import (
	"github.com/beakeyz/dadjoke-gen/pkg/web"
)

type UserCreateEvent struct {
	Ctx        *web.ReqContext
	isCanceled bool
}

func (self *UserCreateEvent) IsCanceled() bool {
	return self.isCanceled
}

func (self *UserCreateEvent) SetCanceled(is bool) {
	self.isCanceled = is
}
