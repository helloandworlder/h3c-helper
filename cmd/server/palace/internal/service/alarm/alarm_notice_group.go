package alarm

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	alarmyapi "github.com/aide-family/moon/api/admin/alarm"
	hookapi "github.com/aide-family/moon/api/rabbit/hook"
	"github.com/aide-family/moon/cmd/server/palace/internal/biz"
	"github.com/aide-family/moon/cmd/server/palace/internal/biz/bo"
	"github.com/aide-family/moon/cmd/server/palace/internal/service/builder"
	"github.com/aide-family/moon/pkg/helper/middleware"
	"github.com/aide-family/moon/pkg/merr"
	"github.com/aide-family/moon/pkg/util/types"
)

// GroupService 告警管理服务
type GroupService struct {
	alarmyapi.UnimplementedAlarmServer
	alarmGroupBiz *biz.AlarmGroupBiz
	alertBiz      *biz.AlarmBiz
}

// NewAlarmService 创建告警管理服务
func NewAlarmService(alarmGroupBiz *biz.AlarmGroupBiz, alertBiz *biz.AlarmBiz) *GroupService {
	return &GroupService{
		alarmGroupBiz: alarmGroupBiz,
		alertBiz:      alertBiz,
	}
}

// CreateAlarmGroup 创建告警组
func (s *GroupService) CreateAlarmGroup(ctx context.Context, req *alarmyapi.CreateAlarmGroupRequest) (*alarmyapi.CreateAlarmGroupReply, error) {
	// 校验通知人是否重复
	if has := types.SlicesHasDuplicates(req.GetNoticeMember(), func(request *alarmyapi.CreateNoticeMemberRequest) string {
		var sb strings.Builder
		sb.WriteString(strconv.FormatInt(int64(request.GetMemberId()), 10))
		return sb.String()
	}); has {
		return nil, merr.ErrorI18nAlertAlertObjectDuplicate(ctx)
	}
	param := builder.NewParamsBuild(ctx).
		AlarmNoticeGroupModuleBuilder().
		WithCreateAlarmGroupRequest(req).
		ToBo()
	if _, err := s.alarmGroupBiz.CreateAlarmGroup(ctx, param); !types.IsNil(err) {
		return nil, err
	}
	return &alarmyapi.CreateAlarmGroupReply{}, nil
}

// DeleteAlarmGroup 删除告警组
func (s *GroupService) DeleteAlarmGroup(ctx context.Context, req *alarmyapi.DeleteAlarmGroupRequest) (*alarmyapi.DeleteAlarmGroupReply, error) {
	if err := s.alarmGroupBiz.DeleteAlarmGroup(ctx, req.GetId()); !types.IsNil(err) {
		return nil, err
	}
	return &alarmyapi.DeleteAlarmGroupReply{}, nil
}

// ListAlarmGroup 获取告警组列表
func (s *GroupService) ListAlarmGroup(ctx context.Context, req *alarmyapi.ListAlarmGroupRequest) (*alarmyapi.ListAlarmGroupReply, error) {
	param := builder.NewParamsBuild(ctx).AlarmNoticeGroupModuleBuilder().WithListAlarmGroupRequest(req).ToBo()
	alarmGroups, err := s.alarmGroupBiz.ListPage(ctx, param)
	if !types.IsNil(err) {
		return nil, err
	}
	return &alarmyapi.ListAlarmGroupReply{
		Pagination: builder.NewParamsBuild(ctx).PaginationModuleBuilder().ToAPI(param.Page),
		List:       builder.NewParamsBuild(ctx).AlarmNoticeGroupModuleBuilder().DoAlarmNoticeGroupItemBuilder().ToAPIs(alarmGroups),
	}, nil
}

// GetAlarmGroup 获取告警组详细信息
func (s *GroupService) GetAlarmGroup(ctx context.Context, req *alarmyapi.GetAlarmGroupRequest) (*alarmyapi.GetAlarmGroupReply, error) {
	detail, err := s.alarmGroupBiz.GetAlarmGroupDetail(ctx, req.GetId())
	if !types.IsNil(err) {
		return nil, err
	}
	return &alarmyapi.GetAlarmGroupReply{
		Detail: builder.NewParamsBuild(ctx).AlarmNoticeGroupModuleBuilder().DoAlarmNoticeGroupItemBuilder().ToAPI(detail),
	}, nil
}

