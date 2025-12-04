// handler/kwok_cluster_handler.go
package handler

import (
	"github.com/labstack/echo/v4"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"

	"net/http"
	"sigs.k8s.io/kube-scheduler-simulator/simulator/server/di"
)

const (
	name = "name"
)

type KwokClusterHandler struct {
	service di.KwokService
}

func NewKwokClusterHandler(s di.KwokService) *KwokClusterHandler {
	return &KwokClusterHandler{service: s}
}

// AddNode add node
func (h *KwokClusterHandler) AddNode(c echo.Context) error {
	ctx := c.Request().Context()
	var node corev1.Node
	if err := c.Bind(&node); err != nil {
		klog.Errorf("Failed to parse request body: %+v", err)
	}
	if err := h.service.AddNode(ctx, &node); err != nil {
		klog.Errorf("failed to add node: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusAccepted)
}

// DeleteNode delete node
func (h *KwokClusterHandler) DeleteNode(c echo.Context) error {
	ctx := c.Request().Context()
	nodeName := c.Param(name)
	if c.Param(nodeName) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "node name cannot be empty")
	}
	if err := h.service.DeleteNode(ctx, nodeName); err != nil {
		klog.Errorf("failed to delete node: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusAccepted)
}
