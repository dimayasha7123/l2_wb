package main

import (
	"reflect"
	"testing"
)

func Test_sortUtility(t *testing.T) {
	type args struct {
		data []byte
		set  settings
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sortUtility(tt.args.data, tt.args.set)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortUtility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortUtility() got = %v, want %v", got, tt.want)
			}
		})
	}
}
