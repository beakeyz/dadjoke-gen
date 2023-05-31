package subapps

import "context"

type SubService interface {
	Run(ctx context.Context) error
}

type SubServiceStruct struct{}
