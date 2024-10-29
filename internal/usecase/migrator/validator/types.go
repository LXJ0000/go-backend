package validator

import (
	"context"
	"log/slog"

	"github.com/LXJ0000/go-backend/internal/usecase/migrator"
	"gorm.io/gorm"
)

type Validator[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB
}

// Validate 用户可以通过 context 来控制校验程序
// 全量校验 从数据库中一条一条查询出来，然后校验
func (v *Validator[T]) Validate(ctx context.Context) error {
	var offset int
	for {
		var b T
		err := v.base.Offset(offset).Order("id").First(&b).Error
		switch err {
		case nil:
			// 和 target 比较
			var t T
			err := v.target.Where("id = ?", b.ID()).First(&t).Error
			switch err {
			case nil: // 找到了 开始比较
			case gorm.ErrRecordNotFound: // target 中缺少数据
			default:
				// 做法1：认为大概率数据一致的 记录一下日志 continue
				slog.Error("校验数据 - 查询 target 出现错误", "error", err.Error())
				// 做法2：出于保险	考虑，报 数据不一致 尝试去修复
			}
		case gorm.ErrRecordNotFound: // 全量校验结束
			return nil
		default:
			slog.Error("校验数据 - 查询 base 出现错误", "error", err.Error())
		}
		offset++
	}
}
