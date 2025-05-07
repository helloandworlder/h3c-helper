package realtime

import (
	"context"

	realtimeapi "github.com/aide-family/moon/api/admin/realtime"
	"github.com/aide-family/moon/cmd/server/palace/internal/biz"
	"github.com/aide-family/moon/cmd/server/palace/internal/service/builder"
)

// AlarmService 实时告警数据服务
type AlarmService struct {
	realtimeapi.UnimplementedAlarmServer

	alarmBiz *biz.AlarmBiz
}

// NewAlarmService 实时告警数据服务
func NewAlarmService(alarmBiz *biz.AlarmBiz) *AlarmService {
	return &AlarmService{
		alarmBiz: alarmBiz,
	}
}

// GetAlarm 获取实时告警数据
func (s *AlarmService) GetAlarm(ctx context.Context, req *realtimeapi.GetAlarmRequest) (*realtimeapi.GetAlarmReply, error) {
	params := builder.NewParamsBuild(ctx).RealtimeAlarmModuleBuilder().WithGetAlarmRequest(req).ToBo()
	realtimeAlarmDetail, err := s.alarmBiz.GetRealTimeAlarm(ctx, params)
	if err != nil {
		return nil, err
	}
	return &realtimeapi.GetAlarmReply{
		Detail: builder.NewParamsBuild(ctx).RealtimeAlarmModuleBuilder().DoRealtimeAlarmBuilder().ToAPI(realtimeAlarmDetail),
	}, nil
}

// ListAlarm 获取实时告警数据列表
func (s *AlarmService) ListAlarm(ctx context.Context, req *realtimeapi.ListAlarmRequest) (*realtimeapi.ListAlarmReply, error) {
	params := builder.NewParamsBuild(ctx).RealtimeAlarmModuleBuilder().WithListAlarmRequest(req).ToBo()
	list, err := s.alarmBiz.ListRealTimeAlarms(ctx, params)
	if err != nil {
		return nil, err
	}
	return &realtimeapi.ListAlarmReply{
		List:       builder.NewParamsBuild(ctx).RealtimeAlarmModuleBuilder().DoRealtimeAlarmBuilder().ToAPIs(list),
		Pagination: builder.NewParamsBuild(ctx).PaginationModuleBuilder().ToAPI(params.Pagination),
	}, nil
}

// MarkAlarm 告警标记
func (s *AlarmService) MarkAlarm(ctx context.Context, req *realtimeapi.MarkAlarmRequest) (*realtimeapi.MarkAlarmReply, error) {
	params := builder.NewParamsBuild(ctx).RealtimeAlarmModuleBuilder().WithMarkAlarmRequest(req).ToBo()
	if err := s.alarmBiz.MarkRealTimeAlarm(ctx, params); err != nil {
		return nil, err
	}
	return &realtimeapi.MarkAlarmReply{}, nil
}
