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

	webappv1 "kube-dev-start/api/v1"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	labelKeyGuestbook          = "guestbook"
	labelKeyGuestbookNamespace = "guestbook-ns"
	finalizerNameGuestbook     = "guestbook.finalizer.webapp.my.domain"
)

func hasFinalizer(o metav1.Object, finalizer string) bool {
	f := o.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return true
		}
	}
	return false
}

// GuestbookReconciler reconciles a Guestbook object
type GuestbookReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.my.domain,resources=guestbooks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=guestbooks/status,verbs=get;update;patch

func (r *GuestbookReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("guestbook", req.NamespacedName)

	// your logic here
	guestbook := &webappv1.Guestbook{}
	err := r.Client.Get(ctx, req.NamespacedName, guestbook)
	if err != nil {
		if apierrs.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch guestbook")
		return ctrl.Result{}, err
	}
	if guestbook.ObjectMeta.DeletionTimestamp.IsZero() {
		if !hasFinalizer(&guestbook.ObjectMeta, finalizerNameGuestbook) {
			controllerutil.AddFinalizer(&guestbook.ObjectMeta, finalizerNameGuestbook)
			if err = r.Update(ctx, guestbook); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if hasFinalizer(&guestbook.ObjectMeta, finalizerNameGuestbook) {
			crb := &rbacv1.ClusterRoleBinding{}
			crb.Name = fmt.Sprintf("%s-%s", guestbook.Namespace, guestbook.Name)
			if err = r.Delete(ctx, crb); err != nil && !apierrs.IsNotFound(err) {
				return ctrl.Result{}, err
			}
			cr := &rbacv1.ClusterRole{}
			cr.Name = fmt.Sprintf("%s-%s", guestbook.Namespace, guestbook.Name)
			if err = r.Delete(ctx, cr); err != nil && !apierrs.IsNotFound(err) {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(&guestbook.ObjectMeta, finalizerNameGuestbook)
			if err = r.Update(ctx, guestbook); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	sa := &corev1.ServiceAccount{}
	sa.Namespace = guestbook.Namespace
	sa.Name = guestbook.Name
	if _, err = controllerutil.CreateOrUpdate(ctx, r.Client, sa, r.serviceAccountMutate(sa, guestbook)); err != nil {
		return ctrl.Result{}, err
	}
	cr := &rbacv1.ClusterRole{}
	cr.Name = fmt.Sprintf("%s-%s", guestbook.Namespace, guestbook.Name)
	if _, err = controllerutil.CreateOrUpdate(ctx, r.Client, cr, r.clusterRoleMutate(cr, guestbook)); err != nil {
		return ctrl.Result{}, err
	}
	crb := &rbacv1.ClusterRoleBinding{}
	crb.Name = fmt.Sprintf("%s-%s", guestbook.Namespace, guestbook.Name)
	if _, err = controllerutil.CreateOrUpdate(ctx, r.Client, crb, r.clusterRoleBindingMutate(crb, cr, sa, guestbook)); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *GuestbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Guestbook{}).
		Complete(r)
}

func (r *GuestbookReconciler) relativeResourcesShareLabels(guestbook *webappv1.Guestbook) map[string]string {
	ls := make(map[string]string)
	ls[labelKeyGuestbook] = guestbook.Name
	return ls
}

func (r *GuestbookReconciler) serviceAccountMutate(sa *corev1.ServiceAccount,
	guestbook *webappv1.Guestbook) controllerutil.MutateFn {
	return func() error {
		sa.Labels = r.relativeResourcesShareLabels(guestbook)
		sa.SetOwnerReferences(nil)
		return controllerutil.SetControllerReference(guestbook, sa, r.Scheme)
	}
}

func (r *GuestbookReconciler) clusterRoleMutate(cr *rbacv1.ClusterRole,
	guestbook *webappv1.Guestbook) controllerutil.MutateFn {
	return func() error {
		cr.Labels = r.relativeResourcesShareLabels(guestbook)
		cr.Labels[labelKeyGuestbookNamespace] = guestbook.Namespace
		cr.Rules = []rbacv1.PolicyRule{{
			APIGroups: []string{""},
			Resources: []string{"guestbooks"},
			Verbs:     []string{"get", "list", "watch"},
		}}
		return nil
	}
}

func (r *GuestbookReconciler) clusterRoleBindingMutate(crb *rbacv1.ClusterRoleBinding,
	cr *rbacv1.ClusterRole, sa *corev1.ServiceAccount,
	guestbook *webappv1.Guestbook) controllerutil.MutateFn {
	return func() error {
		crb.Labels = r.relativeResourcesShareLabels(guestbook)
		crb.Labels[labelKeyGuestbookNamespace] = guestbook.Namespace
		crb.RoleRef = rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     cr.Name,
		}
		crb.Subjects = []rbacv1.Subject{{
			Kind:      rbacv1.ServiceAccountKind,
			Name:      sa.Name,
			Namespace: sa.Namespace,
		}}
		return nil
	}
}
