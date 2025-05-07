package builder

import (
	"context"

	"github.com/aide-family/moon/api"
	"github.com/aide-family/moon/cmd/server/palace/internal/biz/bo"
	"github.com/aide-family/moon/pkg/util/types"
)

var _ IMetricDataModuleBuilder = (*metricDataModuleBuilder)(nil)

type (
	metricDataModuleBuilder struct {
		ctx context.Context
	}

	// IMetricDataModuleBuilder 指标数据模块构造器
	IMetricDataModuleBuilder interface {
		// BoMetricDataBuilder 指标数据条目构造器
		BoMetricDataBuilder() IBoMetricDataBuilder
		// BoMetricDataValueBuilder 指标数据值条目构造器
		BoMetricDataValueBuilder() IBoMetricDataValueBuilder
	}

	// IBoMetricDataBuilder 指标数据条目构造器
	IBoMetricDataBuilder interface {
		// ToAPI 转换为API对象
		ToAPI(*bo.MetricQueryData) *api.MetricQueryResult
		// ToAPIs 转换为API对象列表
		ToAPIs([]*bo.MetricQueryData) []*api.MetricQueryResult
	}

	boMetricDataBuilder struct {
		ctx context.Context
	}

	// IBoMetricDataValueBuilder 指标数据值条目构造器
	IBoMetricDataValueBuilder interface {
		// ToAPI 转换为API对象
		ToAPI(*bo.DatasourceQueryValue) *api.MetricQueryValue
		// ToAPIs 转换为API对象列表
		ToAPIs([]*bo.DatasourceQueryValue) []*api.MetricQueryValue
	}

	boMetricDataValueBuilder struct {
		ctx context.Context
	}
)

func (b *boMetricDataValueBuilder) ToAPI(value *bo.DatasourceQueryValue) *api.MetricQueryValue {
	if types.IsNil(b) || types.IsNil(value) {
		return nil
	}

	return &api.MetricQueryValue{
		Value:     value.Value,
		Timestamp: value.Timestamp,
	}
}

func (b *boMetricDataValueBuilder) ToAPIs(values []*bo.DatasourceQueryValue) []*api.MetricQueryValue {
	if types.IsNil(b) || types.IsNil(values) {
		return nil
	}

	return types.SliceTo(values, func(value *bo.DatasourceQueryValue) *api.MetricQueryValue {
		return b.ToAPI(value)
	})
}

func (m *metricDataModuleBuilder) BoMetricDataValueBuilder() IBoMetricDataValueBuilder {
	return &boMetricDataValueBuilder{ctx: m.ctx}
}

func (b *boMetricDataBuilder) ToAPI(data *bo.MetricQueryData) *api.MetricQueryResult {
	if types.IsNil(b) || types.IsNil(data) {
		return nil
	}

	return &api.MetricQueryResult{
		Labels:     data.Labels,
		ResultType: data.ResultType,
		Values:     NewParamsBuild(b.ctx).MetricDataModuleBuilder().BoMetricDataValueBuilder().ToAPIs(data.Values),
		Value:      NewParamsBuild(b.ctx).MetricDataModuleBuilder().BoMetricDataValueBuilder().ToAPI(data.Value),
	}
}

func (b *boMetricDataBuilder) ToAPIs(data []*bo.MetricQueryData) []*api.MetricQueryResult {
	if types.IsNil(b) || types.IsNil(data) {
		return nil
	}

	return types.SliceTo(data, func(item *bo.MetricQueryData) *api.MetricQueryResult {
		return b.ToAPI(item)
	})
}

func (m *metricDataModuleBuilder) BoMetricDataBuilder() IBoMetricDataBuilder {
	return &boMetricDataBuilder{ctx: m.ctx}
}
