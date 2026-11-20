package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
)

func TestEasyPayQueryOrderStatusMapping(t *testing.T) {
	t.Parallel()

	const orderID = "order-123"
	wantQuery := url.Values{
		"act": {"order"}, "pid": {"pid-1"}, "key": {"pkey-1"}, "out_trade_no": {orderID},
	}
	tests := []struct {
		name        string
		body        string
		wantStatus  string
		wantTradeNo string
		wantAmount  float64
		wantErr     bool
	}{
		{
			name:        "top level trade success is paid",
			body:        `{"code":1,"trade_status":"TRADE_SUCCESS","status":0,"money":"12.34","trade_no":"gateway-123"}`,
			wantStatus:  payment.ProviderStatusPaid,
			wantTradeNo: "gateway-123",
			wantAmount:  12.34,
		},
		{
			name:        "waiting trade status with paid numeric status stays pending",
			body:        `{"code":1,"trade_status":"WAITING","status":1,"money":"12.34","trade_no":"gateway-123"}`,
			wantStatus:  payment.ProviderStatusPending,
			wantTradeNo: "gateway-123",
			wantAmount:  12.34,
		},
		{
			name:        "empty trade status with paid numeric status stays pending",
			body:        `{"code":1,"trade_status":"","status":1,"money":"12.34"}`,
			wantStatus:  payment.ProviderStatusPending,
			wantTradeNo: orderID,
			wantAmount:  12.34,
		},
		{
			name:        "nested data trade success is paid",
			body:        `{"code":1,"data":{"trade_status":"TRADE_SUCCESS","status":0,"money":"9.99","trade_no":"data-456"}}`,
			wantStatus:  payment.ProviderStatusPaid,
			wantTradeNo: "data-456",
			wantAmount:  9.99,
		},
		{
			name:        "legacy numeric paid status remains compatible",
			body:        `{"code":1,"status":1,"money":"3.21"}`,
			wantStatus:  payment.ProviderStatusPaid,
			wantTradeNo: orderID,
			wantAmount:  3.21,
		},
		{
			name:        "legacy numeric non paid status is pending",
			body:        `{"code":1,"status":0,"money":"3.21"}`,
			wantStatus:  payment.ProviderStatusPending,
			wantTradeNo: orderID,
			wantAmount:  3.21,
		},
		{
			name:    "query failure code is an error",
			body:    `{"code":0,"msg":"订单不存在"}`,
			wantErr: true,
		},
		{
			name:    "explicit failure code wins over paid status fields",
			body:    `{"code":-1,"msg":"签名错误","status":1,"trade_status":"TRADE_SUCCESS","money":"12.34"}`,
			wantErr: true,
		},
		{
			name:    "string failure code is an error",
			body:    `{"code":"-1","msg":"error","status":1}`,
			wantErr: true,
		},
		{
			name:    "foreign order echo is rejected",
			body:    `{"code":1,"pid":"pid-1","out_trade_no":"order-999","trade_status":"TRADE_SUCCESS","money":"12.34"}`,
			wantErr: true,
		},
		{
			name:    "foreign order echo in nested data is rejected",
			body:    `{"code":1,"data":{"out_trade_no":"order-999","trade_status":"TRADE_SUCCESS","money":"12.34"}}`,
			wantErr: true,
		},
		{
			name:        "string numeric status from PHP deployments is paid",
			body:        `{"code":"1","msg":"查询订单号成功！","trade_no":"2026101600001","out_trade_no":"order-123","pid":"1001","money":"12.34","status":"1"}`,
			wantStatus:  payment.ProviderStatusPaid,
			wantTradeNo: "2026101600001",
			wantAmount:  12.34,
		},
		{
			name:        "string numeric unpaid status is pending",
			body:        `{"code":1,"trade_no":"2026101600001","money":"12.34","status":"0"}`,
			wantStatus:  payment.ProviderStatusPending,
			wantTradeNo: "2026101600001",
			wantAmount:  12.34,
		},
		{
			name:        "documented caihong response with integer pid and status is paid",
			body:        `{"code":1,"msg":"查询订单号成功！","trade_no":"2016080622555342651","out_trade_no":"order-123","api_trade_no":"4200001","type":"alipay","pid":1001,"addtime":"2016-08-06 22:55:52","endtime":"2016-08-06 22:55:52","name":"VIP会员","money":"1.00","status":1,"param":"","buyer":""}`,
			wantStatus:  payment.ProviderStatusPaid,
			wantTradeNo: "2016080622555342651",
			wantAmount:  1,
		},
		{
			name:        "null top-level fields fall back to nested data",
			body:        `{"code":1,"status":null,"money":null,"trade_no":null,"data":{"status":1,"money":"5.00","trade_no":"data-789"}}`,
			wantStatus:  payment.ProviderStatusPaid,
			wantTradeNo: "data-789",
			wantAmount:  5,
		},
		{
			name:        "matching identity echo with numeric fields is paid",
			body:        `{"code":"1","pid":"pid-1","out_trade_no":"order-123","trade_status":"TRADE_SUCCESS","money":12.34,"trade_no":2026101600001}`,
			wantStatus:  payment.ProviderStatusPaid,
			wantTradeNo: "2026101600001",
			wantAmount:  12.34,
		},
		{
			name:        "missing fields are pending",
			body:        `{}`,
			wantStatus:  payment.ProviderStatusPending,
			wantTradeNo: orderID,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodPost {
					t.Errorf("method = %q, want %q", r.Method, http.MethodPost)
				}
				if r.URL.Path != "/api.php" {
					t.Errorf("path = %q, want /api.php", r.URL.Path)
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("ParseForm: %v", err)
				}
				if got := r.URL.Query(); !reflect.DeepEqual(got, wantQuery) || len(r.PostForm) != 0 {
					t.Errorf("expected query-only params %v, got query %v and form %v", wantQuery, got, r.PostForm)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			provider := newTestEasyPay(t, server.URL)
			resp, err := provider.QueryOrder(context.Background(), orderID)
			if count := calls.Load(); count != 1 {
				t.Fatalf("query attempts = %d, want 1", count)
			}
			if tt.wantErr {
				if err == nil {
					t.Fatalf("QueryOrder succeeded, want error (response=%+v)", resp)
				}
				return
			}
			if err != nil {
				t.Fatalf("QueryOrder returned error: %v", err)
			}
			if resp.Status != tt.wantStatus {
				t.Fatalf("status = %q, want %q (response=%+v)", resp.Status, tt.wantStatus, resp)
			}
			if resp.TradeNo != tt.wantTradeNo {
				t.Fatalf("trade_no = %q, want %q", resp.TradeNo, tt.wantTradeNo)
			}
			if resp.Amount != tt.wantAmount {
				t.Fatalf("amount = %v, want %v", resp.Amount, tt.wantAmount)
			}
		})
	}
}

