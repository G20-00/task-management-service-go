package utils

import (
	"fmt"

	"github.com/G20-00/task-management-service-go/pkg/logger"
)

func RecoverPanic(layer, method string, err *error) {
	if rec := recover(); rec != nil {
		logger.GetLogger().WithFields(map[string]interface{}{
			"layer":  layer,
			"method": method,
			"error":  rec,
		}).Error(fmt.Sprintf("Panic in %s.%s", layer, method))
		*err = fmt.Errorf("panic occurred in %s", method)
	}
}
