// +build !js

/*
 * Copyright 2017 Workiva
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package frugal provides the library APIs used by the Frugal code generator.
package frugal

import (
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/Sirupsen/logrus"
	"github.com/nats-io/nuid"
)

var (
	packageLogger = logrus.StandardLogger()
	loggerMu      sync.RWMutex
)

// SetLogger sets the Logger used by Frugal.
func SetLogger(logger *logrus.Logger) {
	loggerMu.Lock()
	packageLogger = logger
	loggerMu.Unlock()
}

// logger returns the global Logger. Do not mutate the pointer returned, use
// SetLogger instead.
func logger() *logrus.Logger {
	loggerMu.RLock()
	logger := packageLogger
	loggerMu.RUnlock()
	return logger
}

// generateCorrelationID returns a random string id. It's assigned to a var for
// testability purposes.
var generateCorrelationID = func() string {
	return nuid.Next()
}

var getNextOpID = func() string {
	return strconv.FormatUint(atomic.AddUint64(&nextOpID, 1), 10)
}
