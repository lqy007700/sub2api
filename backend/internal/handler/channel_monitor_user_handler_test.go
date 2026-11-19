//go:build unit

package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type userMonitorTimelineRepo struct {
	service.ChannelMonitorRepository
	historyCalls int
}

func (*userMonitorTimelineRepo) ListEnabled(context.Context) ([]*service.ChannelMonitor, error) {
	return []*service.ChannelMonitor{{ID: 1, PrimaryModel: "model"}}, nil
}

func (*userMonitorTimelineRepo) ListLatestForMonitorIDs(context.Context, []int64) (map[int64][]*service.ChannelMonitorLatest, error) {
	return map[int64][]*service.ChannelMonitorLatest{
		1: {{Model: "model", Status: service.MonitorStatusOperational}},
	}, nil
}

func (*userMonitorTimelineRepo) ComputeAvailabilityForMonitors(context.Context, []int64, int) (map[int64][]*service.ChannelMonitorAvailability, error) {
	return nil, nil
}

func (r *userMonitorTimelineRepo) ListRecentHistoryForMonitors(context.Context, []int64, map[int64]string, int) (map[int64][]*service.ChannelMonitorHistoryEntry, error) {
	r.historyCalls++
	return map[int64][]*service.ChannelMonitorHistoryEntry{
		1: {{Status: service.MonitorStatusFailed}},
	}, nil
}

func TestChannelMonitorUserListCanSkipRecentHistory(t *testing.T) {
	repo := &userMonitorTimelineRepo{}
	handler := NewChannelMonitorUserHandler(service.NewChannelMonitorService(repo, nil), nil)
	router := gin.New()
	router.GET("/channel-monitors", handler.List)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/channel-monitors?include_timeline=false", nil))
	require.Equal(t, http.StatusOK, response.Code)
	require.Zero(t, repo.historyCalls)
	require.Contains(t, response.Body.String(), `"primary_status":"operational"`)
	require.Contains(t, response.Body.String(), `"timeline":[]`)

	response = httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/channel-monitors", nil))
	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, 1, repo.historyCalls)
	require.Contains(t, response.Body.String(), `"timeline":[{"status":"failed"`)
}
