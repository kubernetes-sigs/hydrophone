package service

const (
	containerImage         = "registry.k8s.io/conformance-amd64:v1.25.0"
	Namespace              = "conformance"
	PodName                = "e2e-conformance-test"
	ClusterRoleBindingName = "conformance-serviceaccount-role"
	ClusterRoleName        = "conformance-serviceaccount"
	ServiceAccountName     = "conformance-serviceaccount"
)
