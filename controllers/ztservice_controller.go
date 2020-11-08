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

func (r *ZTServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("ztservice", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *ZTServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ztv1beta1.ZTService{}).
		Complete(r)
}
