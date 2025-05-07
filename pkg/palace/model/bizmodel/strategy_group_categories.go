package bizmodel

import (
	"github.com/aide-family/moon/pkg/util/types"
)

const tableNameStrategyGroupCategories = "strategy_group_categories"

// StrategyGroupCategories 策略分组类型中间表
type StrategyGroupCategories struct {
	StrategyGroupID uint32 `gorm:"primaryKey;column:strategy_group_id;type:int unsigned;primaryKey" json:"strategy_group_id"`
	SysDictID       uint32 `gorm:"primaryKey;column:sys_dict_id;type:int unsigned;primaryKey" json:"sys_dict_id"`
}

// UnmarshalBinary redis存储实现
func (c *StrategyGroupCategories) UnmarshalBinary(data []byte) error {
	return types.Unmarshal(data, c)
}

// MarshalBinary redis存储实现
func (c *StrategyGroupCategories) MarshalBinary() (data []byte, err error) {
	return types.Marshal(c)
}

// TableName StrategyGroupCategories 's table name
func (*StrategyGroupCategories) TableName() string {
	return tableNameStrategyGroupCategories
}
