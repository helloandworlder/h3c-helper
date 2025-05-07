package builder

import (
	"context"

	"github.com/aide-family/moon/api"
	adminapi "github.com/aide-family/moon/api/admin"
	"github.com/aide-family/moon/cmd/server/palace/internal/biz"
	"github.com/aide-family/moon/pkg/palace/model"
	"github.com/aide-family/moon/pkg/palace/model/bizmodel"
	"github.com/aide-family/moon/pkg/util/types"
)

var _ ITeamMemberModuleBuilder = (*teamMemberModuleBuilder)(nil)

type (
	teamMemberModuleBuilder struct {
		ctx context.Context
	}

	// ITeamMemberModuleBuilder 团队成员模块条目构造器
	ITeamMemberModuleBuilder interface {
		// DoTeamMemberBuilder 获取团队成员条目构造器
		DoTeamMemberBuilder() IDoTeamMemberBuilder
	}

	// IDoTeamMemberBuilder 团队成员条目构造器
	IDoTeamMemberBuilder interface {
		// ToAPI 转换为API对象
		ToAPI(*bizmodel.SysTeamMember) *adminapi.TeamMemberItem
		// ToAPIs 转换为API对象列表
		ToAPIs([]*bizmodel.SysTeamMember) []*adminapi.TeamMemberItem
		// ToSelect 转换为选择对象
		ToSelect(*bizmodel.SysTeamMember) *adminapi.SelectItem
		// ToSelects 转换为选择对象列表
		ToSelects([]*bizmodel.SysTeamMember) []*adminapi.SelectItem
	}

	doTeamMemberBuilder struct {
		ctx context.Context
	}
)

func (d *doTeamMemberBuilder) ToAPI(member *bizmodel.SysTeamMember) *adminapi.TeamMemberItem {
	if types.IsNil(d) || types.IsNil(member) {
		return nil
	}

	return &adminapi.TeamMemberItem{
		UserId:    member.UserID,
		Id:        member.ID,
		Role:      api.Role(member.Role),
		Status:    api.Status(member.Status),
		CreatedAt: member.CreatedAt.String(),
		UpdatedAt: member.UpdatedAt.String(),
		User:      NewParamsBuild(d.ctx).UserModuleBuilder().DoUserBuilder().ToAPI(biz.RuntimeCache.GetUser(d.ctx, member.UserID)),
	}
}

func (d *doTeamMemberBuilder) ToAPIs(members []*bizmodel.SysTeamMember) []*adminapi.TeamMemberItem {
	if types.IsNil(d) || types.IsNil(members) {
		return nil
	}

	return types.SliceTo(members, func(member *bizmodel.SysTeamMember) *adminapi.TeamMemberItem {
		return d.ToAPI(member)
	})
}

func (d *doTeamMemberBuilder) ToSelect(member *bizmodel.SysTeamMember) *adminapi.SelectItem {
	if types.IsNil(d) || types.IsNil(member) {
		return nil
	}
	userInfo := new(model.SysUser) // TODO 获取实际的user
	return &adminapi.SelectItem{
		Value:    member.ID,
		Label:    userInfo.Username,
		Children: nil,
		Disabled: member.DeletedAt > 0 || !member.Status.IsEnable() || userInfo.DeletedAt > 0 || !userInfo.Status.IsEnable(),
		Extend: &adminapi.SelectExtend{
			Remark: userInfo.Remark,
			Image:  userInfo.Avatar,
		},
	}
}

func (d *doTeamMemberBuilder) ToSelects(members []*bizmodel.SysTeamMember) []*adminapi.SelectItem {
	if types.IsNil(d) || types.IsNil(members) {
		return nil
	}

	return types.SliceTo(members, func(member *bizmodel.SysTeamMember) *adminapi.SelectItem {
		return d.ToSelect(member)
	})
}

func (t *teamMemberModuleBuilder) DoTeamMemberBuilder() IDoTeamMemberBuilder {
	return &doTeamMemberBuilder{ctx: t.ctx}
}