func TestEasyPayQueryOrderRetriesWithFormOnlyWhenActionMissing(t *testing.T) {
	const orderID = "order-123"
	wantParams := url.Values{
		"act": {"order"}, "pid": {"pid-1"}, "key": {"pkey-1"}, "out_trade_no": {orderID},
	}
	for _, tc := range []struct {
		name     string
		response string
	}{
		{name: "plain text", response: "No Act!"},
		{name: "JSON failure", response: `{"code":-1,"msg":"No Act!"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Errorf("ParseForm: %v", err)
				}
				if calls.Add(1) == 1 {
					_, _ = w.Write([]byte(tc.response))
					return
				}
				if r.Method != http.MethodPost || r.URL.Path != "/api.php" || r.URL.RawQuery != "" || !reflect.DeepEqual(r.PostForm, wantParams) {
					t.Errorf("expected form-only query retry, got %s %s form %v", r.Method, r.URL, r.PostForm)
				}
				_, _ = w.Write([]byte(`{"code":1,"out_trade_no":"order-123","money":"12.34","status":1}`))
			}))
			defer server.Close()

			provider := newTestEasyPay(t, server.URL)
			resp, err := provider.QueryOrder(context.Background(), orderID)
			if err != nil || resp.Status != payment.ProviderStatusPaid {
				t.Fatalf("QueryOrder = %+v, %v, want paid order", resp, err)
			}
			if count := calls.Load(); count != 2 {
				t.Fatalf("query attempts = %d, want 2", count)
			}
		})
	}
}
