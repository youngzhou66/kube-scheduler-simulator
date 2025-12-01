package kwok

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

type Service struct {
	k8sClient clientset.Interface
}

func NewKwokService(k8sClient clientset.Interface) *Service {
	s := &Service{
		k8sClient: k8sClient,
	}
	return s
}

func (s *Service) AddNode(ctx context.Context, node *corev1.Node) error {
	nodeName := node.Name
	if node.Name == "" {

		randomNumber := rand.Intn(90000) + 10000
		nodeName = fmt.Sprintf("node-%d", randomNumber)
		node = &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: nodeName,
			},
			Status: corev1.NodeStatus{
				Capacity: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("4"),
					corev1.ResourceMemory: resource.MustParse("16Gi"),
					corev1.ResourcePods:   resource.MustParse("110"),
				},
				Allocatable: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse("4"),
					corev1.ResourceMemory: resource.MustParse("16Gi"),
					corev1.ResourcePods:   resource.MustParse("110"),
				},
				Conditions: []corev1.NodeCondition{
					{
						Type:   corev1.NodeReady,
						Status: corev1.ConditionTrue,
					},
				},
			},
		}
	} else {
		klog.Errorf("node info : %+v", node)
	}
	// TODO
	_, err := s.k8sClient.CoreV1().Nodes().Create(ctx, node, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	err = s.createDynamicDeviceInfoConfigMap(ctx, nodeName)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) createDynamicDeviceInfoConfigMap(ctx context.Context, nodeName string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("mindx-dl-deviceinfo-%s", nodeName),
			Namespace: "kube-system",
			Labels: map[string]string{
				"mx-consumer-cim": "true",
			},
		},
		Data: map[string]string{
			"DeviceInfoCfg":       `{"DeviceInfo":{"DeviceList":{"huawei.com/Ascend910":"Ascend910-0,Ascend910-1,Ascend910-2,Ascend910-3,Ascend910-4,Ascend910-5,Ascend910-6,Ascend910-7","huawei.com/Ascend910-Fault":"[]","huawei.com/Ascend910-NetworkUnhealthy":"","huawei.com/Ascend910-Recovering":"","huawei.com/Ascend910-Unhealthy":""},"UpdateTime":1763713955},"SuperPodID":5,"ServerIndex":0,"RackID":6,"TopoCheck":"OK","CheckCode":"e5cc7a2c30df99b05fb3415484369006515105b0228fc26b097744dceede93ca"}`,
			"ManuallySeparateNPU": "",
		},
	}

	_, err := s.k8sClient.CoreV1().ConfigMaps("kube-system").Create(ctx, configMap, metav1.CreateOptions{})
	return err
}
