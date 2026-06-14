package controller

import (
	"context"
	"fmt"

	valkeyv1alpha1 "github.com/rodrigomicrosiga/valkey-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ValkeyInstanceReconciler reconciles a ValkeyInstance object
type ValkeyInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=valkey.cloud104.io,resources=valkeyinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=valkey.cloud104.io,resources=valkeyinstances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=valkey.cloud104.io,resources=valkeyinstances/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services;persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete

func (r *ValkeyInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1. Busca o pedido (ValkeyInstance)
	instance := &valkeyv1alpha1.ValkeyInstance{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Garante o StatefulSet
	if err := r.reconcileStatefulSet(ctx, instance); err != nil {
		logger.Error(err, "Falha ao reconciliar o StatefulSet")
		return ctrl.Result{}, err
	}

	// 3. Garante o Service
	if err := r.reconcileService(ctx, instance); err != nil {
		logger.Error(err, "Falha ao reconciliar o Service")
		return ctrl.Result{}, err
	}

	// 4. Atualiza o Status (As colunas mágicas do K9s)
	instance.Status.Phase = "Ready"
	instance.Status.Endpoint = fmt.Sprintf("%s.%s.svc.cluster.local:%d", instance.Name, instance.Namespace, instance.Spec.Port)
	if err := r.Status().Update(ctx, instance); err != nil {
		logger.Error(err, "Falha ao atualizar o Status")
		return ctrl.Result{}, err
	}

	logger.Info("Valkey reconciliado com sucesso!", "Endpoint", instance.Status.Endpoint)
	return ctrl.Result{}, nil
}

// SetupWithManager registra o controlador e avisa quais recursos ele "Vigia"
func (r *ValkeyInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&valkeyv1alpha1.ValkeyInstance{}).
		Owns(&appsv1.StatefulSet{}). // O operador acorda se alguém mexer no StatefulSet
		Owns(&corev1.Service{}).     // O operador acorda se alguém mexer no Service
		Complete(r)
}
