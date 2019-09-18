package scriptrunner

import (
	"context"

	scriptrunnerv1alpha1 "github.com/MoserMichael/scriptrunner/pkg/apis/scriptrunner/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	//labels "k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/source"
        //"github.com/google/uuid"
        "strconv"
        "time"
        "strings"
)

var log = logf.Log.WithName("controller_scriptrunner")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ScriptRunner Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileScriptRunner{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
        log.Info("adding controller")
	// Create a new controller
	c, err := controller.New("scriptrunner-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
                log.Info("adding controller failed", "Error", err.Error())
		return err
	}

	// Watch for changes to primary resource ScriptRunner
	log.Info("Watching resource ScriptRunner")
	err = c.Watch(&source.Kind{Type: &scriptrunnerv1alpha1.ScriptRunner{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
                log.Info("Watching for ScriptRunner resource failed", "Error", err.Error())
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ScriptRunner
        log.Info("Watching for resource Pod")
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &scriptrunnerv1alpha1.ScriptRunner{},
	})
	if err != nil {
                log.Info("Watching for resource Pod failed", "Error", err.Error())
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileScriptRunner implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileScriptRunner{}

// ReconcileScriptRunner reconciles a ScriptRunner object
type ReconcileScriptRunner struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ScriptRunner object and makes changes based on the state read
// and what is in the ScriptRunner.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileScriptRunner) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ScriptRunner")

	// Fetch the ScriptRunner instance
	instance := &scriptrunnerv1alpha1.ScriptRunner{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

        if instance.Status.StatusDescription != "" {
            log.Info("StatusDescription set", "status", instance.Status.StatusDescription)
            return reconcile.Result{}, nil
        }

        /*
        // set property that will identfy this instance,
        if instance.Status.InstanceName == "" {
            log.Info( "instance name not yet assigned;", "iname", instance.Status.InstanceName  )
            nuid, err := uuid.NewUUID()
            if err != nil {
               log.Error(err, "Can't get uuid")
               return reconcile.Result{}, err
            }

            instanceName := nuid.String()
            instance.Status.InstanceName = instanceName
            instance.Status.CommandStatus = make(map[string]*scriptrunnerv1alpha1.CommandRunnerStatus)

            // update the object.
            log.Info("Updating client object", "instanceName", instanceName)
            err = r.client.Status().Update(context.TODO(), instance)
            //err = r.client.Update(context.TODO(), instance)
            if err != nil {
              log.Info("Failed to update instance name ", "error", err.Error())
              return reconcile.Result{}, err
            }

            log.Info("after Updating client object", "instanceName", instanceName)

            return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 2 }, nil
        }

        log.Info( "instance name has been assigned", "instanceName", instance.Status.InstanceName, "NodeLabelSelector", instance.Spec.NodeLabelSelector, "Namespace", request.NamespacedName.Namespace)
        */

        opts := &client.ListOptions{}

        if instance.Spec.NodeLabelSelector != "" {
	    opts.SetLabelSelector( instance.Spec.NodeLabelSelector )
        }

        if request.NamespacedName.Namespace != "default" {
            opts.InNamespace(request.NamespacedName.Namespace)
        }

        nodeList := &corev1.NodeList{}
        err = r.client.List(context.TODO(), opts, nodeList)

        if err != nil{
            log.Error(err, "failed to list nodes")
            // apparently a transient error. (?)
            return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 2 }, nil
        }

        log.Info("after list nodes", "length", len(nodeList.Items), "LabelSelector", instance.Spec.NodeLabelSelector, "Namespace", request.NamespacedName.Namespace)
        for _, node := range nodeList.Items {
            if  !node.Spec.Unschedulable {

                log.Info( "need pod on node", "nodeName", node.ObjectMeta.Name, "instanceName", instance.Status.InstanceName)


                privileged := false
                if instance.Spec.PodType != "normal"  {
                    privileged = true
                }
                mountFS := false
                if instance.Spec.PodType == "elevatedWithFS"  {
                    mountFS = true
                }

                // Define a new Pod object
                pod := newPodForCR(instance, &node, privileged, mountFS)


                // Set ScriptRunner instance as the owner and controller
                if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
                        log.Error( err, "Failed call SetControllerReference" )
                        return reconcile.Result{}, err
                }

                // Check if this Pod already exists
                found := &corev1.Pod{}
                err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
                if err != nil && errors.IsNotFound(err) {
                        reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
                        err = r.client.Create(context.TODO(), pod)
                        if err != nil {
                                log.Info("can't create pod", "podError", err)
                                return reconcile.Result{}, err
                        }

                        reqLogger.Info("Pod successfully created", "Pod.Name", pod.Name, "Pod.Namespace", pod.Namespace, "Node", node.ObjectMeta.Name)

                        // Pod created successfully - don't requeue
                        return reconcile.Result{}, nil
                } else if err != nil {
                        reqLogger.Info("Get pod failed", "error", err)
                        return reconcile.Result{}, err
                }

                log.Info("pod already exists", "Pod.Name", pod.Name, "Pod.Namespace", pod.Namespace, "Node", node.ObjectMeta.Name)

            } else {
                log.Info("Node not schedulable", "Node", node.ObjectMeta.Name)
            }
        }

	// Pod already exists - don't requeue
	//reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

func newPodForCR(cr *scriptrunnerv1alpha1.ScriptRunner, node *corev1.Node, privileged bool, mountVolume bool ) *corev1.Pod {
        instanceName := cr.ObjectMeta.Name
        nodename :=  node.ObjectMeta.Name
        scriptRunnerObjName := cr.Status.InstanceName

	labels := map[string]string{
		"app": "ScriptRunner",
                "ScriptRunnerObjInstance": scriptRunnerObjName,
	}

        privVal := new(bool)
        *privVal =  privileged


        var volumeMounts []corev1.VolumeMount 
        var volumes []corev1.Volume

        if privileged && mountVolume {
            volumeMounts = make([]corev1.VolumeMount,1)
            volumeMounts[0] =  corev1.VolumeMount{
                Name: "rootfs",
                ReadOnly: false,
                MountPath: "/host",
                SubPath:"",
            }
            volumes = make([]corev1.Volume,1)
            volumes[0]  = corev1.Volume{
                Name: "rootfs",
            }

            hostPath := &corev1.HostPathVolumeSource{ 
                Path: "/",
            }

            volumes[0].VolumeSource = corev1.VolumeSource{
                HostPath: hostPath,
            }

        } else {
            volumeMounts = make([]corev1.VolumeMount,0)
            volumes = make([]corev1.Volume,0)
        }

        scriptVal := strings.Join(cr.Spec.PythonScript,"\n")
        scriptVal = strings.ReplaceAll(scriptVal,"\\","\\\\")
        scriptVal = strings.ReplaceAll(scriptVal,"$","\\$")
        scriptVal = strings.ReplaceAll(scriptVal,"\"","\\\"")

        log.Info("script value2","script", scriptVal)

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-" + nodename,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "scriptrunnerpod",
					Image:   "scriptrunnerpod:latest",
                                      //ImagePullPolicy: "IfNotPresent",           // get container image from local registry.
                                        ImagePullPolicy: "Never",           // get container image from local registry.
					Command: []string{"/root/commandrunner"},
                                        Env:    []corev1.EnvVar{
/*                                        
                                            {
                                                Name: "EXTRACT_INFO_OBJECT_NAME",
                                                Value: scriptRunnerObjName,
                                            },
 */
                                            {
                                                Name: "EXTRACT_INFO_NAME",
                                                Value: instanceName,
                                            },
                                            {
                                                Name: "EXTRACT_INFO_NODE_NAME",
                                                Value: nodename,
                                            },
                                            {
                                                Name: "INITIAL_WAIT",
                                                Value:  strconv.FormatInt( int64( cr.Spec.InitialWait ), 10 ),
                                            },
                                            {
                                                Name: "RUN_PERIOD",
                                                Value: strconv.FormatInt( int64( cr.Spec.RunPeriod ), 10 ),
                                            },
                                            {
                                                Name: "NUM_REPETITIONS",
                                                Value: strconv.FormatInt( int64( cr.Spec.NumRepetitions ), 10 ),
                                            },
                                            {
                                                Name: "HISTORY_SIZE",
                                                Value: strconv.FormatInt( int64( cr.Spec.HistorySize ), 10 ),
                                            },
                                            {
                                                Name: "PYTHON_SCRIPT",
                                                Value: scriptVal,
                                            },
                                            {
                                                Name: "PACKAGES_TO_INSTALL",
                                                Value: cr.Spec.PackagesToInstall,
                                            },
                                            {
                                                Name: "PIP_TO_INSTALL",
                                                Value: cr.Spec.PipToInstall,
                                            },
                                           },
                                        SecurityContext:  &corev1.SecurityContext{
                                                Privileged: privVal,
                                            },

                                        VolumeMounts:  volumeMounts,
				},
	                },
                        NodeName: nodename,
                        RestartPolicy: "Never",
                        Volumes: volumes,
		},
            }
}

