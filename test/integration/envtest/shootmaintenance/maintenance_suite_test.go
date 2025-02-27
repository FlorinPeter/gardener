// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shoot_maintenance_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	logzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/gardener/gardener/pkg/client/kubernetes"
	controllermanagerfeatures "github.com/gardener/gardener/pkg/controllermanager/features"
	"github.com/gardener/gardener/pkg/envtest"
	log "github.com/gardener/gardener/pkg/logger"
	. "github.com/gardener/gardener/pkg/utils/test/matchers"
	"github.com/gardener/gardener/test/framework"
)

func TestShootMaintenance(t *testing.T) {
	controllermanagerfeatures.RegisterFeatureGates()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoot Maintenance Controller Integration Test Suite")
}

var (
	logger *logrus.Logger

	ctx        = context.Background()
	testEnv    *envtest.GardenerTestEnvironment
	restConfig *rest.Config
	mgrCancel  context.CancelFunc
	testClient client.Client
)

var _ = BeforeSuite(func() {
	logf.SetLogger(logzap.New(logzap.UseDevMode(true), logzap.WriteTo(GinkgoWriter)))

	By("starting test environment")
	testEnv = &envtest.GardenerTestEnvironment{
		GardenerAPIServer: &envtest.GardenerAPIServer{
			Args: []string{
				"--disable-admission-plugins=ResourceReferenceManager,ExtensionValidator,ShootQuotaValidator,ShootValidator,ShootTolerationRestriction",
				"--feature-gates=WorkerPoolKubernetesVersion=true",
			},
		},
	}
	var err error
	restConfig, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())

	testClient, err = client.New(restConfig, client.Options{Scheme: kubernetes.GardenScheme})
	Expect(err).ToNot(HaveOccurred())

	logger = log.AddWriter(log.NewLogger("", ""), GinkgoWriter)

	By("setup manager")
	mgr, err := manager.New(restConfig, manager.Options{
		Scheme:             kubernetes.GardenScheme,
		MetricsBindAddress: "0",
	})
	Expect(err).ToNot(HaveOccurred())

	Expect(addShootMaintenanceControllerToManager(mgr)).To(Succeed())

	var mgrContext context.Context
	mgrContext, mgrCancel = context.WithCancel(ctx)

	By("start manager")
	go func() {
		Expect(mgr.Start(mgrContext)).To(Succeed())
	}()

	By("create shoot namespace")
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "garden-dev"},
	}
	Expect(testClient.Create(ctx, namespace)).To(Or(Succeed(), BeAlreadyExistsError()))
})

var _ = AfterSuite(func() {
	By("stopping manager")
	mgrCancel()

	By("running cleanup actions")
	framework.RunCleanupActions()

	By("stopping test environment")
	Expect(testEnv.Stop()).To(Succeed())
})
