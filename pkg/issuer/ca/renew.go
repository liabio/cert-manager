package ca

import (
	"context"

	"github.com/jetstack-experimental/cert-manager/pkg/apis/certmanager/v1alpha1"
	"github.com/jetstack-experimental/cert-manager/pkg/util/kube"
	"github.com/jetstack-experimental/cert-manager/pkg/util/pki"
)

const (
	errorRenewCert = "ErrRenewCert"

	successCertRenewed = "CertIssueSuccess"

	messageErrorRenewCert = "Error issuing TLS certificate: "

	messageCertRenewed = "Certificate issued successfully"
)

func (c *CA) Renew(ctx context.Context, crt *v1alpha1.Certificate) (v1alpha1.CertificateStatus, []byte, []byte, error) {
	update := crt.DeepCopy()

	signeeKey, err := kube.SecretTLSKey(c.secretsLister, crt.Namespace, crt.Spec.SecretName)

	if err != nil {
		s := messageErrorGetCertKeyPair + err.Error()
		update.UpdateStatusCondition(v1alpha1.CertificateConditionReady, v1alpha1.ConditionFalse, errorGetCertKeyPair, s)
		return update.Status, nil, nil, err
	}

	certPem, err := c.obtainCertificate(crt, signeeKey)

	if err != nil {
		s := messageErrorRenewCert + err.Error()
		update.UpdateStatusCondition(v1alpha1.CertificateConditionReady, v1alpha1.ConditionFalse, errorRenewCert, s)
		return update.Status, nil, nil, err
	}

	update.UpdateStatusCondition(v1alpha1.CertificateConditionReady, v1alpha1.ConditionTrue, successCertRenewed, messageCertRenewed)

	return update.Status, pki.EncodePKCS1PrivateKey(signeeKey), certPem, nil
}
