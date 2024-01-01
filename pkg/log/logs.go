/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func init() {
	// set global logger with custom options
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			NoColor:    !isatty.IsTerminal(os.Stderr.Fd()),
		}),
	))
}

func Fatal(v ...any) {
	slog.Error(fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...any) {
	slog.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Printf(format string, v ...any) {
	slog.Info(fmt.Sprintf(format, v...))
}

func Println(v ...any) {
	slog.Info(fmt.Sprint(v...))
}
