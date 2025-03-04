package ran_du_system_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-goinfra/pkg/namespace"
	"github.com/openshift-kni/eco-systemtests/tests/internal/polarion"
	"github.com/openshift-kni/eco-systemtests/tests/internal/await"
	"github.com/openshift-kni/eco-systemtests/tests/internal/shell"
	. "github.com/openshift-kni/eco-systemtests/tests/ran-du/internal/randuinittools"
	"github.com/openshift-kni/eco-systemtests/tests/ran-du/internal/randuparams"
	"github.com/openshift-kni/eco-systemtests/tests/ran-du/internal/randutestworkload"
)

var _ = Describe(
	"LaunchWorkload",
	Ordered,
	ContinueOnFailure,
	Label(randuparams.LabelLaunchWorkloadTestCases), func() {
		BeforeAll(func() {
			By("Preparing workload")

			if namespace.NewBuilder(APIClient, RanDuTestConfig.TestWorkload.Namespace).Exists() {
				err := randutestworkload.CleanNameSpace(randuparams.DefaultTimeout, RanDuTestConfig.TestWorkload.Namespace)
				Expect(err).ToNot(HaveOccurred(), "Failed to clean workload test namespace objects")
			}

			if RanDuTestConfig.TestWorkload.CreateMethod == "shell" {
				By("Launching workload using shell method")
				_, err := shell.ExecuteCmd(RanDuTestConfig.TestWorkload.CreateShellCmd)
				Expect(err).ToNot(HaveOccurred(), "Failed to launch workload")
			}

			By("Waiting for deployment replicas to become ready")
			_, err := await.WaitUntilAllDeploymentsReady(APIClient, RanDuTestConfig.TestWorkload.Namespace,
				randuparams.DefaultTimeout)
			Expect(err).ToNot(HaveOccurred(), "error while waiting for deployment to become ready")

			By("Waiting for statefulset replicas to become ready")
			_, err = await.WaitUntilAllStatefulSetsReady(APIClient, RanDuTestConfig.TestWorkload.Namespace,
				randuparams.DefaultTimeout)
			Expect(err).ToNot(HaveOccurred(), "error while waiting for statefulsets to become ready")

		})
		It("Assert all pods are ready", polarion.ID("55465"), Label("launch-workload"), func() {
			_, err := await.WaitUntilAllPodsReady(APIClient, RanDuTestConfig.TestWorkload.Namespace, 10*time.Second)
			Expect(err).ToNot(HaveOccurred(), "pod not ready: %s", err)

		})
		AfterAll(func() {
			By("Cleaning up test workload resources")
			err := randutestworkload.CleanNameSpace(randuparams.DefaultTimeout, RanDuTestConfig.TestWorkload.Namespace)
			Expect(err).ToNot(HaveOccurred(), "Failed to clean workload test namespace objects")
		})
	})
