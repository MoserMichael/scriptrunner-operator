package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ScriptRunnerSpec defines the desired state of ScriptRunner
// +k8s:openapi-gen=true
type ScriptRunnerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

        // PythonScript is the script to run on the node
        PythonScript []string `json:"pythonScript"`

        // PackagesToInstall if string not empty: prior to running the script it will install the following distro packages on the docker image 
        PackagesToInstall string `json:"packagesToInstall"`

        // PipToInstall if string not empty: prior to running the script it will install the following pip packages on the docker image
        PipToInstall string `json:"pipToInstall"`

        // InitialWait if not zero then the number of milliseconds to wait before first run
        InitialWait int32 `json:"initialWait"`

        // RunPeriod specifies the interval between invocations of the script (in milliseconds)
        RunPeriod int32 `json:"runPeriod"`

        // NumRepetitions how many times the script is run (-1 means infinite loop) The pod is stopped on last invocation (if not -1)
        NumRepetitions int32 `json:"numRepetitions"`

        // HistorySize number of entries in command history
        HistorySize  int32 `json:"historySize"`

        // Selector selects the set of nodes that this pod is run on (empty string - all nodes)
        NodeLabelSelector string `json:"nodeLabelSelector"`

        // PodType selects the type of pod to run ("normal" "elevated" "elevatedWitFS")
        PodType  string `json:"podType"`
}

// CommandStatus defines the observed state of CommandStatus
// +k8s:openapi-gen=true
type CommandStatus struct {

    // SerialNo serial  number of command invocation; incremented by one with each invocation
    SerialNo    int64   `json:"serialNo"`

    // CommandStatus the exit status of running the script command
    CommandStatus int  `json:"commandStatus"`

    // CommandStdOut the standard output of the command
    CommandStdOut string `json:"commandStdOut"`

    // CommandStdErr the standard error of the command
    CommandStdErr string `json:"commandStdErr"`

    // TimeStart time when command was started
    TimeStart string `json:"timeStart"`

    // TimeStart time when command was started
    TimeEnd string `json:"timeEnd"`
}

// CommandStatus defines the observed state of CommandStatus
// +k8s:openapi-gen=true
// +kubebuilder:subresource:cmdRunHistory
type CommandRunnerStatus struct {

    // PodInstanceName
    PodInstanceName string           `json:"podInstanceName"`

    // CmdRunHistory a history of command invocations 
    // +kubebuilder:validation:MinItems=0
    CmdRunHistory []*CommandStatus  `json:"cmdRunHistory,omitempty"`
}

// ScriptRunnerStatus defines the observed state of ScriptRunner
// +k8s:openapi-gen=true
// +kubebuilder:subresource:commandStatus
type ScriptRunnerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.htmla

        // StatusDescription a description of the status
        StatusDescription string `json:"statusDescription"`

        // InstanceName the internal name of the node (assigned by operator)
        InstanceName string `json:"instanceName"`

        // CommandNodes the result of the script infocations that ran on a particular node
        // +kubebuilder:validation:MinItems=0
        CommandStatus map[string]*CommandRunnerStatus  `json:"commandStatus,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ScriptRunner is the Schema for the scriptrunners API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ScriptRunner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScriptRunnerSpec   `json:"spec,omitempty"`
	Status ScriptRunnerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ScriptRunnerList contains a list of ScriptRunner
type ScriptRunnerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScriptRunner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScriptRunner{}, &ScriptRunnerList{})
}
