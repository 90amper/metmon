package utils

import (
	"fmt"
	"testing"
)

func TestRetryer(t *testing.T) {
	// type args struct {
	// 	fn func() error
	// }

	called := 0
	tests := []struct {
		name    string
		fn      func() error
		wantErr bool
	}{
		{
			name: "aaaa",
			fn: func() error {
				err := fmt.Errorf("caller %d", called)
				fmt.Println(err.Error())
				if called == 2 { // only 5th call returns ok
					return nil
				}
				called++
				return err

			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Retryer(tt.fn); (err != nil) != tt.wantErr {
				t.Errorf("Retryer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
