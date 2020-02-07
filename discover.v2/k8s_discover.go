package discover

import (
	"context"
	"github.com/Sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/goodrain/rainbond/cmd/node/option"
	"github.com/goodrain/rainbond/discover/config"
)

type k8sDiscover struct {
	ctx       context.Context
	cancel    context.CancelFunc
	lock      sync.Mutex
	clientset kubernetes.Interface
	cfg       *option.Conf
	projects  map[string]CallbackUpdate
}

func NewK8sDiscover(ctx context.Context, clientset kubernetes.Interface, cfg *option.Conf) Discover {
	ctx, cancel := context.WithCancel(ctx)
	return &k8sDiscover{
		ctx:       ctx,
		cancel:    cancel,
		clientset: clientset,
		cfg:       cfg,
		projects:  make(map[string]CallbackUpdate),
	}
}

func (k *k8sDiscover) Stop() {
	k.cancel()
}

func (k *k8sDiscover) AddProject(name string, callback Callback) {
	k.lock.Lock()
	defer k.lock.Unlock()

	if _, ok := k.projects[name]; !ok {
		cal := &defaultCallBackUpdate{
			callback:  callback,
			endpoints: make(map[string]*config.Endpoint),
		}
		k.projects[name] = cal
		go k.discover(name, cal)
	}
}

func (k *k8sDiscover) AddUpdateProject(name string, callback CallbackUpdate) {
	k.lock.Lock()
	defer k.lock.Unlock()

	if _, ok := k.projects[name]; !ok {
		k.projects[name] = callback
		go k.discover(name, callback)
	}
}

func (k *k8sDiscover) discover(name string, callback CallbackUpdate) {
	ctx, cancel := context.WithCancel(k.ctx)
	defer cancel()

	endpoints := k.list(name)
	if len(endpoints) > 0 {
		callback.UpdateEndpoints(config.SYNC, endpoints...)
	}

	w, err := k.clientset.CoreV1().Pods(k.cfg.RbdNamespace).Watch(metav1.ListOptions{
		LabelSelector: "name=" + name,
	})
	if err != nil {
		k.rewatchWithErr(name, callback, err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-w.ResultChan():
			pod := event.Object.(*corev1.Pod)
			if !isPodReady(pod) {
				continue
			}
			ep := endpointForPod(pod)
			switch event.Type {
			case watch.Deleted:
				callback.UpdateEndpoints(config.DELETE, ep)
			case watch.Added:
				callback.UpdateEndpoints(config.ADD, ep)
			case watch.Modified:
				callback.UpdateEndpoints(config.UPDATE, ep)
			case watch.Error:
				k.rewatchWithErr(name, callback, err)
			}
		}
	}
}

func (k *k8sDiscover) removeProject(name string) {
	k.lock.Lock()
	defer k.lock.Unlock()
	if _, ok := k.projects[name]; ok {
		delete(k.projects, name)
	}
}

func (k *k8sDiscover) rewatchWithErr(name string, callback CallbackUpdate, err error) {
	logrus.Debugf("name: %s; monitor discover get watch error: %s, remove this watch target first, and then sleep 10 sec, we will re-watch it", name, err.Error())
	callback.Error(err)
	k.removeProject(name)
	time.Sleep(10 * time.Second)
	k.AddUpdateProject(name, callback)
}

func (k *k8sDiscover) list(name string) []*config.Endpoint {
	podList, err := k.clientset.CoreV1().Pods(k.cfg.RbdNamespace).List(metav1.ListOptions{
		LabelSelector: "name=" + name,
	})
	if err != nil {
		logrus.Warningf("list pods for %s: %v", name, err)
		return nil
	}

	var endpoints []*config.Endpoint
	var notReadyEp *config.Endpoint
	for _, pod := range podList.Items {
		ep := endpointForPod(&pod)
		if isPodReady(&pod) {
			endpoints = append(endpoints, ep)
			continue
		}
		if notReadyEp == nil {
			notReadyEp = endpointForPod(&pod)
		}
	}

	// If there are no ready endpoints, a not ready endpoint is used
	if len(endpoints) == 0 && notReadyEp != nil {
		endpoints = append(endpoints, notReadyEp)
	}

	return endpoints
}

func endpointForPod(pod *corev1.Pod) *config.Endpoint {
	return &config.Endpoint{
		Name: pod.Name,
		URL:  pod.Status.PodIP,
	}
}

func isPodReady(pod *corev1.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
