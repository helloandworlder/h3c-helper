package biz

import (
	"context"
	"time"

	"github.com/aide-family/moon/cmd/server/houyi/internal/biz/bo"
	"github.com/aide-family/moon/cmd/server/houyi/internal/biz/repository"
)

// NewAlertBiz new AlertBiz
func NewAlertBiz(alertRepository repository.Alert, cacheRepo repository.CacheRepo) *AlertBiz {
	return &AlertBiz{
		alertRepository: alertRepository,
		cacheRepo:       cacheRepo,
	}
}

// AlertBiz .
type AlertBiz struct {
	alertRepository repository.Alert
	cacheRepo       repository.CacheRepo
}

// SaveAlarm 保存告警数据
func (a *AlertBiz) SaveAlarm(ctx context.Context, alarm *bo.Alarm) error {
	return a.alertRepository.SaveAlarm(ctx, alarm)
}

// PushAlarm 告警推送
func (a *AlertBiz) PushAlarm(ctx context.Context, alarm *bo.Alarm) (err error) {
	defer a.deleteAlarm(ctx, alarm, err)
	// 缓存告警数据推送标记
	alarm, err = a.cacheAlarm(ctx, alarm)
	if err != nil {
		return err
	}
	if len(alarm.Alerts) == 0 {
		return nil
	}
	return a.alertRepository.PushAlarm(ctx, alarm)
}

// cacheAlarm 缓存告警数据
func (a *AlertBiz) cacheAlarm(ctx context.Context, alarm *bo.Alarm) (*bo.Alarm, error) {
	alerts := alarm.Alerts
	pushAlerts := make([]*bo.Alert, 0, len(alerts))
	for _, alert := range alerts {
		ok, err := a.cacheRepo.Cacher().Client().SetNX(ctx, alert.PushedFlag(), "1", 60*time.Second).Result()
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		pushAlerts = append(pushAlerts, alert)
	}
	alarm.Alerts = pushAlerts
	return alarm, nil
}

// deleteAlarm 删除告警数据
func (a *AlertBiz) deleteAlarm(ctx context.Context, alarm *bo.Alarm, err error) {
	if err == nil {
		return
	}
	// 删除缓存
	for _, alert := range alarm.Alerts {
		_ = a.cacheRepo.Cacher().Client().Del(ctx, alert.PushedFlag()).Err()
	}
}
