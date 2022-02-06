package memory

import "task/pkg/model"

type Memory interface {
	Set(string, model.Order)
	Get(string) (model.Order, bool)
}
