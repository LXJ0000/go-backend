package event

import "context"

const (
	MigratorEventBase               = "base"
	MigratorEventTarget             = "target"
	MigratorEventTypeNotEqual       = "not_equal"
	MigratorEventTypeBaseNotFound   = "base_not_found"
	MigratorEventTypeTargetNotFound = "target_not_found"
)

type MigratorEventProducer interface {
	Produce(ctx context.Context, event MigratorEvent) error
}

type MigratorEvent struct {
	ID        int64
	Direction string // 取值为 base 以源表为准 取值为 target 以目标表为准
	Type      string // 事件类型
}
