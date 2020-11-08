package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ztv1beta1 "github.com/takutakahashi/ztsvc-controller/api/v1beta1"
)

// ZTServiceReconciler reconciles a ZTService object
type ZTServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=zt.takutakahashi.dev,resources=ztservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=zt.takutakahashi.dev,resources=ztservices/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;create;update;patch

func (r *ZTServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("ztservice", req.NamespacedName)
	// add finalizer
	var zts *ztv1beta1.ZTService
	if err := r.Get(ctx, req.NamespacedName, zts); err != nil {
		return ctrl.Result{}, err
	}
	ztsAfter, err := r.reconcile(zts)
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Info("reconciled", "after", ztsAfter)
	return ctrl.Result{}, nil
}

func (r *ZTServiceReconciler) reconcile(zts *ztv1beta1.ZTService) (*ztv1beta1.ZTService, error) {
	var ztsAfter *ztv1beta1.ZTService
	// create service if not
	// inject iptables to daemons
	return ztsAfter, nil
}

func (r *ZTServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ztv1beta1.ZTService{}).
		Complete(r)
}
