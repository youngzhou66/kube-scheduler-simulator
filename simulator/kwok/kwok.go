package kwok

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

const (
	configMapPrefix = "mindx-dl-deviceinfo-"
	kubeSystemNS    = "kube-system"
	consumerCIMKey  = "mx-consumer-cim"
	consumerCIMVal  = "true"
	deviceInfoKey   = "DeviceInfoCfg"
	separateNPUKey  = "ManuallySeparateNPU"
)

// deviceInfo
const deviceInfoConfigTemplate = `{"DeviceInfo":{"DeviceList":{"huawei.com/Ascend910":"Ascend910-0,Ascend910-1,Ascend910-2,Ascend910-3,Ascend910-4,Ascend910-5,Ascend910-6,Ascend910-7","huawei.com/Ascend910-Fault":"[]","huawei.com/Ascend910-NetworkUnhealthy":"","huawei.com/Ascend910-Recovering":"","huawei.com/Ascend910-Unhealthy":""},"UpdateTime":1763713955},"SuperPodID":5,"ServerIndex":0,"RackID":6,"TopoCheck":"OK","CheckCode":"e5cc7a2c30df99b05fb3415484369006515105b0228fc26b097744dceede93ca"}`

type Service struct {
	k8sClient clientset.Interface
}

func NewKwokService(k8sClient clientset.Interface) *Service {
	s := &Service{
		k8sClient: k8sClient,
	}
	return s
}

// AddNode add node
func (s *Service) AddNode(ctx context.Context, node *corev1.Node) error {
	nodeName := s.ensureNodeName(node)
	if err := s.createOrUpdateNode(ctx, node, nodeName); err != nil {
		return fmt.Errorf("failed to create or update node %s: %w", nodeName, err)
	}
	return nil
}

// ensureNodeName ensure node name
func (s *Service) ensureNodeName(node *corev1.Node) string {
	if node.Name != "" {
		return node.Name
	}
	nodeName := fmt.Sprintf("node-%d", rand.Intn(90000)+10000)
	node.Name = nodeName
	return nodeName
}

// createOrUpdateNode create or update node
func (s *Service) createOrUpdateNode(ctx context.Context, node *corev1.Node, nodeName string) error {
	existingNode, err := s.k8sClient.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return s.createNodeAndConfigMap(ctx, node, nodeName)
	}
	if err != nil {
		return err
	}
	return s.updateNode(ctx, node, existingNode)
}

// createNodeWithConfigMap create node and cm
func (s *Service) createNodeAndConfigMap(ctx context.Context, node *corev1.Node, nodeName string) error {
	if _, err := s.k8sClient.CoreV1().Nodes().Create(ctx, node, metav1.CreateOptions{}); err != nil {
		return err
	}
	if err := s.createDeviceInfoConfigMap(ctx, nodeName); err != nil {
		_ = s.k8sClient.CoreV1().Nodes().Delete(ctx, nodeName, metav1.DeleteOptions{})
		return fmt.Errorf("failed to create configmap for node %s: %w", nodeName, err)
	}
	return nil
}

// updateNode update node
func (s *Service) updateNode(ctx context.Context, node *corev1.Node, existingNode *corev1.Node) error {
	node.SetResourceVersion(existingNode.ResourceVersion)
	node.SetUID(existingNode.UID)
	_, err := s.k8sClient.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	return err
}

// createDeviceInfoConfigMap create device info configmap
func (s *Service) createDeviceInfoConfigMap(ctx context.Context, nodeName string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getDeviceInfoConfigMapName(nodeName),
			Namespace: kubeSystemNS,
			Labels: map[string]string{
				consumerCIMKey: consumerCIMVal,
			},
		},
		Data: map[string]string{
			deviceInfoKey:  deviceInfoConfigTemplate,
			separateNPUKey: "",
		},
	}
	_, err := s.k8sClient.CoreV1().ConfigMaps(kubeSystemNS).Create(ctx, configMap, metav1.CreateOptions{})
	return err
}

func getDeviceInfoConfigMapName(nodeName string) string {
	return fmt.Sprintf("%s%s", configMapPrefix, nodeName)
}

// DeleteNode delete node
func (s *Service) DeleteNode(ctx context.Context, nodeName string) error {
	_, err := s.k8sClient.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	err = s.k8sClient.CoreV1().ConfigMaps(kubeSystemNS).Delete(ctx, getDeviceInfoConfigMapName(nodeName), metav1.DeleteOptions{})
	klog.Errorf("failed to delete node config map: %+v", err)
	return s.k8sClient.CoreV1().Nodes().Delete(ctx, nodeName, metav1.DeleteOptions{})
}
