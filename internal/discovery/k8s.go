package discovery

import (
	"fmt"

	"github.com/newrelic/nri-flex/internal/load"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func getK8Labels(podName string, podNamespace string) map[string]string {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		load.Logrus.Debug("k8s: unable to set InClusterConfig config :  ", err)
		return nil
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		load.Logrus.Debug("k8s: unable to create clientset  ", err)
		return nil
	}

	//Get the POD and return Labels
	//listOptions := metav1.ListOptions{}
	//pods, err := clientset.CoreV1().Pods(podNamespace).List(listOptions)

	getOptions := metav1.GetOptions{}

	p, err := clientset.CoreV1().Pods(podNamespace).Get(podName, getOptions)
	if err != nil {
		load.Logrus.Debug("k8s: error getting pod: " + err.Error())
		return nil
	}
	load.Logrus.Debug(fmt.Sprintf("k8s: namespace: %v, pods: %v, labels : %v ", podNamespace, podName, p.Labels))
	return p.Labels
}
