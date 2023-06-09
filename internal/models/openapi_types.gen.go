// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package models

// Metric defines model for Metric.
type Metric struct {
	Delta *int64   `json:"delta,omitempty"`
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Value *float64 `json:"value,omitempty"`
}

// UpdateMetricValueJSONRequestBody defines body for UpdateMetricValue for application/json ContentType.
type UpdateMetricValueJSONRequestBody = Metric

// GetMetricValueJSONRequestBody defines body for GetMetricValue for application/json ContentType.
type GetMetricValueJSONRequestBody = Metric
