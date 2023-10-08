package poster

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ccfos/nightingale/v6/conf"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

func TestPostByUrls(t *testing.T) {
	type args struct {
		ctx  *ctx.Context
		path string
		v    interface{}
	}

	info := struct {
		a string
		b int
	}{
		a: "aaa",
		b: 888,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := DataResponse[interface{}]{Dat: "", Err: ""}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	tests := []struct {
		name string
		args args
	}{
		{"a", args{ctx: &ctx.Context{
			CenterApi: conf.CenterApi{
				Addrs: []string{server.URL},
			},
		}, path: "/v1/n9e/server-heartbeat", v: info}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PostByUrls(tt.args.ctx, tt.args.path, tt.args.v); err != nil {
				t.Errorf("PostByUrls() error = %v ", err)
			}
		})
	}
}

func TestPostByUrlsWithResp(t *testing.T) {
	type args struct {
		ctx  *ctx.Context
		path string
		v    interface{}
	}
	type testCase[T any] struct {
		name string
		args args
	}
	info := struct {
		a string
		b int
	}{
		a: "aaa",
		b: 888,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := DataResponse[int64]{Dat: 123, Err: ""}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	tests := []testCase[int64]{{"a-resp", args{ctx: &ctx.Context{
		CenterApi: conf.CenterApi{
			Addrs: []string{server.URL},
		}}, path: "/v1/n9e/event-persist",
		v: info}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := PostByUrlsWithResp[int64](tt.args.ctx, tt.args.path, tt.args.v)
			if err != nil {
				t.Errorf("PostByUrlsWithResp() error = %v", err)
				return
			}
			t.Logf("PostByUrlsWithResp() gotT = %v", gotT)
		})
	}
}
