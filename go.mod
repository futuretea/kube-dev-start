module kube-dev-start

go 1.14

require (
	github.com/go-logr/logr v0.2.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/prometheus/client_model v0.2.0 // indirect
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.4.0
	sigs.k8s.io/controller-runtime v0.5.0
)

replace (
	k8s.io/api v0.0.0 => k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver v0.0.0 => k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery v0.0.0 => k8s.io/apimachinery v0.18.7-rc.0
	k8s.io/apiserver v0.0.0 => k8s.io/apiserver v0.18.6
	k8s.io/cli-runtime v0.0.0 => k8s.io/cli-runtime v0.18.6
	k8s.io/client-go v0.0.0 => k8s.io/client-go v0.18.6
	k8s.io/cloud-provider v0.0.0 => k8s.io/cloud-provider v0.18.6
	k8s.io/cluster-bootstrap v0.0.0 => k8s.io/cluster-bootstrap v0.18.6
	k8s.io/code-generator v0.0.0 => k8s.io/code-generator v0.18.7-rc.0
	k8s.io/component-base v0.0.0 => k8s.io/component-base v0.18.6
	k8s.io/cri-api v0.0.0 => k8s.io/cri-api v0.18.7-rc.0
	k8s.io/csi-translation-lib v0.0.0 => k8s.io/csi-translation-lib v0.18.6
	k8s.io/kube-aggregator v0.0.0 => k8s.io/kube-aggregator v0.18.6
	k8s.io/kube-controller-manager v0.0.0 => k8s.io/kube-controller-manager v0.18.6
	k8s.io/kube-proxy v0.0.0 => k8s.io/kube-proxy v0.18.6
	k8s.io/kube-scheduler v0.0.0 => k8s.io/kube-scheduler v0.18.6
	k8s.io/kubectl v0.0.0 => k8s.io/kubectl v0.18.6
	k8s.io/kubelet v0.0.0 => k8s.io/kubelet v0.18.6
	k8s.io/legacy-cloud-providers v0.0.0 => k8s.io/legacy-cloud-providers v0.18.6
	k8s.io/metrics v0.0.0 => k8s.io/metrics v0.18.6
	k8s.io/sample-apiserver v0.0.0 => k8s.io/sample-apiserver v0.18.6
	k8s.io/sample-cli-plugin v0.0.0 => k8s.io/sample-cli-plugin v0.18.6
	k8s.io/sample-controller v0.0.0 => k8s.io/sample-controller v0.18.6
)
