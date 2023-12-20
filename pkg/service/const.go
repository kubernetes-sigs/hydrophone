package service

const (
	containerImage         = "registry.k8s.io/conformance:v1.28.0"
	Namespace              = "conformance"
	PodName                = "e2e-conformance-test"
	ClusterRoleBindingName = "conformance-serviceaccount-role"
	ClusterRoleName        = "conformance-serviceaccount"
	ServiceAccountName     = "conformance-serviceaccount"
)