// UpdateAlarmGroup 更新告警组信息
func (s *GroupService) UpdateAlarmGroup(ctx context.Context, req *alarmyapi.UpdateAlarmGroupRequest) (*alarmyapi.UpdateAlarmGroupReply, error) {
	// 校验通知人是否重复
	if has := types.SlicesHasDuplicates(req.GetUpdate().GetNoticeMember(), func(request *alarmyapi.CreateNoticeMemberRequest) string {
		var sb strings.Builder
		sb.WriteString(strconv.FormatInt(int64(request.GetMemberId()), 10))
		return sb.String()
	}); has {
		return nil, merr.ErrorI18nAlertAlertObjectDuplicate(ctx)
	}
	param := builder.NewParamsBuild(ctx).AlarmNoticeGroupModuleBuilder().WithUpdateAlarmGroupRequest(req).ToBo()
	err := s.alarmGroupBiz.UpdateAlarmGroup(ctx, param)
	if !types.IsNil(err) {
		return nil, err
	}
	return &alarmyapi.UpdateAlarmGroupReply{}, nil
}

// UpdateAlarmGroupStatus 更新告警组状态
func (s *GroupService) UpdateAlarmGroupStatus(ctx context.Context, req *alarmyapi.UpdateAlarmGroupStatusRequest) (*alarmyapi.UpdateAlarmGroupStatusReply, error) {
	param := builder.NewParamsBuild(ctx).AlarmNoticeGroupModuleBuilder().WithUpdateAlarmGroupStatusRequest(req).ToBo()
	err := s.alarmGroupBiz.UpdateStatus(ctx, param)
	if !types.IsNil(err) {
		return nil, err
	}
	return &alarmyapi.UpdateAlarmGroupStatusReply{}, nil
}

// ListAlarmGroupSelect 获取告警组下拉列表
func (s *GroupService) ListAlarmGroupSelect(ctx context.Context, req *alarmyapi.ListAlarmGroupRequest) (*alarmyapi.ListAlarmGroupSelectReply, error) {
	param := builder.NewParamsBuild(ctx).AlarmNoticeGroupModuleBuilder().WithListAlarmGroupRequest(req).ToBo()
	alarmGroups, err := s.alarmGroupBiz.ListPage(ctx, param)
	if !types.IsNil(err) {
		return nil, err
	}
	return &alarmyapi.ListAlarmGroupSelectReply{
		List: builder.NewParamsBuild(ctx).AlarmNoticeGroupModuleBuilder().DoAlarmNoticeGroupItemBuilder().ToSelects(alarmGroups),
	}, nil
}

// MyAlarmGroupList 获取我的告警组
func (s *GroupService) MyAlarmGroupList(ctx context.Context, req *alarmyapi.MyAlarmGroupRequest) (*alarmyapi.MyAlarmGroupReply, error) {
	param := builder.NewParamsBuild(ctx).
		AlarmNoticeGroupModuleBuilder().
		WithAPIMyAlarmGroupListRequest(req).
		ToBo()

	myAlarmGroup, err := s.alarmGroupBiz.MyAlarmGroups(ctx, param)
	if !types.IsNil(err) {
		return nil, err
	}

	return &alarmyapi.MyAlarmGroupReply{
		Pagination: builder.NewParamsBuild(ctx).
			PaginationModuleBuilder().
			ToAPI(param.Page),
		List: builder.NewParamsBuild(ctx).
			AlarmNoticeGroupModuleBuilder().
			DoAlarmNoticeGroupItemBuilder().
			ToAPIs(myAlarmGroup),
	}, nil
}

// MessageTest 发送测试消息
func (s *GroupService) MessageTest(ctx context.Context, req *alarmyapi.MessageTestRequest) (*alarmyapi.MessageTestReply, error) {
	msg := &bo.SendMsg{
		SendMsgRequest: &hookapi.SendMsgRequest{
			Json:      req.GetMessage(),
			Route:     fmt.Sprintf("team_%d_%d", middleware.GetTeamID(ctx), req.GetId()),
			RequestID: types.MD5(time.Now().Format("2006-01-02 15:04")), // 限定一分钟发送一次
		},
	}
	s.alertBiz.SendAlertMsg(ctx, msg)
	return &alarmyapi.MessageTestReply{}, nil
}
