package controllers

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	funsies "github.com/carbonrelay/kind-controller/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func TestReconcile(t *testing.T) {
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	cfg, err := testEnv.Start()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// +kubebuilder:scaffold:scheme
	err = funsies.AddToScheme(scheme.Scheme)
	require.NoError(t, err)

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	require.NoError(t, err)
	require.NotNil(t, k8sClient)

	reconciler := &EchoReconciler{
		Client: k8sClient,
		Scheme: scheme.Scheme,
		Log:    ctrl.Log,
	}

	testCases := []struct {
		desc    string
		echoRes *funsies.Echo
	}{
		{
			desc: "default",
			echoRes: &funsies.Echo{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testecho",
					Namespace: "default",
				},
				Spec: funsies.EchoSpec{
					Message: "trololololo",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.desc), func(t *testing.T) {
			var err error
			assert.NoError(t, err)

			err = k8sClient.Create(context.Background(), tc.echoRes)
			assert.NoError(t, err)

			_, err = reconciler.Reconcile(ctrl.Request{NamespacedName: types.NamespacedName{Namespace: tc.echoRes.ObjectMeta.Namespace, Name: tc.echoRes.ObjectMeta.Name}})
			assert.NoError(t, err)

			// Wait for resource to be created/reconciled
			echoRes := &funsies.Echo{}
			count := 0
			for len(echoRes.Status.Conditions) != 0 {
				err = k8sClient.Get(context.Background(), client.ObjectKey{Namespace: tc.echoRes.ObjectMeta.Namespace, Name: tc.echoRes.ObjectMeta.Name}, echoRes)
				assert.NoError(t, err)

				if count >= 1000 {
					break
				}

				time.Sleep(10 * time.Millisecond)
				count++
			}

			assert.Equal(t, echoRes.Status.Message, tc.echoRes.Status.Message)
		})
	}
}
