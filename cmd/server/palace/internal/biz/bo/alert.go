package bo

import (
	"github.com/aide-family/moon/pkg/label"
	"github.com/aide-family/moon/pkg/util/types"
	"github.com/aide-family/moon/pkg/vobj"
	"github.com/aide-family/moon/pkg/watch"
)

var (
	_ watch.Indexer = (*Alarm)(nil)
	_ watch.Indexer = (*Alert)(nil)
)

type (
	// Alarm alarm detail info
	Alarm struct {
		Receiver          string            `json:"receiver"`
		Status            vobj.AlertStatus  `json:"status"`
		Alerts            []*Alert          `json:"alerts"`
		GroupLabels       *label.Labels     `json:"groupLabels"`
		CommonLabels      *label.Labels     `json:"commonLabels"`
		CommonAnnotations label.Annotations `json:"commonAnnotations"`
		ExternalURL       string            `json:"externalURL"`
		Version           string            `json:"version"`
		GroupKey          string            `json:"groupKey"`
		TruncatedAlerts   int32             `json:"truncatedAlerts"`
	}

	alarmInfo struct {
		Receiver          string            `json:"receiver"`
		Status            string            `json:"status"`
		Alerts            []*alertInfo      `json:"alerts"`
		GroupLabels       map[string]string `json:"groupLabels"`
		CommonLabels      map[string]string `json:"commonLabels"`
		CommonAnnotations map[string]string `json:"commonAnnotations"`
		ExternalURL       string            `json:"externalURL"`
		Version           string            `json:"version"`
		GroupKey          string            `json:"groupKey"`
	}

	// Alert alert detail info
	Alert struct {
		Status       vobj.AlertStatus  `json:"status"`
		Labels       *label.Labels     `json:"labels"`
		Annotations  label.Annotations `json:"annotations"`
		StartsAt     *types.Time       `json:"startsAt"`
		EndsAt       *types.Time       `json:"endsAt"`
		GeneratorURL string            `json:"generatorURL"`
		Fingerprint  string            `json:"fingerprint"`
		Value        float64           `json:"value"`
	}

	alertInfo struct {
		Status       string            `json:"status"`
		Labels       map[string]string `json:"labels"`
		Annotations  map[string]string `json:"annotations"`
		StartsAt     string            `json:"startsAt"`
		EndsAt       string            `json:"endsAt"`
		GeneratorURL string            `json:"generatorURL"`
		Fingerprint  string            `json:"fingerprint"`
		Value        float64           `json:"value"`
	}

	// LatestAlarmEvent 最新事件
	LatestAlarmEvent struct {
		Fingerprint string           `json:"fingerprint"`
		Level       string           `json:"level"`
		EventTime   string           `json:"eventTime"`
		Summary     string           `json:"summary"`
		Status      vobj.AlertStatus `json:"status"`
		TeamID      string           `json:"teamID"`
	}

	// LatestInterventionEvent 最新干预事件
	LatestInterventionEvent struct {
		Fingerprint string           `json:"fingerprint"`
		Level       string           `json:"level"`
		EventTime   string           `json:"eventTime"`
		Summary     string           `json:"summary"`
		Status      vobj.AlertStatus `json:"status"`
		TeamID      string           `json:"teamID"`
		HandlerID   uint32           `json:"handlerID"`
		HandledAt   string           `json:"handledAt"`
	}
)

// BinaryUnmarshaler LatestAlarmEvent encoding.BinaryUnmarshaler
func (a *LatestAlarmEvent) UnmarshalBinary(data []byte) error {
	return types.Unmarshal(data, a)
}

// MarshalBinary LatestAlarmEvent encoding.BinaryMarshaler
func (a *LatestAlarmEvent) MarshalBinary() ([]byte, error) {
	return types.Marshal(a)
}

// BinaryUnmarshaler LatestInterventionEvent encoding.BinaryUnmarshaler
func (a *LatestInterventionEvent) UnmarshalBinary(data []byte) error {
	return types.Unmarshal(data, a)
}

// MarshalBinary LatestInterventionEvent encoding.BinaryMarshaler
func (a *LatestInterventionEvent) MarshalBinary() ([]byte, error) {
	return types.Marshal(a)
}

func (a *Alarm) String() string {
	alarm := alarmInfo{
		Receiver: a.Receiver,
		Status:   a.Status.String(),
		Alerts: types.SliceTo(a.Alerts, func(alert *Alert) *alertInfo {
			return &alertInfo{
				Status:       alert.Status.String(),
				Labels:       alert.Labels.Map(),
				Annotations:  alert.Annotations,
				StartsAt:     alert.StartsAt.String(),
				EndsAt:       alert.EndsAt.String(),
				GeneratorURL: alert.GeneratorURL,
				Fingerprint:  alert.Fingerprint,
				Value:        alert.Value,
			}
		}),
		GroupLabels:       a.GroupLabels.Map(),
		CommonLabels:      a.CommonLabels.Map(),
		CommonAnnotations: a.CommonAnnotations,
		ExternalURL:       a.ExternalURL,
		Version:           a.Version,
		GroupKey:          a.GroupKey,
	}
	bs, _ := types.Marshal(alarm)
	return string(bs)
}

func (a *Alert) String() string {
	alert := alertInfo{
		Status:       a.Status.String(),
		Labels:       a.Labels.Map(),
		Annotations:  a.Annotations,
		StartsAt:     a.StartsAt.String(),
		EndsAt:       a.EndsAt.String(),
		GeneratorURL: a.GeneratorURL,
		Fingerprint:  a.Fingerprint,
		Value:        a.Value,
	}
	bs, _ := types.Marshal(alert)
	return string(bs)
}

// Index gen alert index
func (a *Alert) Index() string {
	return types.TextJoin("houyi:alert:", types.MD5(a.Labels.String()))
}

// Index gen alarm index
func (a *Alarm) Index() string {
	return types.TextJoin("houyi:alarm:", types.MD5(a.GroupLabels.String()))
}

// Message gen alarm message
func (a *Alarm) Message() *watch.Message {
	return watch.NewMessage(a, vobj.TopicAlarm, watch.WithMessageRetryMax(3))
}

// Message gen alert message
func (a *Alert) Message() *watch.Message {
	return watch.NewMessage(a, vobj.TopicAlert, watch.WithMessageRetryMax(3))
}
