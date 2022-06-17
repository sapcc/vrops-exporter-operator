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
	"strconv"

	// "fmt"

	monitoringv1 "cloud.sap/project/api/v1" // github.com/sapcc/vrops-exporter-operator/api/v1
	v1 "cloud.sap/project/api/v1"

	// "github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	// netv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	// "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// "sigs.k8s.io/controller-runtime/pkg/controller"
	// "sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	// "sigs.k8s.io/controller-runtime/pkg/source"
)

// VropsExporterReconciler reconciles a VropsExporter object
type VropsExporterReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.cloud.sap;apps,resources=vropsexporters;deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.cloud.sap,resources=vropsexporters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.cloud.sap,resources=vropsexporters/finalizers,verbs=*

// Reconcile on VropsExporterReconciler reconciles a VropsExporter object
func (r *VropsExporterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Fetching vropsExporter resource")

	vropsExporter := &monitoringv1.VropsExporter{}
	if err := r.Client.Get(ctx, req.NamespacedName, vropsExporter); err != nil {
		// Object not found, return.  Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		if apierrors.IsNotFound(err) {
			logger.Error(err, "failed to get vropsExporter resource")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	vropsExporterList := &monitoringv1.VropsExporterList{}
	if err := r.Client.List(context.TODO(), vropsExporterList); err != nil {
		// Object not found, return.  Created objects are automatically garbage collected.
		// For additional cleanup logic use finalizers.
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	var vropsExporterInstances []v1.VropsExporter
	for _, ve := range vropsExporterList.Items {
		if ve.Name != vropsExporter.Name {
			vropsExporterInstances = append(vropsExporterInstances, ve)
		}
	}

	for _, exporterType := range vropsExporter.Spec.ExporterTypes {
		vropsExporterFullName := vropsExporter.Spec.Name + "-" + exporterType.Name

		logger.Info("checking if an existing Deployment exists for this resource")
		deployment := appsv1.Deployment{}

		if err := r.Client.Get(ctx, client.ObjectKey{Namespace: vropsExporter.Namespace, Name: vropsExporterFullName}, &deployment); apierrors.IsNotFound(err) {

			logger.Info("could not find existing Deployment for ", vropsExporterFullName, " creating one...")

			deployment = *buildDeployment(*vropsExporter, exporterType, vropsExporterFullName)
			if err := r.Client.Create(ctx, &deployment); err != nil {
				logger.Error(err, "failed to create Deployment resource")
				return ctrl.Result{}, err
			}

			if err != nil {
				logger.Error(err, "failed to get Deployment for vropsExporter resource")
				return ctrl.Result{}, err
			}
			logger.Info("created Deployment resource for VropsExporter")
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

func buildDeployment(vropsExporter monitoringv1.VropsExporter, exporterType monitoringv1.ExporterType, vropsExporterFullName string) *appsv1.Deployment {

	args := []string{
		"-m",
		"/config/collector_config.yaml",
		"-t",
		vropsExporter.Spec.Target,
	}

	for _, collector := range exporterType.Collectors {
		args = append(args, "-c", collector)
	}

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            vropsExporterFullName,
			Namespace:       vropsExporter.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(&vropsExporter, monitoringv1.GroupVersion.WithKind("vropsExporter"))},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"vrops-exporter-operator.monitoring.cloud.sap": vropsExporterFullName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"vrops-exporter-operator.monitoring.cloud.sap": vropsExporterFullName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  vropsExporterFullName,
							Image: vropsExporter.Spec.Image,
							Ports: []corev1.ContainerPort{
								{Name: "metrics", ContainerPort: vropsExporter.Spec.Port},
							},
							Command: []string{"./exporter.py"},
							Args:    args,
							Env: []corev1.EnvVar{
								{
									Name:  "PORT",
									Value: strconv.Itoa(int(vropsExporter.Spec.Port))},
								{
									Name:  "DEBUG",
									Value: vropsExporter.Spec.Debug},
								{
									Name:  "INVENTORY",
									Value: "vrops-inventory"},
							},
							Resources: exporterType.Resources,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "vrops-config",
									MountPath: "/config",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "vrops-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "vrops-exporter-collector-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return &deployment
}

var (
	deploymentOwnerKey = ".metadata.controller"
)

// SetupWithManager sets up the controller with the Manager.
func (r *VropsExporterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1.VropsExporter{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
