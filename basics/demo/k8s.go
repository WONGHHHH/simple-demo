package demo

import (
	"encoding/json"
    "flag"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

type DeploymentInfo struct {
    Namespace string `json:"namespace"`
    Name      string `json:"name"`
    Replicas  int32  `json:"replicas"`
}

type PodCount struct {
    Count int `json:"count"`
}

func main() {
    // 获取 Kubernetes 配置
    kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    flag.Parse()

    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        log.Fatal(err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatal(err)
    }

    router := mux.NewRouter()

    // 处理 GET /deployment/pod-count 请求
    router.HandleFunc("/deployment/pod-count", func(w http.ResponseWriter, r *http.Request) {
        deploymentInfo := DeploymentInfo{}
        err := json.NewDecoder(r.Body).Decode(&deploymentInfo)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        pods, err := clientset.CoreV1().Pods(deploymentInfo.Namespace).List(r.Context(), v1.ListOptions{
            LabelSelector: fmt.Sprintf("app=%s", deploymentInfo.Name),
        })
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        podCount := PodCount{Count: len(pods.Items)}
        json.NewEncoder(w).Encode(podCount)
    }).Methods("GET")

    // 处理 POST /deployment/scale 请求
    router.HandleFunc("/deployment/scale", func(w http.ResponseWriter, r *http.Request) {
        deploymentInfo := DeploymentInfo{}
        err := json.NewDecoder(r.Body).Decode(&deploymentInfo)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        deployment, err := clientset.AppsV1().Deployments(deploymentInfo.Namespace).Get(r.Context(), deploymentInfo.Name, v1.GetOptions{})
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        deployment.Spec.Replicas = &deploymentInfo.Replicas
        _, err = clientset.AppsV1().Deployments(deploymentInfo.Namespace).Update(r.Context(), deployment, v1.UpdateOptions{})
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // 等待 Deployment 的 Pod 数量达到预期
        for {
            pods, err := clientset.CoreV1().Pods(deploymentInfo.Namespace).List(r.Context(), v1.ListOptions{
                LabelSelector: fmt.Sprintf("app=%s", deploymentInfo.Name),
            })
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            if int32(len(pods.Items)) == deploymentInfo.Replicas {
                break
            }
        }

        w.WriteHeader(http.StatusOK)
    }).Methods("POST")

    log.Fatal(http.ListenAndServe(":8080", router))
}