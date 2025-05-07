package repository

import (
	"context"

	"github.com/aide-family/moon/cmd/server/palace/internal/biz/bo"
	"github.com/aide-family/moon/pkg/palace/imodel"
	"github.com/aide-family/moon/pkg/vobj"
)

// Resource 资源管理接口
type Resource interface {
	// GetByID get resource by id
	GetByID(context.Context, uint32) (imodel.IResource, error)

	// CheckPath check resource path
	CheckPath(context.Context, string) (imodel.IResource, error)

	// FindByPage find resource by page
	FindByPage(context.Context, *bo.QueryResourceListParams) ([]imodel.IResource, error)

	// UpdateStatus update resource status
	UpdateStatus(context.Context, vobj.Status, ...uint32) error

	// FindSelectByPage find select resource by page
	FindSelectByPage(context.Context, *bo.QueryResourceListParams) ([]imodel.IResource, error)
}
