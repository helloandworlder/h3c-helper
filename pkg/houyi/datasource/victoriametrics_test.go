package datasource_test

import (
	"context"
	"testing"
	"time"

	"github.com/aide-family/moon/pkg/houyi/datasource"
	"github.com/aide-family/moon/pkg/util/types"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestNewVictoriametricsDatasource(t *testing.T) {
	opts := []datasource.VictoriaMetricsDatasourceOption{
		datasource.WithVictoriaMetricsEndpoint("https://victoriametrics.aide-cloud.cn"),
	}
	vmData := datasource.NewVictoriaMetricsDatasource(opts...)
	durationT := types.NewDuration(durationpb.New(60 * time.Second))
	endAt := time.Now()
	startAt := types.NewTime(endAt.Add(-durationT.Duration.AsDuration()))
	queryRange, err := vmData.QueryRange(context.Background(), "up", startAt.Unix(), endAt.Unix(), 10)
	if err != nil {
		t.Fatal(err)
	}
	bs, _ := types.Marshal(queryRange)
	t.Log(string(bs))

	eval, err := datasource.MetricEval(vmData)(context.Background(), "up", durationT)
	if err != nil {
		t.Fatal(err)
	}
	for indexer, point := range eval {
		t.Log("idnex", indexer)
		pointBs, _ := types.Marshal(point)
		t.Log("point", string(pointBs))
	}
}
