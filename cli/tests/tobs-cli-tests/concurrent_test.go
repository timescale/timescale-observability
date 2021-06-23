package tobs_cli_tests

import (
	"testing"
	"time"

	test_utils "github.com/timescale/tobs/cli/tests/test-utils"
)

func changeRelease(t testing.TB) {
	if RELEASE_NAME == "test1" {
		RELEASE_NAME = "test2"
	} else if RELEASE_NAME == "test2" {
		RELEASE_NAME = "test1"
	} else {
		t.Fatalf("Unexpected release name %v", RELEASE_NAME)
	}

	if NAMESPACE == "test1" {
		NAMESPACE = "test2"
	} else if NAMESPACE == "test2" {
		NAMESPACE = "test1"
	} else {
		t.Fatalf("Unexpected namespace %v", RELEASE_NAME)
	}
}

func TestConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent tests")
	}

	// skipping concurrent tests for now
	if true {
		t.Skip("Skipping concurrent tests as it needs changes on multiple tobs installations on the same cluster.")
	}

	u := test_utils.TestUnInstallSpec{
		ReleaseName: RELEASE_NAME,
		Namespace:   NAMESPACE,
		DeleteData:  false,
	}
	u.TestUninstall(t)

	oldname := RELEASE_NAME
	oldspace := NAMESPACE

	RELEASE_NAME = "test1"
	NAMESPACE = "test1"

	i1 := test_utils.TestInstallSpec{
		PathToChart:  PATH_TO_CHART,
		ReleaseName:  RELEASE_NAME,
		Namespace:    NAMESPACE,
		PathToValues: PATH_TO_TEST_VALUES,
		EnableBackUp: false,
		SkipWait:     true,
		OnlySecrets:  false,
	}
	i1.TestInstall(t)
	changeRelease(t)

	i2 := test_utils.TestInstallSpec{
		PathToChart:  PATH_TO_CHART,
		ReleaseName:  RELEASE_NAME,
		Namespace:    NAMESPACE,
		PathToValues: PATH_TO_TEST_VALUES,
		EnableBackUp: false,
		SkipWait:     true,
		OnlySecrets:  false,
	}
	i2.TestInstall(t)

	TestGrafana(t)
	TestMetrics(t)
	TestPortForward(t)
	TestPrometheus(t)
	TestTimescale(t)

	changeRelease(t)
	TestGrafana(t)
	TestMetrics(t)
	TestPortForward(t)
	TestPrometheus(t)
	TestTimescale(t)

	u1 := test_utils.TestUnInstallSpec{
		ReleaseName: RELEASE_NAME,
		Namespace:   NAMESPACE,
		DeleteData:  false,
	}
	u1.TestUninstall(t)
	changeRelease(t)
	u2 := test_utils.TestUnInstallSpec{
		ReleaseName: RELEASE_NAME,
		Namespace:   NAMESPACE,
		DeleteData:  false,
	}
	u2.TestUninstall(t)

	RELEASE_NAME = oldname
	NAMESPACE = oldspace

	i3 := test_utils.TestInstallSpec{
		PathToChart:  PATH_TO_CHART,
		ReleaseName:  RELEASE_NAME,
		Namespace:    NAMESPACE,
		PathToValues: PATH_TO_TEST_VALUES,
		EnableBackUp: false,
		SkipWait:     true,
		OnlySecrets:  false,
	}
	i3.TestInstall(t)

	time.Sleep(10 * time.Second)

	t.Logf("Waiting for pods to initialize...")
	pods, err := kubeClient.K8s.KubeGetAllPods(NAMESPACE, RELEASE_NAME)
	if err != nil {
		t.Logf("Error getting all pods")
		t.Fatal(err)
	}

	for _, pod := range pods {
		err = kubeClient.K8s.KubeWaitOnPod(NAMESPACE, pod.Name)
		if err != nil {
			t.Logf("Error while waiting on pod")
			t.Fatal(err)
		}
	}

	time.Sleep(30 * time.Second)
}
