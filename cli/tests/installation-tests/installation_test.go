package installation_tests

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	test_utils "github.com/timescale/tobs/cli/tests/test-utils"
)

func testHelmDeleteData(t testing.TB, name, namespace string) {
	cmds := []string{"helm", "delete-data"}
	if name != "" {
		cmds = append(cmds, "-n", name)
	} else {
		cmds = append(cmds, "-n", RELEASE_NAME)
	}
	if namespace != "" {
		cmds = append(cmds, "--namespace", namespace)
	} else {
		cmds = append(cmds, "--namespace", NAMESPACE)
	}

	t.Logf("Running '%v'", "tobs "+strings.Join(cmds, " "))
	deletedata := exec.Command(PATH_TO_TOBS, cmds...)

	out, err := deletedata.CombinedOutput()
	if err != nil {
		t.Logf(string(out))
		t.Fatal(err)
	}

	pvcs, err := kubeClient.K8s.KubeGetPVCNames("default", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if len(pvcs) != 0 {
		t.Fatal("PVC remaining")
	}
}

func testHelmShowValues(t testing.TB) {
	var showvalues *exec.Cmd

	t.Logf("Running 'tobs helm show-values'")

	showvalues = exec.Command(PATH_TO_TOBS, "helm", "show-values", "-c", PATH_TO_CHART)
	out, err := showvalues.CombinedOutput()
	if err != nil {
		t.Logf(string(out))
		t.Fatal(err)
	}
}

func TestInstallation(t *testing.T) {

	// abc Install
	abcInstall := test_utils.TestInstallSpec{
		PathToChart:  PATH_TO_CHART,
		ReleaseName:  "abc",
		Namespace:    NAMESPACE,
		PathToValues: PATH_TO_TEST_VALUES,
		EnableBackUp: false,
		SkipWait:     true,
		OnlySecrets:  false,
	}
	abcInstall.TestInstall(t)

	// abc Uninstall
	abcUninstall := test_utils.TestUnInstallSpec{
		ReleaseName:  "abc",
		Namespace:    NAMESPACE,
		DeleteData:   false,
	}
	abcUninstall.TestHelmUninstall(t)

	// def helm cmd install
	defInstall := test_utils.TestInstallSpec{
		PathToChart: PATH_TO_CHART,
		ReleaseName:  "def",
		Namespace:    NAMESPACE,
		PathToValues: PATH_TO_TEST_VALUES,
		EnableBackUp: false,
		SkipWait:     true,
		OnlySecrets:  false,
	}
	defInstall.TestHelmInstall(t)

	// def uninstall
	defUninstall := test_utils.TestUnInstallSpec{
		ReleaseName: "def",
		Namespace:   NAMESPACE,
		DeleteData:  false,
	}
	defUninstall.TestUninstall(t)
	testHelmDeleteData(t, "def", "")

	// f1 install
	f1Install := test_utils.TestInstallSpec{
		PathToChart:  PATH_TO_CHART,
		ReleaseName:  "f1",
		Namespace:    NAMESPACE,
		PathToValues: "./../testdata/f1.yaml",
		EnableBackUp: false,
		SkipWait:     true,
		OnlySecrets:  false,
	}
	f1Install.TestInstall(t)

	// f1 uninstall
	f1Uninstall := test_utils.TestUnInstallSpec{
		ReleaseName: "f1",
		Namespace:   NAMESPACE,
		DeleteData:  false,
	}
	f1Uninstall.TestHelmUninstall(t)

	// f2 install
	f2Install := test_utils.TestInstallSpec{
		PathToChart:  PATH_TO_CHART,
		ReleaseName:  "f2",
		Namespace:    NAMESPACE,
		PathToValues: "./../testdata/f2.yaml",
		EnableBackUp: false,
		SkipWait:     true,
		OnlySecrets:  false,
	}
	f2Install.TestHelmInstall(t)
	// f2 uninstall
	f2Uninstall := test_utils.TestUnInstallSpec{
		ReleaseName: "f2",
		Namespace:   NAMESPACE,
		DeleteData:  false,
	}
	f2Uninstall.TestUninstall(t)

	// install --only-secrets
	f5Install := test_utils.TestInstallSpec{
		PathToChart:  PATH_TO_CHART,
		ReleaseName:  "f5",
		Namespace:    "secrets",
		PathToValues: PATH_TO_TEST_VALUES,
		EnableBackUp: false,
		SkipWait:     true,
		OnlySecrets:  true,
	}
	f5Install.TestInstall(t)
	pods, err := kubeClient.K8s.KubeGetAllPods("secrets", "f5")
	if err != nil {
		t.Log("failed to get all tobs pods")
		t.Fatal(err)
	}
	if len(pods) != 0 {
		t.Fatal("failed to install tobs with --only-secrets. We see other pods by tobs install")
	}
	err = kubeClient.DeleteNamespace("secrets")
	if err != nil {
		t.Fatal(err)
	}

	// kubectl get pods -A
	test_utils.ShowAllPods(t)

	// This installation is used to run all tests in tobs-cli-tests
	dInstall := test_utils.TestInstallSpec{
		PathToChart:  PATH_TO_CHART,
		ReleaseName:  RELEASE_NAME,
		Namespace:    NAMESPACE,
		PathToValues: PATH_TO_TEST_VALUES,
		EnableBackUp: false,
		SkipWait:     false,
		OnlySecrets:  false,
	}
	dInstall.TestInstall(t)

	time.Sleep(3 * time.Minute)

	// test show values
	testHelmShowValues(t)

	// kubectl get pods -A
	test_utils.ShowAllPods(t)

	t.Logf("Waiting for pods to initialize...")
	pods, err = kubeClient.K8s.KubeGetAllPods(NAMESPACE, RELEASE_NAME)
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
