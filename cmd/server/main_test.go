package main

import (
	"net/http"
	"testing"
)

func Test_collectorHandler(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		method string
		want   int
	}{
		{
			name:   "incorrect method",
			path:   "/update/gauge/g/3.14",
			method: "GET",
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "incorrect type",
			path:   "/update/test/g/3.14",
			method: "POST",
			want:   http.StatusBadRequest,
		},
		{
			name:   "correct path gauge",
			path:   "/update/gauge/g/3.14",
			method: "POST",
			want:   http.StatusOK,
		},
		// {
		// 	name:   "incorrect path gauge: without name 1",
		// 	path:   "/update/gauge//3.14",
		// 	method: "POST",
		// 	want:   http.StatusNotFound,
		// },
		{
			name:   "incorrect path gauge: without name 2",
			path:   "/update/gauge/3.14",
			method: "POST",
			want:   http.StatusNotFound,
		},
		{
			name:   "incorrect path gauge: incorrect value 1",
			path:   "/update/gauge/g/text",
			method: "POST",
			want:   http.StatusBadRequest,
		},
		{
			name:   "correct path counter",
			path:   "/update/counter/c/5",
			method: "POST",
			want:   http.StatusOK,
		},
		// {
		// 	name:   "incorrect path counter: without name 1",
		// 	path:   "/update/counter//5",
		// 	method: "POST",
		// 	want:   http.StatusNotFound,
		// },
		{
			name:   "incorrect path counter: without name 2",
			path:   "/update/counter/5",
			method: "POST",
			want:   http.StatusNotFound,
		},
		{
			name:   "incorrect path counter: incorrect value 1",
			path:   "/update/counter/g/3.14",
			method: "POST",
			want:   http.StatusBadRequest,
		},
		{
			name:   "incorrect path counter: incorrect value 2",
			path:   "/update/counter/g/text",
			method: "POST",
			want:   http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// r := httptest.NewRequest(tt.method, tt.path, nil)
			// w := httptest.NewRecorder()
			// collectorHandler(w, r)
			// result := w.Result()
			// defer result.Body.Close()
			// assert.Equal(t, tt.want, result.StatusCode)
		})
	}
}
