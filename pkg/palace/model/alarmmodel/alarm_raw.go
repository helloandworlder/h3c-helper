package alarmmodel

import (
	"encoding/json"

	"github.com/aide-family/moon/pkg/palace/model"

	"github.com/go-kratos/kratos/v2/log"
)

const tableNameAlarmRaws = "alarm_raws"

// AlarmRaw 告警原始信息
type AlarmRaw struct {
	model.EasyModel
	// 接收对象
	Receiver string `gorm:"column:receiver;type:text;not null;comment:接收对象"`
	// 原始信息json
	RawInfo string `gorm:"column:raw_info;type:text;not null;comment:原始信息json"`
	// 指纹
	Fingerprint string `gorm:"column:fingerprint;type:varchar(255);not null;comment:fingerprint;uniqueIndex"`
}

// String json string
func (a *AlarmRaw) String() string {
	bs, err := json.Marshal(a)
	if err != nil {
		log.Warnw("method", "AlarmRaw.string", "err", err)
		return ""
	}
	return string(bs)
}

// UnmarshalBinary redis存储实现
func (a *AlarmRaw) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

// MarshalBinary redis存储实现
func (a *AlarmRaw) MarshalBinary() (data []byte, err error) {
	return json.Marshal(a)
}

// TableName AlarmRaw's table name
func (a *AlarmRaw) TableName() string {
	return tableNameAlarmRaws
}
