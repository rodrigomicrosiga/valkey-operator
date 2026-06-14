package controller

import (
	"context"

	valkeyv1alpha1 "github.com/rodrigomicrosiga/valkey-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// reconcileStatefulSet garante que o disco e o container do Valkey existam
func (r *ValkeyInstanceReconciler) reconcileStatefulSet(ctx context.Context, instance *valkeyv1alpha1.ValkeyInstance) error {
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}

	// CreateOrUpdate aplica a idempotência
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, sts, func() error {
		// Define a instância como "Dona" deste STS (se apagar a instância, o K8s apaga o STS junto)
		if err := ctrl.SetControllerReference(instance, sts, r.Scheme); err != nil {
			return err
		}

		labels := map[string]string{
			"app.kubernetes.io/name":     "valkey",
			"app.kubernetes.io/instance": instance.Name,
		}

		sts.Labels = labels
		sts.Spec.Replicas = &instance.Spec.Replicas
		sts.Spec.ServiceName = instance.Name
		sts.Spec.Selector = &metav1.LabelSelector{MatchLabels: labels}
		sts.Spec.Template.ObjectMeta.Labels = labels

		storageQty, err := resource.ParseQuantity(instance.Spec.StorageSize)
		if err != nil {
			storageQty = resource.MustParse("2Gi")
		}

		sts.Spec.Template.Spec.Containers = []corev1.Container{
			{
				Name:         "valkey",
				Image:        instance.Spec.Image,
				Ports:        []corev1.ContainerPort{{ContainerPort: instance.Spec.Port, Name: "valkey-port"}},
				VolumeMounts: []corev1.VolumeMount{{Name: "valkey-data", MountPath: "/data"}},
			},
		}

		sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "valkey-data"},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources: corev1.VolumeResourceRequirements{
						Requests: corev1.ResourceList{corev1.ResourceStorage: storageQty},
					},
				},
			},
		}

		return nil
	})

	return err
}

// reconcileService garante que a rede interna do Valkey esteja roteando corretamente
func (r *ValkeyInstanceReconciler) reconcileService(ctx context.Context, instance *valkeyv1alpha1.ValkeyInstance) error {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, svc, func() error {
		if err := ctrl.SetControllerReference(instance, svc, r.Scheme); err != nil {
			return err
		}

		labels := map[string]string{
			"app.kubernetes.io/name":     "valkey",
			"app.kubernetes.io/instance": instance.Name,
		}

		svc.Labels = labels
		svc.Spec.Selector = labels
		svc.Spec.Ports = []corev1.ServicePort{
			{
				Name:       "valkey-port",
				Port:       instance.Spec.Port,
				TargetPort: intstr.FromInt(int(instance.Spec.Port)),
			},
		}

		return nil
	})

	return err
}
