package e2e

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/exec"
)

func Test_Upgrade(t *testing.T) {
	if e2eSuite == nil {
		t.Skip("e2eSuite not defined")
		return
	}

	defer func() {
		err := e2eSuite.kubeclient.Rbac().ClusterRoleBindings().Delete("test-username-binding", nil)
		if err != nil {
			t.Errorf("failed to delete test cluster rolebinding: %s", err)
		}
	}()

	_, err := e2eSuite.kubeclient.Rbac().ClusterRoleBindings().Create(
		&rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-username-binding",
				Namespace: namespaceTokenTest,
			},
			Subjects: []rbacv1.Subject{
				{
					Name: "test-username",
					Kind: "User",
				},
			},
			RoleRef: rbacv1.RoleRef{
				Name: "cluster-admin",
				Kind: "ClusterRole",
			},
		})
	if err != nil {
		t.Fatal(err)
	}

	restConfig, err := clientcmd.BuildConfigFromFlags("", "/home/josh/.kube/kind-config-kube-oidc-proxy-e2e")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", restConfig)
	restConfig.GroupVersion = &corev1.SchemeGroupVersion
	restConfig.NegotiatedSerializer = scheme.Codecs

	////coreclient := kubeclient.Core()
	//coreclient := kubeclient.Core()

	//pod, err := coreclient.Pods("kube-system").Get("kube-proxy-q6c48", metav1.GetOptions{})
	//if err != nil {
	//	t.Fatal(err)
	//}

	pod, err := e2eSuite.kubeclient.Core().Pods("kube-system").Get("kube-proxy-qz9ht", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	ioStreams := genericclioptions.IOStreams{In: nil, Out: execOut, ErrOut: execErr}

	kubeConfigFlags := genericclioptions.NewConfigFlags()
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)

	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	os.Args = []string{"exec", "kube-proxy-gke-kube-oidc-proxy-default-pool-bcef1072-3tc8", "-it", "/bin/sh"}
	cmd := exec.NewCmdExec(f, ioStreams)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s\n", execOut.String())
	fmt.Printf("%s\n", execErr.String())

	//	kubeclient, err := kubernetes.NewForConfig(restConfig)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	options := &exec.ExecOptions{
	//		StreamOptions: exec.StreamOptions{
	//			IOStreams: genericclioptions.IOStreams{
	//				Out:    &execOut,
	//				ErrOut: &execErr,
	//			},
	//
	//			Namespace:     pod.Namespace,
	//			ContainerName: pod.Spec.Containers[0].Name,
	//			PodName:       pod.Name,
	//		},
	//
	//		Command:  []string{"echo", "hello", "world"},
	//		Executor: &exec.DefaultRemoteExecutor{},
	//		Config:   restConfig,
	//
	//		PodClient: kubeclient.Core(),
	//	}
	//	if err := options.Validate(); err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	if err := options.Run(); err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	fmt.Printf("%s\n", execOut.String())
	//	fmt.Printf("%s\n", execErr.String())

	//	err = e2eSuite.setAuthHeader(e2eSuite.validToken())
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	//URL := &url.URL{
	//	//	Host:   fmt.Sprintf("https://127.0.0.1:%s", e2eSuite.proxyPort),
	//	//	Scheme: "https",
	//	//}
	//	restconfig := new(rest.Config)
	//	if err := rest.SetKubernetesDefaults(restconfig); err != nil {
	//		t.Fatal(err)
	//	}
	//	restconfig.Host = fmt.Sprintf("https://127.0.0.1:%s", e2eSuite.proxyPort)
	//	//restconfig.Host = fmt.Sprintf("https://127.0.0.1:%s", e2eSuite.proxyPort)
	//	//restconfig.WrapTransport(e2eSuite.proxyClient.Transport)
	//	//restconfig := &rest.Config{
	//	//	Host: fmt.Sprintf("https://127.0.0.1:%s", e2eSuite.proxyPort),
	//	//	//Host:      fmt.Sprintf("https://localhost:%s", "35499"),
	//	//	Transport: e2eSuite.proxyClient.Transport,
	//	//	//TLSClientConfig: rest.TLSClientConfig{
	//	//	//	CAData: e2eSuite.kubeRestConfig.CAData,
	//	//	//},
	//	//}
	//	//restconfig2, err := rest.RESTClientFor(restconfig)
	//	//if err != nil {
	//	//	t.Fatal(err)
	//	//}
	//
	//	restconfig.CAData = append(append(e2eSuite.kubeRestConfig.CAData, '\n'),
	//		e2eSuite.proxyCert...)
	//	restconfig.BearerToken = e2eSuite.wrappedRT.token
	//	kubeclient, err := kubernetes.NewForConfig(restconfig)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	////coreclient := kubeclient.Core()
	//	//coreclient := kubeclient.Core()
	//
	//	//pod, err := coreclient.Pods("kube-system").Get("kube-proxy-q6c48", metav1.GetOptions{})
	//	//if err != nil {
	//	//	t.Fatal(err)
	//	//}
	//
	//	pod, err := e2eSuite.kubeclient.Core().Pods("kube-system").Get("kube-proxy-qz9ht", metav1.GetOptions{})
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	container := pod.Spec.Containers[0].Name
	//	fmt.Printf(">>%s\n", container)
	//
	//	//req := coreclient.RESTClient().Post().
	//	req := kubeclient.RESTClient().Post().
	//		//req := e2eSuite.kubeclient.RESTClient().Post().
	//		Resource("pods").
	//		Name(pod.Name).
	//		Namespace(pod.Namespace).
	//		SubResource("exec")
	//	attachOptions := &corev1.PodExecOptions{
	//		Stderr: true,
	//		//Stdin:   true,
	//		Stdout:  true,
	//		Command: []string{"echo hello world"},
	//
	//		//Container: container,
	//
	//		//Container: "kube-proxy",
	//		Container: container,
	//		//TTY:       true,
	//	}
	//
	//	//url := req.URL()
	//	//url.Host = e2eSuite.kubeRestConfig.Host
	//	//url.Host = "localhost:35499"
	//
	//	req.VersionedParams(attachOptions, scheme.ParameterCodec)
	//
	//	//fmt.Printf(">>%s\n", restconfig.CAData)
	//	e2eSuite.setAuthHeader(e2eSuite.validToken())
	//	url := req.URL()
	//	url.Path = "/api/v1" + url.Path
	//	url.RawQuery = "command=echo&command=hello&command=world&container=kube-proxy&container=kube-proxy&stderr=true&stdout=true"
	//
	//	//fmt.Printf("%s\n", req.URL())
	//	//URL, err := url.Parse(fmt.Sprintf("https://127.0.0.1:%s/api/v1/namespaces/kube-system/pods/kube-proxy-qz9ht/exec?command=echo&command=hello&command=world&container=kube-proxy&container=kube-proxy&stderr=true&stdout=true", e2eSuite.proxyPort))
	//	//if err != nil {
	//	//	t.Fatal(err)
	//	//}
	//	//URL := &url.URL{
	//	//	Scheme: "https",
	//	//	Path:   "",
	//	//}
	//	//fmt.Printf("%s\n", URL.String())
	//	fmt.Printf("%s\n", url)
	//	//https: //localhost:32789
	//
	//	//target := fmt.Sprintf("https://127.0.0.1:%s/namespaces/kube-system/pods/kube-proxy-q6c48/exec?timeout=32s", e2eSuite.proxyPort)
	//
	//	//if err := e2eSuite.setAuthHeader(e2eSuite.validToken()); err != nil {
	//	//	t.Error(err)
	//	//	return
	//	//}
	//
	//	//resp, err := e2eSuite.proxyClient.Get(target)
	//	//if err != nil {
	//	//	t.Error(err)
	//	//	return
	//	//}
	//
	//	//body, err := ioutil.ReadAll(resp.Body)
	//	//if err != nil {
	//	//	t.Error(err)
	//	//	return
	//	//}
	//
	//	//fmt.Printf(">%s\n", body)
	//
	//	exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", url)
	//	//exec, err := remotecommand.NewSPDYExecutor(restconfig, "POST", URL)
	//	//exec, err := remotecommand.NewSPDYExecutor(e2eSuite.kubeRestConfig, "POST", req.URL())
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	//wrapper, upgradeRoundTripper, err := spdy.RoundTripperFor(restconfig)
	//	//if err != nil {
	//	//	t.Fatal(err)
	//	//}
	//	//e2eSuite.wrappedRT.transport = wrapper
	//	//exec, err := remotecommand.NewSPDYExecutorForTransports(e2eSuite.wrappedRT, upgradeRoundTripper, "POST", req.URL())
	//	////exec, err := remotecommand.NewSPDYExecutorForTransports(e2eSuite.kubeRestConfig, "POST", req.URL())
	//	//if err != nil {
	//	//	t.Fatal(err)
	//	//}
	//
	//	var (
	//		execOut bytes.Buffer
	//		execErr bytes.Buffer
	//	)
	//
	//	streamOptions := remotecommand.StreamOptions{
	//		Stderr: &execErr,
	//		//Stdin:  os.Stdin,
	//		Stdout: &execOut,
	//		//Tty:    true,
	//	}
	//
	//	err = exec.Stream(streamOptions)
	//	fmt.Printf("%s\n", execOut.String())
	//	fmt.Printf("%s\n", execErr.String())
	//
	//	if err != nil {
	//		fmt.Printf(">%s<\n", err)
	//		t.Fatal(err)
	//	}

	//streamOptions := getStreamOptions(attachOptions, stdin, stdout, stderr)

	//err := attach(cmd.Cfg.Kubeconfig, cmd.pod, attachOptions, cmd.Stdin, cmd.Stdout, cmd.Stderr)
	//if err != nil {
	//	return fmt.Errorf("cannot attach: %v", err)
}
