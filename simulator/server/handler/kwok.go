// handler/kwok_cluster_handler.go
package handler

import (
	"github.com/labstack/echo/v4"
	"k8s.io/klog/v2"

	"net/http"
	"sigs.k8s.io/kube-scheduler-simulator/simulator/server/di"
)

type KwokClusterHandler struct {
	service di.KwokService
}

func NewKwokClusterHandler(s di.KwokService) *KwokClusterHandler {
	return &KwokClusterHandler{service: s}
}

// AddNode 向集群添加节点
func (h *KwokClusterHandler) AddNode(c echo.Context) error {
	ctx := c.Request().Context()
	if err := h.service.AddNode(ctx); err != nil {
		klog.Errorf("failed to reset all resources and schediler configuration: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusAccepted)
}
