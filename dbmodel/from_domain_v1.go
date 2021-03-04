package dbmodel

import (
	"fmt"
	"strings"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/influxdata/influxdb1-client/models"
	"github.com/kitabisa/jaeger-influxdb/common"
	"github.com/jaegertracing/jaeger/model"
)

// SpanToPointsV1 converts a Jaeger span to InfluxDB v1.x points.
// One point for the span itself, and one point for each log entry on the span.
func SpanToPointsV1(span *model.Span, spanMeasurement, logMeasurement string, logger hclog.Logger) ([]models.Point, error) {
	var tags models.Tags

	tags.SetString(common.TraceIDKey, span.TraceID.String())
	tags.SetString(common.ServiceNameKey, span.Process.ServiceName)
	tags.SetString(common.OperationNameKey, span.OperationName)

	for _, tag := range append(span.Tags, span.Process.Tags...) {
		key, value, err := keyValueAsStrings(&tag)
		if err != nil {
			logger.Warn(err.Error(),
				"skipped-key-and-type", fmt.Sprintf("%s:%s", tag.Key, tag.VType.String()))
			continue
		}

		tags.SetString(key, value)
	}

	fields := models.Fields{}

	fields[common.SpanIDKey] = span.SpanID.String()
	// The 3 least significant digits are always 0. Jaeger uses µs, not ns
	fields[common.DurationKey] = span.Duration.Nanoseconds()
	fields[common.FlagsKey] = uint32(span.Flags)

	var processTagKeys []string
	for _, tag := range span.Process.Tags {
		processTagKeys = append(processTagKeys, tag.Key)
	}
	if len(processTagKeys) > 0 {
		// TODO escape commas
		fields[common.ProcessTagKeysKey] = strings.Join(processTagKeys, ",")
	}

	var references []string
	for _, spanRef := range span.References {
		if spanRef.SpanID == 0 {
			continue
		}

		var referenceType string
		switch spanRef.RefType {
		case model.SpanRefType_CHILD_OF:
			referenceType = common.ReferenceTypeChildOf
		case model.SpanRefType_FOLLOWS_FROM:
			referenceType = common.ReferenceTypeFollowsFrom
		default:
			logger.Warn("skipped unrecognized span reference type",
				"skipped-spanref-id-and-type", fmt.Sprintf("%s:%s", spanRef.SpanID.String(), spanRef.RefType.String()))
			continue
		}
		references = append(references, fmt.Sprintf("%s:%s", spanRef.SpanID.String(), referenceType))
	}
	if len(references) > 0 {
		// TODO escape colons and commas
		fields[common.ReferencesKey] = strings.Join(references, ",")
	}

	startTime := mergeTimeAndSpanID(span.StartTime, span.SpanID)
	spanPoint, err := models.NewPoint(spanMeasurement, tags, fields, startTime)
	if err != nil {
		return nil, err
	}
	points := append(make([]models.Point, 0, len(span.Logs)+1), spanPoint)

	if len(span.Logs) > 0 {
		var tags models.Tags
		tags.SetString(common.TraceIDKey, span.TraceID.String())

		for i := range span.Logs {
			spanLog := &span.Logs[i]

			fields := make(map[string]interface{}, len(spanLog.Fields)+1)
			fields[common.SpanIDKey] = span.SpanID.String()
			for j := 0; j < len(spanLog.Fields); j++ {
				spanField := &spanLog.Fields[j]
				key, value, err := keyValueAsStringAndInterface(spanField)
				if err != nil {
					logger.Warn("skipping span log field",
						common.TraceIDKey, span.TraceID.String(),
						common.SpanIDKey, span.SpanID.String(),
						"field-key", spanField.Key,
						"error", err)
					continue
				}
				if key == common.TraceIDKey || key == common.SpanIDKey {
					logger.Warn("skipping span log field because field key is reserved",
						common.TraceIDKey, span.TraceID.String(),
						common.SpanIDKey, span.SpanID.String(),
						"field-key", spanField.Key)
					continue
				}
				fields[key] = value
			}

			point, err := models.NewPoint(logMeasurement, tags, fields, spanLog.Timestamp)
			if err != nil {
				logger.Warn("skipping span log", common.SpanIDKey, span.SpanID.String())
				continue
			}
			points = append(points, point)
		}
	}

	return points, nil
}
