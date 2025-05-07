package builder

import (
	"context"

	"github.com/aide-family/moon/api"
	adminapi "github.com/aide-family/moon/api/admin"
	resourceapi "github.com/aide-family/moon/api/admin/resource"
	"github.com/aide-family/moon/cmd/server/palace/internal/biz/bo"
	"github.com/aide-family/moon/pkg/palace/imodel"
	"github.com/aide-family/moon/pkg/util/types"
	"github.com/aide-family/moon/pkg/vobj"
)

var _ IResourceModuleBuilder = (*resourceModuleBuilder)(nil)

type (
	resourceModuleBuilder struct {
		ctx context.Context
	}

	// IResourceModuleBuilder 资源模块构造器
	IResourceModuleBuilder interface {
		// WithListResourceRequest 设置获取资源列表请求参数
		WithListResourceRequest(*resourceapi.ListResourceRequest) IListResourceRequestBuilder
		// DoResourceBuilder 资源条目构造器
		DoResourceBuilder() IDoResourceBuilder
	}

	// IListResourceRequestBuilder 获取资源列表请求参数构造器
	IListResourceRequestBuilder interface {
		// ToBo 转换为业务对象
		ToBo() *bo.QueryResourceListParams
	}

	listResourceRequestBuilder struct {
		ctx context.Context
		*resourceapi.ListResourceRequest
	}

	// IDoResourceBuilder 资源条目构造器
	IDoResourceBuilder interface {
		// ToAPI 转换为API对象
		ToAPI(imodel.IResource) *adminapi.ResourceItem
		// ToAPIs 转换为API对象列表
		ToAPIs([]imodel.IResource) []*adminapi.ResourceItem
		// ToSelect 转换为选择对象
		ToSelect(imodel.IResource) *adminapi.SelectItem
		// ToSelects 转换为选择对象列表
		ToSelects([]imodel.IResource) []*adminapi.SelectItem
	}

	doResourceBuilder struct {
		ctx context.Context
	}
)

func (d *doResourceBuilder) ToAPI(sysAPI imodel.IResource) *adminapi.ResourceItem {
	if types.IsNil(d) || types.IsNil(sysAPI) {
		return nil
	}

	return &adminapi.ResourceItem{
		Id:        sysAPI.GetID(),
		Name:      sysAPI.GetName(),
		Path:      sysAPI.GetPath(),
		Status:    api.Status(sysAPI.GetStatus()),
		Remark:    sysAPI.GetRemark(),
		CreatedAt: sysAPI.GetCreatedAt().String(),
		UpdatedAt: sysAPI.GetUpdatedAt().String(),
		Module:    api.ModuleType(sysAPI.GetModule()),
		Domain:    api.DomainType(sysAPI.GetDomain()),
	}
}

func (d *doResourceBuilder) ToAPIs(apis []imodel.IResource) []*adminapi.ResourceItem {
	if types.IsNil(d) || types.IsNil(apis) {
		return nil
	}

	return types.SliceTo(apis, func(api imodel.IResource) *adminapi.ResourceItem {
		return d.ToAPI(api)
	})
}

func (d *doResourceBuilder) ToSelect(sysAPI imodel.IResource) *adminapi.SelectItem {
	if types.IsNil(d) || types.IsNil(sysAPI) {
		return nil
	}

	return &adminapi.SelectItem{
		Value:    sysAPI.GetID(),
		Label:    sysAPI.GetName(),
		Children: nil,
		Disabled: sysAPI.GetDeletedAt() > 0 || !sysAPI.GetStatus().IsEnable(),
		Extend: &adminapi.SelectExtend{
			Remark: sysAPI.GetRemark(),
		},
	}
}

func (d *doResourceBuilder) ToSelects(apis []imodel.IResource) []*adminapi.SelectItem {
	if types.IsNil(d) || types.IsNil(apis) {
		return nil
	}

	return types.SliceTo(apis, func(api imodel.IResource) *adminapi.SelectItem {
		return d.ToSelect(api)
	})
}

func (l *listResourceRequestBuilder) ToBo() *bo.QueryResourceListParams {
	if types.IsNil(l) || types.IsNil(l.ListResourceRequest) {
		return nil
	}

	return &bo.QueryResourceListParams{
		Keyword: l.GetKeyword(),
		Page:    types.NewPagination(l.GetPagination()),
		IsAll:   l.GetIsAll(),
		Status:  vobj.Status(l.GetStatus()),
	}
}

func (r *resourceModuleBuilder) WithListResourceRequest(request *resourceapi.ListResourceRequest) IListResourceRequestBuilder {
	return &listResourceRequestBuilder{ctx: r.ctx, ListResourceRequest: request}
}

func (r *resourceModuleBuilder) DoResourceBuilder() IDoResourceBuilder {
	return &doResourceBuilder{ctx: r.ctx}
}
