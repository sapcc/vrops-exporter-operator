/*
Copyright 2022 SAP SE.

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

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	monitoringv1 "cloud.sap/project/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VropsExporterReconciler reconciles a VropsExporter object
type VropsExporterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.cloud.sap,resources=vropsexporters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.cloud.sap,resources=vropsexporters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.cloud.sap,resources=vropsexporters/finalizers,verbs=update

// Reconcile on VropsExporterReconciler reconciles a VropsExporter object
func (r *VropsExporterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	logger.Info("Fetching vropsExporter resource")
	vropsExporter := &monitoringv1.VropsExporter{}
	if err := r.Client.Get(ctx, req.NamespacedName, vropsExporter); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Error(err, "failed to get vropsExporter resource")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	vropsExporterList := &monitoringv1.VropsExporterList{}
	if err := r.Client.List(context.TODO(), vropsExporterList); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	logger.Info("checking if an existing Deployment exists for this resource")
	deployment := appsv1.Deployment{}
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: vropsExporter.Namespace, Name: vropsExporter.Spec.Name}, &deployment); apierrors.IsNotFound(err) {
		logger.Info("could not find existing Deployment for Vrops Exporter, creating one...")

		for _, exporterType := range vropsExporter.Spec.ExporterTypes {
			deployment = *buildDeployment(*vropsExporter, exporterType)
		}

	}
	return ctrl.Result{}, nil
}

func buildDeployment(vropsExporter monitoringv1.VropsExporter, exporterTyp monitoringv1.ExporterType) *appsv1.Deployment {
	args := []string{
		"-m",
		"/config/collector_config.yaml",
		"-t",
		vropsExporter.Spec.Target,
	}

	for _, collector := range exporterTyp.Collectors {
		args = append(args, "-c", collector)
	}

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            vropsExporter.Spec.Name,
			Namespace:       vropsExporter.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&vropsExporter, monitoringv1.GroupVersion.WithKind("vropsExporter"))},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"vrops-exporter-operator.monitoring.cloud.sap": vropsExporter.Spec.Name,
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"vrops-exporter-operator.monitoring.cloud.sap": vropsExporter.Spec.Name,
					},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{
						{
							Name:  vropsExporter.Spec.Name,
							Image: vropsExporter.Spec.Image,
							Ports: []v1.ContainerPort{
								{Name: "metrics", ContainerPort: vropsExporter.Spec.Port},
							},
							Command: []string{"./exporter.py"},
							Args:    args,
							Env: []v1.EnvVar{
								{Name: "PORT", Value: string(vropsExporter.Spec.Port)},
								{Name: "DEBUG", Value: vropsExporter.Spec.Debug},
								{Name: "INVENTORY", Value: "vrops-inventory"},
							},
							Resources: exporterTyp.Resources,
							VolumeMounts: []v1.VolumeMount{
								{Name: "vrops-config", MountPath: "/config", ReadOnly: true},
							},
						},
					},
				},
			},
		},
	}
	return &deployment
}

// SetupWithManager sets up the controller with the Manager.
func (r *VropsExporterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1.VropsExporter{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
