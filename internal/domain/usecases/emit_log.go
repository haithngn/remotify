package usecases

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
)

type LogData struct {
	LoggedAt  string      `json:"logged_at"`
	Direction string      `json:"direction"`
	Message   interface{} `json:"message"`
}

type LogDirection string

const (
	LogDirectionSent     LogDirection = ">>>"
	LogDirectionReceived LogDirection = "<<<"
)

func EmitLog(message interface{}, direction LogDirection, ctx context.Context) {
	directionString := string(direction)
	runtime.EventsEmit(ctx, "onLog", LogData{
		LoggedAt:  time.Now().Format(time.RFC3339),
		Direction: directionString,
		Message:   message,
	})
}
