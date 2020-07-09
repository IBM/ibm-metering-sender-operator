// +build !ignore_autogenerated

//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeteringSender) DeepCopyInto(out *MeteringSender) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeteringSender.
func (in *MeteringSender) DeepCopy() *MeteringSender {
	if in == nil {
		return nil
	}
	out := new(MeteringSender)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MeteringSender) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeteringSenderList) DeepCopyInto(out *MeteringSenderList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MeteringSender, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeteringSenderList.
func (in *MeteringSenderList) DeepCopy() *MeteringSenderList {
	if in == nil {
		return nil
	}
	out := new(MeteringSenderList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MeteringSenderList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeteringSenderSpec) DeepCopyInto(out *MeteringSenderSpec) {
	*out = *in
	out.Sender = in.Sender
	out.MongoDB = in.MongoDB
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeteringSenderSpec.
func (in *MeteringSenderSpec) DeepCopy() *MeteringSenderSpec {
	if in == nil {
		return nil
	}
	out := new(MeteringSenderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeteringSenderSpecSender) DeepCopyInto(out *MeteringSenderSpecSender) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeteringSenderSpecSender.
func (in *MeteringSenderSpecSender) DeepCopy() *MeteringSenderSpecSender {
	if in == nil {
		return nil
	}
	out := new(MeteringSenderSpecSender)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeteringSenderStatus) DeepCopyInto(out *MeteringSenderStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeteringSenderStatus.
func (in *MeteringSenderStatus) DeepCopy() *MeteringSenderStatus {
	if in == nil {
		return nil
	}
	out := new(MeteringSenderStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeteringSpecMongoDB) DeepCopyInto(out *MeteringSpecMongoDB) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeteringSpecMongoDB.
func (in *MeteringSpecMongoDB) DeepCopy() *MeteringSpecMongoDB {
	if in == nil {
		return nil
	}
	out := new(MeteringSpecMongoDB)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MeteringStatus) DeepCopyInto(out *MeteringStatus) {
	*out = *in
	if in.PodNames != nil {
		in, out := &in.PodNames, &out.PodNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MeteringStatus.
func (in *MeteringStatus) DeepCopy() *MeteringStatus {
	if in == nil {
		return nil
	}
	out := new(MeteringStatus)
	in.DeepCopyInto(out)
	return out
}