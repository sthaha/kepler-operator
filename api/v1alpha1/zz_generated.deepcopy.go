//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectorSpec) DeepCopyInto(out *CollectorSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectorSpec.
func (in *CollectorSpec) DeepCopy() *CollectorSpec {
	if in == nil {
		return nil
	}
	out := new(CollectorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomHeader) DeepCopyInto(out *CustomHeader) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomHeader.
func (in *CustomHeader) DeepCopy() *CustomHeader {
	if in == nil {
		return nil
	}
	out := new(CustomHeader)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EstimatorSpec) DeepCopyInto(out *EstimatorSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EstimatorSpec.
func (in *EstimatorSpec) DeepCopy() *EstimatorSpec {
	if in == nil {
		return nil
	}
	out := new(EstimatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Kepler) DeepCopyInto(out *Kepler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Kepler.
func (in *Kepler) DeepCopy() *Kepler {
	if in == nil {
		return nil
	}
	out := new(Kepler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Kepler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerList) DeepCopyInto(out *KeplerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Kepler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerList.
func (in *KeplerList) DeepCopy() *KeplerList {
	if in == nil {
		return nil
	}
	out := new(KeplerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KeplerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerSpec) DeepCopyInto(out *KeplerSpec) {
	*out = *in
	if in.ModelServerExporter != nil {
		in, out := &in.ModelServerExporter, &out.ModelServerExporter
		*out = new(ModelServerExporterSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ModelServerFeatures != nil {
		in, out := &in.ModelServerFeatures, &out.ModelServerFeatures
		*out = new(ModelServerFeaturesSpec)
		**out = **in
	}
	if in.Estimator != nil {
		in, out := &in.Estimator, &out.Estimator
		*out = new(EstimatorSpec)
		**out = **in
	}
	if in.Collector != nil {
		in, out := &in.Collector, &out.Collector
		*out = new(CollectorSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerSpec.
func (in *KeplerSpec) DeepCopy() *KeplerSpec {
	if in == nil {
		return nil
	}
	out := new(KeplerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerStatus) DeepCopyInto(out *KeplerStatus) {
	*out = *in
	in.Conditions.DeepCopyInto(&out.Conditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerStatus.
func (in *KeplerStatus) DeepCopy() *KeplerStatus {
	if in == nil {
		return nil
	}
	out := new(KeplerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelServerExporterSpec) DeepCopyInto(out *ModelServerExporterSpec) {
	*out = *in
	if in.ModelServerTrainer != nil {
		in, out := &in.ModelServerTrainer, &out.ModelServerTrainer
		*out = new(ModelServerTrainerSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelServerExporterSpec.
func (in *ModelServerExporterSpec) DeepCopy() *ModelServerExporterSpec {
	if in == nil {
		return nil
	}
	out := new(ModelServerExporterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelServerFeaturesSpec) DeepCopyInto(out *ModelServerFeaturesSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelServerFeaturesSpec.
func (in *ModelServerFeaturesSpec) DeepCopy() *ModelServerFeaturesSpec {
	if in == nil {
		return nil
	}
	out := new(ModelServerFeaturesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelServerTrainerSpec) DeepCopyInto(out *ModelServerTrainerSpec) {
	*out = *in
	if in.PromHeaders != nil {
		in, out := &in.PromHeaders, &out.PromHeaders
		*out = make([]CustomHeader, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelServerTrainerSpec.
func (in *ModelServerTrainerSpec) DeepCopy() *ModelServerTrainerSpec {
	if in == nil {
		return nil
	}
	out := new(ModelServerTrainerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RatioMetrics) DeepCopyInto(out *RatioMetrics) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RatioMetrics.
func (in *RatioMetrics) DeepCopy() *RatioMetrics {
	if in == nil {
		return nil
	}
	out := new(RatioMetrics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Sources) DeepCopyInto(out *Sources) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Sources.
func (in *Sources) DeepCopy() *Sources {
	if in == nil {
		return nil
	}
	out := new(Sources)
	in.DeepCopyInto(out)
	return out
}
