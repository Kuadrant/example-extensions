package controller

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	gatewayapiv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// todo(@adam-cattermole): duplicates internal/kuadrant.ErrTargetNotFound
type ErrTargetNotFound struct {
	Kind      string
	TargetRef gatewayapiv1alpha2.LocalPolicyTargetReference
	Err       error
}

func (e ErrTargetNotFound) Error() string {
	if apierrors.IsNotFound(e.Err) {
		return fmt.Sprintf("%s target %s was not found", e.Kind, e.TargetRef.Name)
	}

	return fmt.Sprintf("%s target %s was not found: %s", e.Kind, e.TargetRef.Name, e.Err.Error())
}

func (e ErrTargetNotFound) Reason() gatewayapiv1alpha2.PolicyConditionReason {
	return gatewayapiv1alpha2.PolicyReasonTargetNotFound
}

func NewErrTargetNotFound(kind string, targetRef gatewayapiv1alpha2.LocalPolicyTargetReference, err error) ErrTargetNotFound {
	return ErrTargetNotFound{
		Kind:      kind,
		TargetRef: targetRef,
		Err:       err,
	}
}
