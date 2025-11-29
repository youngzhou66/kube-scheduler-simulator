package kwok

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	clientset "k8s.io/client-go/kubernetes"
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

func (s *Service) AddNode(ctx context.Context) error {
	randomNumber := rand.Intn(90000) + 10000
	nodeName := fmt.Sprintf("node-%d", randomNumber)
	newNode := &corev1.Node{
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
	// TODO
	_, err := s.k8sClient.CoreV1().Nodes().Create(ctx, newNode, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
