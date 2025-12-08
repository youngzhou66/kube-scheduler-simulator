package handler

import (
	"github.com/labstack/echo/v4"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"strconv"

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
	if nodeName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "node name cannot be empty")
	}
	if err := h.service.DeleteNode(ctx, nodeName); err != nil {
		klog.Errorf("failed to delete node: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *KwokClusterHandler) AddDeployment(c echo.Context) error {
	ctx := c.Request().Context()
	var deployment appsv1.Deployment
	if err := c.Bind(&deployment); err != nil {
		klog.Errorf("Failed to parse request body: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid deployment format")
	}
	if err := h.service.AddDeployment(ctx, &deployment); err != nil {
		klog.Errorf("failed to create deployment: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *KwokClusterHandler) DeleteDeployment(c echo.Context) error {
	ctx := c.Request().Context()
	namespace := c.Param("namespace")
	deploymentName := c.Param("name")
	if namespace == "" || deploymentName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "namespace and name cannot be empty")
	}
	if err := h.service.DeleteDeployment(ctx, namespace, deploymentName); err != nil {
		klog.Errorf("failed to delete deployment: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusAccepted)
}

func (h *KwokClusterHandler) AddDeployments(c echo.Context) error {
	ctx := c.Request().Context()
	countParam := c.Param("count")
	count, err := strconv.Atoi(countParam)
	if err != nil || count <= 0 {
		klog.Errorf("Invalid count parameter: %s, error: %+v", countParam, err)
		return echo.NewHTTPError(http.StatusBadRequest, "count must be a positive integer")
	}
	var deployment appsv1.Deployment
	if err := c.Bind(&deployment); err != nil {
		klog.Errorf("Failed to parse request body: %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid deployment format")
	}
	if err := h.service.AddDeployments(ctx, &deployment, count); err != nil {
		klog.Errorf("failed to create deployment: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusAccepted)
}
