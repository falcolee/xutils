/*
 * @Date: 2024-06-20 09:48:51
 * @LastEditTime: 2024-06-21 09:48:52
 * @Description:
 */
package memory

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func TestStore_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	store := New(1024 * 1024)
	store.Set(context.TODO(), "dddd", "dddd", time.Minute)
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "test1", args: args{ctx: context.TODO(), key: "dddd"}, want: "dddd", wantErr: false},
		{name: "test2", args: args{ctx: context.TODO(), key: "ddd2d"}, want: "dddd", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Store.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
