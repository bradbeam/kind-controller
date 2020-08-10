/*


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
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	funsies "github.com/carbonrelay/kind-controller/api/v1alpha1"
)

// EchoReconciler reconciles a Echo object
type EchoReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=funsies.example.squiggly.dev,resources=echoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=funsies.example.squiggly.dev,resources=echoes/status,verbs=get;update;patch

func (r *EchoReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info(fmt.Sprintf("reconcile loop starting for %v", req.NamespacedName))
	defer r.Log.Info(fmt.Sprintf("reconcile loop finished for %v", req.NamespacedName))

	_ = r.Log.WithValues("echo", req.NamespacedName)

	ctx := context.Background()

	// Look up object
	echoRes := &funsies.Echo{}

	// Identify the resource
	if err := r.Get(ctx, req.NamespacedName, echoRes); err != nil {
		return ctrl.Result{}, err
	}

	if echoRes.ObjectMeta.Generation == echoRes.Status.ObservedGeneration {
		return ctrl.Result{}, nil
	}

	// Update status with echo.Message
	echoRes.Status.ObservedGeneration = echoRes.ObjectMeta.Generation
	echoRes.Status.Message = echoRes.Spec.Message
	echoRes.Status.Conditions = append(echoRes.Status.Conditions, funsies.EchoCondition{Type: "completed", Status: v1.ConditionTrue})

	if err := r.Status().Update(ctx, echoRes); err != nil {
		r.Log.Error(err, "failed to update echo status")
		return ctrl.Result{}, err
	}

	if err := r.Update(ctx, echoRes); err != nil {
		r.Log.Error(err, "failed to update echo status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *EchoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&funsies.Echo{}).
		Complete(r)
}
