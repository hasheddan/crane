/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane-runtime/pkg/logging"
	"github.com/crossplaneio/crossplane-runtime/pkg/resource"
	gcpcontainerv1beta1 "github.com/crossplaneio/stack-gcp/apis/container/v1beta1"

	"github.com/hasheddan/crane/apis/v1alpha1"
)

var (
	log  = logging.Logger.WithName("controller")
	user = "admin"
)

// JobReconciler reconciles a Job object.
type JobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile reconciles Job objects.
func (r *JobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	log.V(logging.Debug).Info("Reconciling", "controller", "JobReconciler", "request", req)

	if err := r.Client.Create(ctx, &gcpcontainerv1beta1.GKECluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Spec: gcpcontainerv1beta1.GKEClusterSpec{
			ForProvider: gcpcontainerv1beta1.GKEClusterParameters{
				Location: "us-central1-a",
				MasterAuth: &gcpcontainerv1beta1.MasterAuth{
					Username: &user,
				},
			},
			ResourceSpec: runtimev1alpha1.ResourceSpec{
				WriteConnectionSecretToReference: &runtimev1alpha1.SecretReference{
					Name:      "job-secret",
					Namespace: "crossplane-system",
				},
				ProviderReference: &corev1.ObjectReference{
					Name: "gcp-provider",
				},
				ReclaimPolicy: runtimev1alpha1.ReclaimDelete,
			},
		},
	}); err != nil {
		return ctrl.Result{}, errors.Wrap(resource.Ignore(kerrors.IsAlreadyExists, err), "unable to create cluster")
	}

	return ctrl.Result{Requeue: false}, nil
}

// SetupWithManager adds the JobReconciler to the Manager.
func (r *JobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Job{}).
		Complete(r)
}
