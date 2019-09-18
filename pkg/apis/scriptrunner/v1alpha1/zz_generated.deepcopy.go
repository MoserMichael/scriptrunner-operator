// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommandRunnerStatus) DeepCopyInto(out *CommandRunnerStatus) {
	*out = *in
	if in.CmdRunHistory != nil {
		in, out := &in.CmdRunHistory, &out.CmdRunHistory
		*out = make([]*CommandStatus, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(CommandStatus)
				**out = **in
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommandRunnerStatus.
func (in *CommandRunnerStatus) DeepCopy() *CommandRunnerStatus {
	if in == nil {
		return nil
	}
	out := new(CommandRunnerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommandStatus) DeepCopyInto(out *CommandStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommandStatus.
func (in *CommandStatus) DeepCopy() *CommandStatus {
	if in == nil {
		return nil
	}
	out := new(CommandStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScriptRunner) DeepCopyInto(out *ScriptRunner) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScriptRunner.
func (in *ScriptRunner) DeepCopy() *ScriptRunner {
	if in == nil {
		return nil
	}
	out := new(ScriptRunner)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ScriptRunner) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScriptRunnerList) DeepCopyInto(out *ScriptRunnerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ScriptRunner, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScriptRunnerList.
func (in *ScriptRunnerList) DeepCopy() *ScriptRunnerList {
	if in == nil {
		return nil
	}
	out := new(ScriptRunnerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ScriptRunnerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScriptRunnerSpec) DeepCopyInto(out *ScriptRunnerSpec) {
	*out = *in
	if in.PythonScript != nil {
		in, out := &in.PythonScript, &out.PythonScript
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScriptRunnerSpec.
func (in *ScriptRunnerSpec) DeepCopy() *ScriptRunnerSpec {
	if in == nil {
		return nil
	}
	out := new(ScriptRunnerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScriptRunnerStatus) DeepCopyInto(out *ScriptRunnerStatus) {
	*out = *in
	if in.CommandStatus != nil {
		in, out := &in.CommandStatus, &out.CommandStatus
		*out = make(map[string]*CommandRunnerStatus, len(*in))
		for key, val := range *in {
			var outVal *CommandRunnerStatus
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(CommandRunnerStatus)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScriptRunnerStatus.
func (in *ScriptRunnerStatus) DeepCopy() *ScriptRunnerStatus {
	if in == nil {
		return nil
	}
	out := new(ScriptRunnerStatus)
	in.DeepCopyInto(out)
	return out
}
