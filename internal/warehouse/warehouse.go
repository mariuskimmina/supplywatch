package warehouse

import "github.com/mariuskimmina/supplywatch/pkg/logger"

type warehouse struct {

}

func NewWarehouse() *warehouse {
    return &warehouse{}
}

func (w *warehouse) Start() {
    logger.Debug("Starting warehouse")
}
