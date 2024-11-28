package gval

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_toFunc(t *testing.T) {
	myError := fmt.Errorf("my error")
	tests := []struct {
		name       string
		function   interface{}
		arguments  []interface{}
		want       interface{}
		wantErr    error
		wantAnyErr bool
	}{
		{
			name:     "empty",
			function: func() {},
		},
		{
			name: "one arg",
			function: func(a interface{}) {
				if a != true {
					panic("fail")
				}
			},
			arguments: []interface{}{true},
		},
		{
			name: "three args",
			function: func(a, b, c interface{}) {
				if a != 1 || b != 2 || c != 3 {
					panic("fail")
				}
			},
			arguments: []interface{}{1, 2, 3},
		},
		{
			name: "input types",
			function: func(a int, b string, c bool) {
				if a != 1 || b != "2" || !c {
					panic("fail")
				}
			},
			arguments: []interface{}{1, "2", true},
		},
		{
			name:       "wronge input type int",
			function:   func(a int, b string, c bool) {},
			arguments:  []interface{}{"1", "2", true},
			wantAnyErr: true,
		},
		{
			name:       "wronge input type string",
			function:   func(a int, b string, c bool) {},
			arguments:  []interface{}{1, 2, true},
			wantAnyErr: true,
		},
		{
			name:       "wronge input type bool",
			function:   func(a int, b string, c bool) {},
			arguments:  []interface{}{1, "2", "true"},
			wantAnyErr: true,
		},
		{
			name:       "wronge input number",
			function:   func(a int, b string, c bool) {},
			arguments:  []interface{}{1, "2"},
			wantAnyErr: true,
		},
		{
			name: "one return",
			function: func() bool {
				return true
			},
			want: true,
		},
		{
			name: "three returns",
			function: func() (bool, string, int) {
				return true, "2", 3
			},
			want: []interface{}{true, "2", 3},
		},
		{
			name: "error",
			function: func() error {
				return myError
			},
			wantErr: myError,
		},
		{
			name: "none error",
			function: func() error {
				return nil
			},
		},
		{
			name: "one return with error",
			function: func() (bool, error) {
				return false, myError
			},
			want:    false,
			wantErr: myError,
		},
		{
			name: "three returns with error",
			function: func() (bool, string, int, error) {
				return false, "", 0, myError
			},
			want:    []interface{}{false, "", 0},
			wantErr: myError,
		},
		{
			name: "context not expiring",
			function: func(ctx context.Context) error {
				return nil
			},
		},
		{
			name: "context expires",
			function: func(ctx context.Context) error {
				time.Sleep(20 * time.Millisecond)
				return nil
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "nil arg",
			function: func(a interface{}) bool {
				return a == nil
			},
			arguments: []interface{}{nil},
			want:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			got, err := toFunc(tt.function)(ctx, tt.arguments...)
			cancel()

			if tt.wantAnyErr {
				if err != nil {
					return
				}
				t.Fatalf("toFunc()(args...) = error(nil), but wantAnyErr")
			}
			if err != tt.wantErr {
				t.Fatalf("toFunc()(args...) = error(%v), wantErr (%v)", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toFunc()(args...) = %v, want %v", got, tt.want)
			}
		})
	}
}
