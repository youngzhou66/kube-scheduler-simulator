package kwok

import (
	"context"
	"fmt"
	"sync"

	appsv1 "k8s.io/api/apps/v1"
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
const deviceInfoConfigTemplate = `{"DeviceInfo":{"DeviceList":{"huawei.com/Ascend910":"Ascend910-0,Ascend910-1,Ascend910-2,Ascend910-3,Ascend910-4,Ascend910-5,Ascend910-6,Ascend910-7","huawei.com/Ascend910-Fault":"[]","huawei.com/Ascend910-NetworkUnhealthy":"","huawei.com/Ascend910-Recovering":"","huawei.com/Ascend910-Unhealthy":""},"UpdateTime":1763713955},"SuperPodID":%s,"ServerIndex":0,"RackID":%s,"TopoCheck":"OK","CheckCode":"e5cc7a2c30df99b05fb3415484369006515105b0228fc26b097744dceede93ca"}`

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
	// TODO 测试日志
	klog.Errorf("add nodeName: %+v", nodeName)
	if err := s.createOrUpdateNode(ctx, node, nodeName); err != nil {
		return fmt.Errorf("failed to create or update node %s: %w", nodeName, err)
	}
	return nil
}

//func (s *Service) AddNodes(ctx context.Context, node *corev1.Node, count int) error {
//	// todo 要不要优化成多线程 现在好慢...
//	nodeName := node.Name
//	for i := 1; i <= count; i++ {
//		node.Name = fmt.Sprintf("%s-%d", nodeName, i)
//		if err := s.createOrUpdateNode(ctx, node, node.Name); err != nil {
//			return fmt.Errorf("failed to create or update node %s: %w", nodeName, err)
//		}
//	}
//	return nil
//}

func (s *Service) AddNodes(ctx context.Context, node *corev1.Node, count int) error {
	nodeName := node.Name
	errCh := make(chan error, (count+7)/8) // 缓冲通道存储错误
	var wg sync.WaitGroup
	threadCount := (count + 7) / 8
	for t := 0; t < threadCount; t++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			start := threadID*8 + 1
			end := start + 7
			if end > count {
				end = count
			}
			nodeCopy := node.DeepCopy()
			nodeCopy.Annotations["rackID"] = fmt.Sprintf("%d", threadID)
			for i := start; i <= end; i++ {
				nodeCopy.Name = fmt.Sprintf("%s-%d", nodeName, i)
				if err := s.createOrUpdateNode(ctx, nodeCopy, nodeCopy.Name); err != nil {
					errCh <- fmt.Errorf("thread %d: failed to create node %s: %w",
						threadID, nodeCopy.Name, err)
					return
				}
			}
		}(t)
	}

	wg.Wait()
	close(errCh)
	if len(errCh) > 0 {
		return <-errCh // 返回第一个错误
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
	if err == nil {
		return s.updateNodeAndConfigMap(ctx, node, existingNode)
	}
	if !errors.IsNotFound(err) {
		return err
	}
	return s.createNodeAndConfigMap(ctx, node, nodeName)
}

// createNodeWithConfigMap create node and cm
func (s *Service) createNodeAndConfigMap(ctx context.Context, node *corev1.Node, nodeName string) error {
	if _, err := s.k8sClient.CoreV1().Nodes().Create(ctx, node, metav1.CreateOptions{}); err != nil {
		return err
	}
	if err := s.createDeviceInfoConfigMap(ctx, node); err != nil {
		_ = s.k8sClient.CoreV1().Nodes().Delete(ctx, nodeName, metav1.DeleteOptions{})
		return fmt.Errorf("failed to create configmap for node %s: %w", nodeName, err)
	}
	return nil
}

// updateNode update node
func (s *Service) updateNodeAndConfigMap(ctx context.Context, node *corev1.Node, existingNode *corev1.Node) error {
	node.SetResourceVersion(existingNode.ResourceVersion)
	node.SetUID(existingNode.UID)
	_, err := s.k8sClient.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	configMap := s.getNodeCm(node)
	_, err = s.k8sClient.CoreV1().ConfigMaps(kubeSystemNS).Update(ctx, configMap, metav1.UpdateOptions{})
	return err
}

func (s *Service) getNodeCm(node *corev1.Node) *corev1.ConfigMap {
	deviceInfo := fmt.Sprintf(deviceInfoConfigTemplate, node.ObjectMeta.Annotations["superPodID"], node.ObjectMeta.Annotations["rackID"])
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getDeviceInfoConfigMapName(node.Name),
			Namespace: kubeSystemNS,
			Labels: map[string]string{
				consumerCIMKey: consumerCIMVal,
			},
		},
		Data: map[string]string{
			deviceInfoKey:  deviceInfo,
			separateNPUKey: "",
		},
	}
	return configMap
}

// createDeviceInfoConfigMap create device info configmap
func (s *Service) createDeviceInfoConfigMap(ctx context.Context, node *corev1.Node) error {
	configMap := s.getNodeCm(node)
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

// AddDeployment add deployment
func (s *Service) AddDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	if deployment.Name == "" {
		return fmt.Errorf("deployment name cannot be empty")
	}
	if deployment.Namespace == "" {
		deployment.Namespace = "default"
	}
	_, err := s.k8sClient.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	return err
}

// AddDeployments add deployments
func (s *Service) AddDeployments(ctx context.Context, deployment *appsv1.Deployment, count int) error {
	deploymentName := deployment.Name
	if deploymentName == "" {
		return fmt.Errorf("deployment name cannot be empty")
	}
	if deployment.Namespace == "" {
		deployment.Namespace = "default"
	}
	// todo 要不要优化成多线程 现在好慢...
	for i := 1; i <= count; i++ {
		deployment.Name = fmt.Sprintf("%s-%d", deploymentName, i)
		_, err := s.k8sClient.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteDeployment delete deployment
func (s *Service) DeleteDeployment(ctx context.Context, namespace, name string) error {
	_, err := s.k8sClient.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("deployment %s/%s not found: %w", namespace, name, err)
	}
	err = s.k8sClient.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete deployment %s/%s: %w", namespace, name, err)
	}
	return err
}
