package middleware

import (
	"context"

	"github.com/sirupsen/logrus"

	logger "github.com/alighm/sample-service/log"
)

func log(ctx context.Context) logrus.FieldLogger {
	return logger.Logger(ctx).WithField("pkg", "middleware")
}
