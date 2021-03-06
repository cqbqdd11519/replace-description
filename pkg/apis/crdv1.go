package apis

import (
	"encoding/json"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CustomResourceDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// spec describes how the user wants the resources to appear
	Spec CustomResourceDefinitionSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// status indicates the actual state of the CustomResourceDefinition
	// +optional
	Status crdv1.CustomResourceDefinitionStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// CustomResourceDefinitionSpec describes how a user wants their resource to appear
type CustomResourceDefinitionSpec struct {
	// group is the API group of the defined custom resource.
	// The custom resources are served under `/apis/<group>/...`.
	// Must match the name of the CustomResourceDefinition (in the form `<names.plural>.<group>`).
	Group string `json:"group" protobuf:"bytes,1,opt,name=group"`
	// names specify the resource and kind names for the custom resource.
	Names crdv1.CustomResourceDefinitionNames `json:"names" protobuf:"bytes,3,opt,name=names"`
	// scope indicates whether the defined custom resource is cluster- or namespace-scoped.
	// Allowed values are `Cluster` and `Namespaced`.
	Scope crdv1.ResourceScope `json:"scope" protobuf:"bytes,4,opt,name=scope,casttype=ResourceScope"`
	// versions is the list of all API versions of the defined custom resource.
	// Version names are used to compute the order in which served versions are listed in API discovery.
	// If the version string is "kube-like", it will sort above non "kube-like" version strings, which are ordered
	// lexicographically. "Kube-like" versions start with a "v", then are followed by a number (the major version),
	// then optionally the string "alpha" or "beta" and another number (the minor version). These are sorted first
	// by GA > beta > alpha (where GA is a version with no suffix such as beta or alpha), and then by comparing
	// major version, then minor version. An example sorted list of versions:
	// v10, v2, v1, v11beta2, v10beta3, v3beta1, v12alpha1, v11alpha2, foo1, foo10.
	Versions []CustomResourceDefinitionVersion `json:"versions" protobuf:"bytes,7,rep,name=versions"`

	// conversion defines conversion settings for the CRD.
	// +optional
	Conversion *crdv1.CustomResourceConversion `json:"conversion,omitempty" protobuf:"bytes,9,opt,name=conversion"`

	// preserveUnknownFields indicates that object fields which are not specified
	// in the OpenAPI schema should be preserved when persisting to storage.
	// apiVersion, kind, metadata and known fields inside metadata are always preserved.
	// This field is deprecated in favor of setting `x-preserve-unknown-fields` to true in `spec.versions[*].schema.openAPIV3Schema`.
	// See https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#pruning-versus-preserving-unknown-fields for details.
	// +optional
	PreserveUnknownFields bool `json:"preserveUnknownFields,omitempty" protobuf:"varint,10,opt,name=preserveUnknownFields"`
}

// CustomResourceDefinitionVersion describes a version for CRD.
type CustomResourceDefinitionVersion struct {
	// name is the version name, e.g. ???v1???, ???v2beta1???, etc.
	// The custom resources are served under this version at `/apis/<group>/<version>/...` if `served` is true.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// served is a flag enabling/disabling this version from being served via REST APIs
	Served bool `json:"served" protobuf:"varint,2,opt,name=served"`
	// storage indicates this version should be used when persisting custom resources to storage.
	// There must be exactly one version with storage=true.
	Storage bool `json:"storage" protobuf:"varint,3,opt,name=storage"`
	// deprecated indicates this version of the custom resource API is deprecated.
	// When set to true, API requests to this version receive a warning header in the server response.
	// Defaults to false.
	// +optional
	Deprecated bool `json:"deprecated,omitempty" protobuf:"varint,7,opt,name=deprecated"`
	// deprecationWarning overrides the default warning returned to API clients.
	// May only be set when `deprecated` is true.
	// The default warning indicates this version is deprecated and recommends use
	// of the newest served version of equal or greater stability, if one exists.
	// +optional
	DeprecationWarning *string `json:"deprecationWarning,omitempty" protobuf:"bytes,8,opt,name=deprecationWarning"`
	// schema describes the schema used for validation, pruning, and defaulting of this version of the custom resource.
	// +optional
	Schema *CustomResourceValidation `json:"schema,omitempty" protobuf:"bytes,4,opt,name=schema"`
	// subresources specify what subresources this version of the defined custom resource have.
	// +optional
	Subresources *crdv1.CustomResourceSubresources `json:"subresources,omitempty" protobuf:"bytes,5,opt,name=subresources"`
	// additionalPrinterColumns specifies additional columns returned in Table output.
	// See https://kubernetes.io/docs/reference/using-api/api-concepts/#receiving-resources-as-tables for details.
	// If no columns are specified, a single column displaying the age of the custom resource is used.
	// +optional
	AdditionalPrinterColumns []crdv1.CustomResourceColumnDefinition `json:"additionalPrinterColumns,omitempty" protobuf:"bytes,6,rep,name=additionalPrinterColumns"`
}

type CustomResourceValidation struct {
	// openAPIV3Schema is the OpenAPI v3 schema to use for validation and pruning.
	// +optional
	OpenAPIV3Schema *JSONSchemaProps `json:"openAPIV3Schema,omitempty" protobuf:"bytes,1,opt,name=openAPIV3Schema"`
}

type JSONSchemaProps struct {
	ID          string              `json:"id,omitempty" protobuf:"bytes,1,opt,name=id"`
	Schema      crdv1.JSONSchemaURL `json:"$schema,omitempty" protobuf:"bytes,2,opt,name=schema"`
	Ref         *string             `json:"$ref,omitempty" protobuf:"bytes,3,opt,name=ref"`
	Description string              `json:"description,omitempty" protobuf:"bytes,4,opt,name=description"`
	Type        string              `json:"type,omitempty" protobuf:"bytes,5,opt,name=type"`

	// format is an OpenAPI v3 format string. Unknown formats are ignored. The following formats are validated:
	//
	// - bsonobjectid: a bson object ID, i.e. a 24 characters hex string
	// - uri: an URI as parsed by Golang net/url.ParseRequestURI
	// - email: an email address as parsed by Golang net/mail.ParseAddress
	// - hostname: a valid representation for an Internet host name, as defined by RFC 1034, section 3.1 [RFC1034].
	// - ipv4: an IPv4 IP as parsed by Golang net.ParseIP
	// - ipv6: an IPv6 IP as parsed by Golang net.ParseIP
	// - cidr: a CIDR as parsed by Golang net.ParseCIDR
	// - mac: a MAC address as parsed by Golang net.ParseMAC
	// - uuid: an UUID that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{4}-?[0-9a-f]{12}$
	// - uuid3: an UUID3 that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?3[0-9a-f]{3}-?[0-9a-f]{4}-?[0-9a-f]{12}$
	// - uuid4: an UUID4 that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?4[0-9a-f]{3}-?[89ab][0-9a-f]{3}-?[0-9a-f]{12}$
	// - uuid5: an UUID5 that allows uppercase defined by the regex (?i)^[0-9a-f]{8}-?[0-9a-f]{4}-?5[0-9a-f]{3}-?[89ab][0-9a-f]{3}-?[0-9a-f]{12}$
	// - isbn: an ISBN10 or ISBN13 number string like "0321751043" or "978-0321751041"
	// - isbn10: an ISBN10 number string like "0321751043"
	// - isbn13: an ISBN13 number string like "978-0321751041"
	// - creditcard: a credit card number defined by the regex ^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11})$ with any non digit characters mixed in
	// - ssn: a U.S. social security number following the regex ^\\d{3}[- ]?\\d{2}[- ]?\\d{4}$
	// - hexcolor: an hexadecimal color code like "#FFFFFF: following the regex ^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$
	// - rgbcolor: an RGB color code like rgb like "rgb(255,255,2559"
	// - byte: base64 encoded binary data
	// - password: any kind of string
	// - date: a date string like "2006-01-02" as defined by full-date in RFC3339
	// - duration: a duration string like "22 ns" as parsed by Golang time.ParseDuration or compatible with Scala duration format
	// - datetime: a date time string like "2014-12-15T19:30:20.000Z" as defined by date-time in RFC3339.
	Format string `json:"format,omitempty" protobuf:"bytes,6,opt,name=format"`

	Title string `json:"title,omitempty" protobuf:"bytes,7,opt,name=title"`
	// default is a default value for undefined object fields.
	// Defaulting is a beta feature under the CustomResourceDefaulting feature gate.
	// Defaulting requires spec.preserveUnknownFields to be false.
	Default              *crdv1.JSON                  `json:"default,omitempty" protobuf:"bytes,8,opt,name=default"`
	Maximum              *float64                     `json:"maximum,omitempty" protobuf:"bytes,9,opt,name=maximum"`
	ExclusiveMaximum     bool                         `json:"exclusiveMaximum,omitempty" protobuf:"bytes,10,opt,name=exclusiveMaximum"`
	Minimum              *float64                     `json:"minimum,omitempty" protobuf:"bytes,11,opt,name=minimum"`
	ExclusiveMinimum     bool                         `json:"exclusiveMinimum,omitempty" protobuf:"bytes,12,opt,name=exclusiveMinimum"`
	MaxLength            *int64                       `json:"maxLength,omitempty" protobuf:"bytes,13,opt,name=maxLength"`
	MinLength            *int64                       `json:"minLength,omitempty" protobuf:"bytes,14,opt,name=minLength"`
	Pattern              string                       `json:"pattern,omitempty" protobuf:"bytes,15,opt,name=pattern"`
	MaxItems             *int64                       `json:"maxItems,omitempty" protobuf:"bytes,16,opt,name=maxItems"`
	MinItems             *int64                       `json:"minItems,omitempty" protobuf:"bytes,17,opt,name=minItems"`
	UniqueItems          bool                         `json:"uniqueItems,omitempty" protobuf:"bytes,18,opt,name=uniqueItems"`
	MultipleOf           *float64                     `json:"multipleOf,omitempty" protobuf:"bytes,19,opt,name=multipleOf"`
	Enum                 []crdv1.JSON                 `json:"enum,omitempty" protobuf:"bytes,20,rep,name=enum"`
	MaxProperties        *int64                       `json:"maxProperties,omitempty" protobuf:"bytes,21,opt,name=maxProperties"`
	MinProperties        *int64                       `json:"minProperties,omitempty" protobuf:"bytes,22,opt,name=minProperties"`
	Required             []string                     `json:"required,omitempty" protobuf:"bytes,23,rep,name=required"`
	Items                *JSONSchemaPropsOrArray      `json:"items,omitempty" protobuf:"bytes,24,opt,name=items"`
	AllOf                []JSONSchemaProps            `json:"allOf,omitempty" protobuf:"bytes,25,rep,name=allOf"`
	OneOf                []JSONSchemaProps            `json:"oneOf,omitempty" protobuf:"bytes,26,rep,name=oneOf"`
	AnyOf                []JSONSchemaProps            `json:"anyOf,omitempty" protobuf:"bytes,27,rep,name=anyOf"`
	Not                  *JSONSchemaProps             `json:"not,omitempty" protobuf:"bytes,28,opt,name=not"`
	Properties           map[string]*JSONSchemaProps  `json:"properties,omitempty" protobuf:"bytes,29,rep,name=properties"`
	AdditionalProperties *crdv1.JSONSchemaPropsOrBool `json:"additionalProperties,omitempty" protobuf:"bytes,30,opt,name=additionalProperties"`
	PatternProperties    map[string]JSONSchemaProps   `json:"patternProperties,omitempty" protobuf:"bytes,31,rep,name=patternProperties"`
	Dependencies         crdv1.JSONSchemaDependencies `json:"dependencies,omitempty" protobuf:"bytes,32,opt,name=dependencies"`
	AdditionalItems      *crdv1.JSONSchemaPropsOrBool `json:"additionalItems,omitempty" protobuf:"bytes,33,opt,name=additionalItems"`
	Definitions          crdv1.JSONSchemaDefinitions  `json:"definitions,omitempty" protobuf:"bytes,34,opt,name=definitions"`
	ExternalDocs         *crdv1.ExternalDocumentation `json:"externalDocs,omitempty" protobuf:"bytes,35,opt,name=externalDocs"`
	Example              *crdv1.JSON                  `json:"example,omitempty" protobuf:"bytes,36,opt,name=example"`
	Nullable             bool                         `json:"nullable,omitempty" protobuf:"bytes,37,opt,name=nullable"`

	// x-kubernetes-preserve-unknown-fields stops the API server
	// decoding step from pruning fields which are not specified
	// in the validation schema. This affects fields recursively,
	// but switches back to normal pruning behaviour if nested
	// properties or additionalProperties are specified in the schema.
	// This can either be true or undefined. False is forbidden.
	XPreserveUnknownFields *bool `json:"x-kubernetes-preserve-unknown-fields,omitempty" protobuf:"bytes,38,opt,name=xKubernetesPreserveUnknownFields"`

	// x-kubernetes-embedded-resource defines that the value is an
	// embedded Kubernetes runtime.Object, with TypeMeta and
	// ObjectMeta. The type must be object. It is allowed to further
	// restrict the embedded object. kind, apiVersion and metadata
	// are validated automatically. x-kubernetes-preserve-unknown-fields
	// is allowed to be true, but does not have to be if the object
	// is fully specified (up to kind, apiVersion, metadata).
	XEmbeddedResource bool `json:"x-kubernetes-embedded-resource,omitempty" protobuf:"bytes,39,opt,name=xKubernetesEmbeddedResource"`

	// x-kubernetes-int-or-string specifies that this value is
	// either an integer or a string. If this is true, an empty
	// type is allowed and type as child of anyOf is permitted
	// if following one of the following patterns:
	//
	// 1) anyOf:
	//    - type: integer
	//    - type: string
	// 2) allOf:
	//    - anyOf:
	//      - type: integer
	//      - type: string
	//    - ... zero or more
	XIntOrString bool `json:"x-kubernetes-int-or-string,omitempty" protobuf:"bytes,40,opt,name=xKubernetesIntOrString"`

	// x-kubernetes-list-map-keys annotates an array with the x-kubernetes-list-type `map` by specifying the keys used
	// as the index of the map.
	//
	// This tag MUST only be used on lists that have the "x-kubernetes-list-type"
	// extension set to "map". Also, the values specified for this attribute must
	// be a scalar typed field of the child structure (no nesting is supported).
	//
	// The properties specified must either be required or have a default value,
	// to ensure those properties are present for all list items.
	//
	// +optional
	XListMapKeys []string `json:"x-kubernetes-list-map-keys,omitempty" protobuf:"bytes,41,rep,name=xKubernetesListMapKeys"`

	// x-kubernetes-list-type annotates an array to further describe its topology.
	// This extension must only be used on lists and may have 3 possible values:
	//
	// 1) `atomic`: the list is treated as a single entity, like a scalar.
	//      Atomic lists will be entirely replaced when updated. This extension
	//      may be used on any type of list (struct, scalar, ...).
	// 2) `set`:
	//      Sets are lists that must not have multiple items with the same value. Each
	//      value must be a scalar, an object with x-kubernetes-map-type `atomic` or an
	//      array with x-kubernetes-list-type `atomic`.
	// 3) `map`:
	//      These lists are like maps in that their elements have a non-index key
	//      used to identify them. Order is preserved upon merge. The map tag
	//      must only be used on a list with elements of type object.
	// Defaults to atomic for arrays.
	// +optional
	XListType *string `json:"x-kubernetes-list-type,omitempty" protobuf:"bytes,42,opt,name=xKubernetesListType"`

	// x-kubernetes-map-type annotates an object to further describe its topology.
	// This extension must only be used when type is object and may have 2 possible values:
	//
	// 1) `granular`:
	//      These maps are actual maps (key-value pairs) and each fields are independent
	//      from each other (they can each be manipulated by separate actors). This is
	//      the default behaviour for all maps.
	// 2) `atomic`: the list is treated as a single entity, like a scalar.
	//      Atomic maps will be entirely replaced when updated.
	// +optional
	XMapType *string `json:"x-kubernetes-map-type,omitempty" protobuf:"bytes,43,opt,name=xKubernetesMapType"`
}

type JSONSchemaPropsOrArray struct {
	Schema      *JSONSchemaProps  `protobuf:"bytes,1,opt,name=schema"`
	JSONSchemas []JSONSchemaProps `protobuf:"bytes,2,rep,name=jSONSchemas"`
}

func (s JSONSchemaPropsOrArray) MarshalJSON() ([]byte, error) {
	if len(s.JSONSchemas) > 0 {
		return json.Marshal(s.JSONSchemas)
	}
	return json.Marshal(s.Schema)
}

func (s *JSONSchemaPropsOrArray) UnmarshalJSON(data []byte) error {
	var nw JSONSchemaPropsOrArray
	var first byte
	if len(data) > 1 {
		first = data[0]
	}
	if first == '{' {
		var sch JSONSchemaProps
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.JSONSchemas); err != nil {
			return err
		}
	}
	*s = nw
	return nil
}
