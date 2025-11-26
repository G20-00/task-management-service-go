// Package utils provides utility functions for common operations across the application.
package utils

import (
	"fmt"

	"github.com/G20-00/task-management-service-go/pkg/logger"
)

// RecoverPanic recovers from panics in the specified layer and method, logging the error and setting the error pointer.
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
