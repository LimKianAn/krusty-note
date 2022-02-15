package crds

import (
	"fmt"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func CRDs() ([]apiextensionsv1.CustomResourceDefinition, error) {
	if err := apiextensionsv1.AddToScheme(scheme.Scheme); err != nil {
		return nil, fmt.Errorf("adding to scheme: %v", err)
	}
	rm, err := krusty.MakeKustomizer(krusty.MakeDefaultOptions()).Run(filesys.MakeFsOnDisk(), ".")
	if err != nil {
		return nil, fmt.Errorf("running kustomizer: %v", err)
	}
	rr := rm.Resources()
	oo := make([]runtime.Object, 0, len(rr))
	decoder := scheme.Codecs.UniversalDeserializer()
	for _, r := range rr {
		bb, err := r.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("marshalling json: %v", err)
		}

		o, _, err := decoder.Decode(bb, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("decoding: %v", err)
		}
		oo = append(oo, o)
	}

	crdList := &apiextensionsv1.CustomResourceDefinitionList{}
	if err := meta.SetList(crdList, oo); err != nil {
		return nil, fmt.Errorf("setting crd list: %v", err)
	}
	return crdList.Items, nil
}
