Warning: truncated output (original token count: 35325)
Total output lines: 3345

package sarif

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

// A physical or virtual address, or a range of addresses, in an 'addressable
// region' (memory or a binary file).
type Address struct {
	// The address expressed as a byte offset from the start of the addressable
	// region.
	AbsoluteAddress int `json:"absoluteAddress,omitempty" yaml:"absoluteAddress,omitempty" mapstructure:"absoluteAddress,omitempty"`

	// A human-readable fully qualified name that is associated with the address.
	FullyQualifiedName *string `json:"fullyQualifiedName,omitempty" yaml:"fullyQualifiedName,omitempty" mapstructure:"fullyQualifiedName,omitempty"`

	// The index within run.addresses of the cached object for this address.
	Index int `json:"index,omitempty" yaml:"index,omitempty" mapstructure:"index,omitempty"`

	// An open-ended string that identifies the address kind. 'data', 'function',
	// 'header','instruction', 'module', 'page', 'section', 'segment', 'stack',
	// 'stackFrame', 'table' are well-known values.
	Kind *string `json:"kind,omitempty" yaml:"kind,omitempty" mapstructure:"kind,omitempty"`

	// The number of bytes in this range of addresses.
	Length *int `json:"length,omitempty" yaml:"length,omitempty" mapstructure:"length,omitempty"`

	// A name that is associated with the address, e.g., '.text'.
	Name *string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`

	// The byte offset of this address from the absolute or relative address of the
	// parent object.
	OffsetFromParent *int `json:"offsetFromParent,omitempty" yaml:"offsetFromParent,omitempty" mapstructure:"offsetFromParent,omitempty"`

	// The index within run.addresses of the parent object.
	ParentIndex int `json:"parentIndex,omitempty" yaml:"parentIndex,omitempty" mapstructure:"parentIndex,omitempty"`

	// Key/value pairs that provide additional information about the address.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The address expressed as a byte offset from the absolute address of the
	// top-most parent object.
	RelativeAddress *int `json:"relativeAddress,omitempty" yaml:"relativeAddress,omitempty" mapstructure:"relativeAddress,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Address) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain Address
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["absoluteAddress"]; !ok || v == nil {
		plain.AbsoluteAddress = -1.0
	}
	if -1 > plain.AbsoluteAddress {
		return fmt.Errorf("field %s: must be >= %v", "absoluteAddress", -1)
	}
	if v, ok := raw["index"]; !ok || v == nil {
		plain.Index = -1.0
	}
	if -1 > plain.Index {
		return fmt.Errorf("field %s: must be >= %v", "index", -1)
	}
	if v, ok := raw["parentIndex"]; !ok || v == nil {
		plain.ParentIndex = -1.0
	}
	if -1 > plain.ParentIndex {
		return fmt.Errorf("field %s: must be >= %v", "parentIndex", -1)
	}
	*j = Address(plain)
	return nil
}

// A single artifact. In some cases, this artifact might be nested within another
// artifact.
type Artifact struct {
	// The contents of the artifact.
	Contents *ArtifactContent `json:"contents,omitempty" yaml:"contents,omitempty" mapstructure:"contents,omitempty"`

	// A short description of the artifact.
	Description Message `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// Specifies the encoding for an artifact object that refers to a text file.
	Encoding *string `json:"encoding,omitempty" yaml:"encoding,omitempty" mapstructure:"encoding,omitempty"`

	// A dictionary, each of whose keys is the name of a hash function and each of
	// whose values is the hashed value of the artifact produced by the specified hash
	// function.
	Hashes map[string]string `json:"hashes,omitempty" yaml:"hashes,omitempty" mapstructure:"hashes,omitempty"`

	// The Coordinated Universal Time (UTC) date and time at which the artifact was
	// most recently modified. See "Date/time properties" in the SARIF spec for the
	// required format.
	LastModifiedTimeUtc *time.Time `json:"lastModifiedTimeUtc,omitempty" yaml:"lastModifiedTimeUtc,omitempty" mapstructure:"lastModifiedTimeUtc,omitempty"`

	// The length of the artifact in bytes.
	Length int `json:"length,omitempty" yaml:"length,omitempty" mapstructure:"length,omitempty"`

	// The location of the artifact.
	Location *ArtifactLocation `json:"location,omitempty" yaml:"location,omitempty" mapstructure:"location,omitempty"`

	// The MIME type (RFC 2045) of the artifact.
	MimeType *string `json:"mimeType,omitempty" yaml:"mimeType,omitempty" mapstructure:"mimeType,omitempty"`

	// The offset in bytes of the artifact within its containing artifact.
	Offset *int `json:"offset,omitempty" yaml:"offset,omitempty" mapstructure:"offset,omitempty"`

	// IDentifies the index of the immediate parent of the artifact, if this artifact
	// is nested.
	ParentIndex int `json:"parentIndex,omitempty" yaml:"parentIndex,omitempty" mapstructure:"parentIndex,omitempty"`

	// Key/value pairs that provide additional information about the artifact.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The role or roles played by the artifact in the analysis.
	Roles []ArtifactRolesElem `json:"roles,omitempty" yaml:"roles,omitempty" mapstructure:"roles,omitempty"`

	// Specifies the source language for any artifact object that refers to a text
	// file that contains source code.
	SourceLanguage *string `json:"sourceLanguage,omitempty" yaml:"sourceLanguage,omitempty" mapstructure:"sourceLanguage,omitempty"`
}

// A change to a single artifact.
type ArtifactChange struct {
	// The location of the artifact to change.
	ArtifactLocation ArtifactLocation `json:"artifactLocation" yaml:"artifactLocation" mapstructure:"artifactLocation"`

	// Key/value pairs that provide additional information about the change.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// An array of replacement objects, each of which represents the replacement of a
	// single region in a single artifact specified by 'artifactLocation'.
	Replacements []Replacement `json:"replacements" yaml:"replacements" mapstructure:"replacements"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ArtifactChange) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["artifactLocation"]; raw != nil && !ok {
		return fmt.Errorf("field artifactLocation in ArtifactChange: required")
	}
	if _, ok := raw["replacements"]; raw != nil && !ok {
		return fmt.Errorf("field replacements in ArtifactChange: required")
	}
	type Plain ArtifactChange
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if plain.Replacements != nil && len(plain.Replacements) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "replacements", 1)
	}
	*j = ArtifactChange(plain)
	return nil
}

// Represents the contents of an artifact.
type ArtifactContent struct {
	// MIME Base64-encoded content from a binary artifact, or from a text artifact in
	// its original encoding.
	Binary *string `json:"binary,omitempty" yaml:"binary,omitempty" mapstructure:"binary,omitempty"`

	// Key/value pairs that provide additional information about the artifact content.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// An alternate rendered representation of the artifact (e.g., a decompiled
	// representation of a binary region).
	Rendered *MultiformatMessageString `json:"rendered,omitempty" yaml:"rendered,omitempty" mapstructure:"rendered,omitempty"`

	// UTF-8-encoded content from a text artifact.
	Text *string `json:"text,omitempty" yaml:"text,omitempty" mapstructure:"text,omitempty"`
}

// Specifies the location of an artifact.
type ArtifactLocation struct {
	// A short description of the artifact location.
	Description Message `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// The index within the run artifacts array of the artifact object associated with
	// the artifact location.
	Index int `json:"index,omitempty" yaml:"index,omitempty" mapstructure:"index,omitempty"`

	// Key/value pairs that provide additional information about the artifact
	// location.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// A string containing a valid relative or absolute URI.
	URI *string `json:"uri,omitempty" yaml:"uri,omitempty" mapstructure:"uri,omitempty"`

	// A string which indirectly specifies the absolute URI with respect to which a
	// relative URI in the "uri" property is interpreted.
	URIBaseID *string `json:"uriBaseID,omitempty" yaml:"uriBaseID,omitempty" mapstructure:"uriBaseID,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ArtifactLocation) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ArtifactLocation
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["index"]; !ok || v == nil {
		plain.Index = -1.0
	}
	if -1 > plain.Index {
		return fmt.Errorf("field %s: must be >= %v", "index", -1)
	}
	*j = ArtifactLocation(plain)
	return nil
}

type ArtifactRolesElem string

const ArtifactRolesElemAdded ArtifactRolesElem = "added"
const ArtifactRolesElemAnalysisTarget ArtifactRolesElem = "analysisTarget"
const ArtifactRolesElemAttachment ArtifactRolesElem = "attachment"
const ArtifactRolesElemDebugOutputFile ArtifactRolesElem = "debugOutputFile"
const ArtifactRolesElemDeleted ArtifactRolesElem = "deleted"
const ArtifactRolesElemDirectory ArtifactRolesElem = "directory"
const ArtifactRolesElemDriver ArtifactRolesElem = "driver"
const ArtifactRolesElemExtension ArtifactRolesElem = "extension"
const ArtifactRolesElemMemoryContents ArtifactRolesElem = "memoryContents"
const ArtifactRolesElemModified ArtifactRolesElem = "modified"
const ArtifactRolesElemPolicy ArtifactRolesElem = "policy"
const ArtifactRolesElemReferencedOnCommandLine ArtifactRolesElem = "referencedOnCommandLine"
const ArtifactRolesElemRenamed ArtifactRolesElem = "renamed"
const ArtifactRolesElemResponseFile ArtifactRolesElem = "responseFile"
const ArtifactRolesElemResultFile ArtifactRolesElem = "resultFile"
const ArtifactRolesElemStandardStream ArtifactRolesElem = "standardStream"
const ArtifactRolesElemTaxonomy ArtifactRolesElem = "taxonomy"
const ArtifactRolesElemToolSpecifiedConfiguration ArtifactRolesElem = "toolSpecifiedConfiguration"
const ArtifactRolesElemTracedFile ArtifactRolesElem = "tracedFile"
const ArtifactRolesElemTranslation ArtifactRolesElem = "translation"
const ArtifactRolesElemUncontrolled ArtifactRolesElem = "uncontrolled"
const ArtifactRolesElemUnmodified ArtifactRolesElem = "unmodified"
const ArtifactRolesElemUserSpecifiedConfiguration ArtifactRolesElem = "userSpecifiedConfiguration"

var enumValues_ArtifactRolesElem = []interface{}{
	"analysisTarget",
	"attachment",
	"responseFile",
	"resultFile",
	"standardStream",
	"tracedFile",
	"unmodified",
	"modified",
	"added",
	"deleted",
	"renamed",
	"uncontrolled",
	"driver",
	"extension",
	"translation",
	"taxonomy",
	"policy",
	"referencedOnCommandLine",
	"memoryContents",
	"directory",
	"userSpecifiedConfiguration",
	"toolSpecifiedConfiguration",
	"debugOutputFile",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ArtifactRolesElem) UnmarshalJSON(value []byte) error {
	var v string
	if err := json.Unmarshal(value, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_ArtifactRolesElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_ArtifactRolesElem, v)
	}
	*j = ArtifactRolesElem(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Artifact) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain Artifact
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["length"]; !ok || v == nil {
		plain.Length = -1.0
	}
	if -1 > plain.Length {
		return fmt.Errorf("field %s: must be >= %v", "length", -1)
	}
	if plain.MimeType != nil {
		if matched, _ := regexp.MatchString(`[^/]+/.+`, string(*plain.MimeType)); !matched {
			return fmt.Errorf("field %s pattern match: must match %s", "MimeType", `[^/]+/.+`)
		}
	}
	if plain.Offset != nil && 0 > *plain.Offset {
		return fmt.Errorf("field %s: must be >= %v", "offset", 0)
	}
	if v, ok := raw["parentIndex"]; !ok || v == nil {
		plain.ParentIndex = -1.0
	}
	if -1 > plain.ParentIndex {
		return fmt.Errorf("field %s: must be >= %v", "parentIndex", -1)
	}
	if v, ok := raw["roles"]; !ok || v == nil {
		plain.Roles = []ArtifactRolesElem{}
	}
	*j = Artifact(plain)
	return nil
}

// An artifact relevant to a result.
type Attachment struct {
	// The location of the attachment.
	ArtifactLocation ArtifactLocation `json:"artifactLocation" yaml:"artifactLocation" mapstructure:"artifactLocation"`

	// A message describing the role played by the attachment.
	Description Message `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// Key/value pairs that provide additional information about the attachment.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// An array of rectangles specifying areas of interest within the image.
	Rectangles []Rectangle `json:"rectangles,omitempty" yaml:"rectangles,omitempty" mapstructure:"rectangles,omitempty"`

	// An array of regions of interest within the attachment.
	Regions []Region `json:"regions,omitempty" yaml:"regions,omitempty" mapstructure:"regions,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Attachment) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["artifactLocation"]; raw != nil && !ok {
		return fmt.Errorf("field artifactLocation in Attachment: required")
	}
	type Plain Attachment
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["rectangles"]; !ok || v == nil {
		plain.Rectangles = []Rectangle{}
	}
	if v, ok := raw["regions"]; !ok || v == nil {
		plain.Regions = []Region{}
	}
	*j = Attachment(plain)
	return nil
}

// A set of threadFlows which together describe a pattern of code execution
// relevant to detecting a result.
type CodeFlow struct {
	// A message relevant to the code flow.
	Message Message `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message,omitempty"`

	// Key/value pairs that provide additional information about the code flow.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// An array of one or more unique threadFlow objects, each of which describes the
	// progress of a program through a thread of execution.
	ThreadFlows []ThreadFlow `json:"threadFlows" yaml:"threadFlows" mapstructure:"threadFlows"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CodeFlow) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["threadFlows"]; raw != nil && !ok {
		return fmt.Errorf("field threadFlows in CodeFlow: required")
	}
	type Plain CodeFlow
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if plain.ThreadFlows != nil && len(plain.ThreadFlows) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "threadFlows", 1)
	}
	*j = CodeFlow(plain)
	return nil
}

// Information about how a specific rule or notification was reconfigured at
// runtime.
type ConfigurationOverride struct {
	// Specifies how the rule or notification was configured during the scan.
	Configuration ReportingConfiguration `json:"configuration" yaml:"configuration" mapstructure:"configuration"`

	// A reference used to locate the descriptor whose configuration was overridden.
	Descriptor ReportingDescriptorReference `json:"descriptor" yaml:"descriptor" mapstructure:"descriptor"`

	// Key/value pairs that provide additional information about the configuration
	// override.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ConfigurationOverride) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["configuration"]; raw != nil && !ok {
		return fmt.Errorf("field configuration in ConfigurationOverride: required")
	}
	if _, ok := raw["descriptor"]; raw != nil && !ok {
		return fmt.Errorf("field descriptor in ConfigurationOverride: required")
	}
	type Plain ConfigurationOverride
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = ConfigurationOverride(plain)
	return nil
}

// Describes how a converter transformed the output of a static analysis tool from
// the analysis tool's native output format into the SARIF format.
type Conversion struct {
	// The locations of the analysis tool's per-run log files.
	AnalysisToolLogFiles []ArtifactLocation `json:"analysisToolLogFiles,omitempty" yaml:"analysisToolLogFiles,omitempty" mapstructure:"analysisToolLogFiles,omitempty"`

	// An invocation object that describes the invocation of the converter.
	Invocation *Invocation `json:"invocation,omitempty" yaml:"invocation,omitempty" mapstructure:"invocation,omitempty"`

	// Key/value pairs that provide additional information about the conversion.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// A tool object that describes the converter.
	Tool Tool `json:"tool" yaml:"tool" mapstructure:"tool"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Conversion) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["tool"]; raw != nil && !ok {
		return fmt.Errorf("field tool in Conversion: required")
	}
	type Plain Conversion
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["analysisToolLogFiles"]; !ok || v == nil {
		plain.AnalysisToolLogFiles = []ArtifactLocation{}
	}
	*j = Conversion(plain)
	return nil
}

// Represents a directed edge in a graph.
type Edge struct {
	// A string that uniquely identifies the edge within its graph.
	ID string `json:"id" yaml:"id" mapstructure:"id"`

	// A short description of the edge.
	Label Message `json:"label,omitempty" yaml:"label,omitempty" mapstructure:"label,omitempty"`

	// Key/value pairs that provide additional information about the edge.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// IDentifies the source node (the node at which the edge starts).
	SourceNodeID string `json:"sourceNodeID" yaml:"sourceNodeID" mapstructure:"sourceNodeID"`

	// IDentifies the target node (the node at which the edge ends).
	TargetNodeID string `json:"targetNodeID" yaml:"targetNodeID" mapstructure:"targetNodeID"`
}

// Represents the traversal of a single edge during a graph traversal.
type EdgeTraversal struct {
	// IDentifies the edge being traversed.
	EdgeID string `json:"edgeID" yaml:"edgeID" mapstructure:"edgeID"`

	// The values of relevant expressions after the edge has been traversed.
	FinalState map[string]MultiformatMessageString `json:"finalState,omitempty" yaml:"finalState,omitempty" mapstructure:"finalState,omitempty"`

	// A message to display to the user as the edge is traversed.
	Message Message `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message,omitempty"`

	// Key/value pairs that provide additional information about the edge traversal.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The number of edge traversals necessary to return from a nested graph.
	StepOverEdgeCount *int `json:"stepOverEdgeCount,omitempty" yaml:"stepOverEdgeCount,omitempty" mapstructure:"stepOverEdgeCount,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *EdgeTraversal) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["edgeID"]; raw != nil && !ok {
		return fmt.Errorf("field edgeID in EdgeTraversal: required")
	}
	type Plain EdgeTraversal
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if plain.StepOverEdgeCount != nil && 0 > *plain.StepOverEdgeCount {
		return fmt.Errorf("field %s: must be >= %v", "stepOverEdgeCount", 0)
	}
	*j = EdgeTraversal(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Edge) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["id"]; raw != nil && !ok {
		return fmt.Errorf("field id in Edge: required")
	}
	if _, ok := raw["sourceNodeID"]; raw != nil && !ok {
		return fmt.Errorf("field sourceNodeID in Edge: required")
	}
	if _, ok := raw["targetNodeID"]; raw != nil && !ok {
		return fmt.Errorf("field targetNodeID in Edge: required")
	}
	type Plain Edge
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = Edge(plain)
	return nil
}

// Describes a runtime exception encountered during the execution of an analysis
// tool.
type Exception struct {
	// An array of exception objects each of which is considered a cause of this
	// exception.
	InnerExceptions []Exception `json:"innerExceptions,omitempty" yaml:"innerExceptions,omitempty" mapstructure:"innerExceptions,omitempty"`

	// A string that identifies the kind of exception, for example, the fully
	// qualified type name of an object that was thrown, or the symbolic name of a
	// signal.
	Kind *string `json:"kind,omitempty" yaml:"kind,omitempty" mapstructure:"kind,omitempty"`

	// A message that describes the exception.
	Message *string `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message,omitempty"`

	// Key/value pairs that provide additional information about the exception.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The sequence of function calls leading to the exception.
	Stack *Stack `json:"stack,omitempty" yaml:"stack,omitempty" mapstructure:"stack,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Exception) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain Exception
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["innerExceptions"]; !ok || v == nil {
		plain.InnerExceptions = []Exception{}
	}
	*j = Exception(plain)
	return nil
}

// The top-level element of an external property file.
type ExternalProperties struct {
	// Addresses that will be merged with a separate run.
	Addresses []Address `json:"addresses,omitempty" yaml:"addresses,omitempty" mapstructure:"addresses,omitempty"`

	// An array of artifact objects that will be merged with a separate run.
	Artifacts []Artifact `json:"artifacts,omitempty" yaml:"artifacts,omitempty" mapstructure:"artifacts,omitempty"`

	// A conversion object that will be merged with a separate run.
	Conversion *Conversion `json:"conversion,omitempty" yaml:"conversion,omitempty" mapstructure:"conversion,omitempty"`

	// The analysis tool object that will be merged with a separate run.
	Driver *ToolComponent `json:"driver,omitempty" yaml:"driver,omitempty" mapstructure:"driver,omitempty"`

	// Tool extensions that will be merged with a separate run.
	Extensions []ToolComponent `json:"extensions,omitempty" yaml:"extensions,omitempty" mapstructure:"extensions,omitempty"`

	// Key/value pairs that provide additional information that will be merged with a
	// separate run.
	ExternalizedProperties *PropertyBag `json:"externalizedProperties,omitempty" yaml:"externalizedProperties,omitempty" mapstructure:"externalizedProperties,omitempty"`

	// An array of graph objects that will be merged with a separate run.
	Graphs []Graph `json:"graphs,omitempty" yaml:"graphs,omitempty" mapstructure:"graphs,omitempty"`

	// A stable, unique identifier for this external properties object, in the form of
	// a GUID.
	Guid *string `json:"guid,omitempty" yaml:"guid,omitempty" mapstructure:"guid,omitempty"`

	// Describes the invocation of the analysis tool that will be merged with a
	// separate run.
	Invocations []Invocation `json:"invocations,omitempty" yaml:"invocations,omitempty" mapstructure:"invocations,omitempty"`

	// An array of logical locations such as namespaces, types or functions that will
	// be merged with a separate run.
	LogicalLocations []LogicalLocation `json:"logicalLocations,omitempty" yaml:"logicalLocations,omitempty" mapstructure:"logicalLocations,omitempty"`

	// Tool policies that will be merged with a separate run.
	Policies []ToolComponent `json:"policies,omitempty" yaml:"policies,omitempty" mapstructure:"policies,omitempty"`

	// Key/value pairs that provide additional information about the external
	// properties.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// An array of result objects that will be merged with a separate run.
	Results []Result `json:"results,omitempty" yaml:"results,omitempty" mapstructure:"results,omitempty"`

	// A stable, unique identifier for the run associated with this external
	// properties object, in the form of a GUID.
	RunGuid *string `json:"runGuid,omitempty" yaml:"runGuid,omitempty" mapstructure:"runGuid,omitempty"`

	// The URI of the JSON schema corresponding to the version of the external
	// property file format.
	Schema *string `json:"schema,omitempty" yaml:"schema,omitempty" mapstructure:"schema,omitempty"`

	// Tool taxonomies that will be merged with a separate run.
	Taxonomies []ToolComponent `json:"taxonomies,omitempty" yaml:"taxonomies,omitempty" mapstructure:"taxonomies,omitempty"`

	// An array of threadFlowLocation objects that will be merged with a separate run.
	ThreadFlowLocations []ThreadFlowLocation `json:"threadFlowLocations,omitempty" yaml:"threadFlowLocations,omitempty" mapstructure:"threadFlowLocations,omitempty"`

	// Tool translations that will be merged with a separate run.
	Translations []ToolComponent `json:"translations,omitempty" yaml:"translations,omitempty" mapstructure:"translations,omitempty"`

	// The SARIF format version of this external properties object.
	Version *ExternalPropertiesVersion `json:"version,omitempty" yaml:"version,omitempty" mapstructure:"version,omitempty"`

	// Requests that will be merged with a separate run.
	WebRequests []WebRequest `json:"webRequests,omitempty" yaml:"webRequests,omitempty" mapstructure:"webRequests,omitempty"`

	// Responses that will be merged with a separate run.
	WebResponses []WebResponse `json:"webResponses,omitempty" yaml:"webResponses,omitempty" mapstructure:"webResponses,omitempty"`
}

type ExternalPropertiesVersion string

const ExternalPropertiesVersionA210 ExternalPropertiesVersion = "2.1.0"

var enumValues_ExternalPropertiesVersion = []interface{}{
	"2.1.0",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExternalPropertiesVersion) UnmarshalJSON(value []byte) error {
	var v string
	if err := json.Unmarshal(value, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_ExternalPropertiesVersion {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_ExternalPropertiesVersion, v)
	}
	*j = ExternalPropertiesVersion(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExternalProperties) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ExternalProperties
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["addresses"]; !ok || v == nil {
		plain.Addresses = []Address{}
	}
	if v, ok := raw["extensions"]; !ok || v == nil {
		plain.Extensions = []ToolComponent{}
	}
	if v, ok := raw["graphs"]; !ok || v == nil {
		plain.Graphs = []Graph{}
	}
	if plain.Guid != nil {
		if matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`, string(*plain.Guid)); !matched {
			return fmt.Errorf("field %s pattern match: must match %s", "Guid", `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
		}
	}
	if v, ok := raw["invocations"]; !ok || v == nil {
		plain.Invocations = []Invocation{}
	}
	if v, ok := raw["logicalLocations"]; !ok || v == nil {
		plain.LogicalLocations = []LogicalLocation{}
	}
	if v, ok := raw["policies"]; !ok || v == nil {
		plain.Policies = []ToolComponent{}
	}
	if v, ok := raw["results"]; !ok || v == nil {
		plain.Results = []Result{}
	}
	if plain.RunGuid != nil {
		if matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`, string(*plain.RunGuid)); !matched {
			return fmt.Errorf("field %s pattern match: must match %s", "RunGuid", `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
		}
	}
	if v, ok := raw["taxonomies"]; !ok || v == nil {
		plain.Taxonomies = []ToolComponent{}
	}
	if v, ok := raw["threadFlowLocations"]; !ok || v == nil {
		plain.ThreadFlowLocations = []ThreadFlowLocation{}
	}
	if v, ok := raw["translations"]; !ok || v == nil {
		plain.Translations = []ToolComponent{}
	}
	if v, ok := raw["webRequests"]; !ok || v == nil {
		plain.WebRequests = []WebRequest{}
	}
	if v, ok := raw["webResponses"]; !ok || v == nil {
		plain.WebResponses = []WebResponse{}
	}
	*j = ExternalProperties(plain)
	return nil
}

// Contains information that enables a SARIF consumer to locate the external
// property file that contains the value of an externalized property associated
// with the run.
type ExternalPropertyFileReference interface{}

// References to external property files that should be inlined with the content of
// a root log file.
type ExternalPropertyFileReferences struct {
	// An array of external property files containing run.addresses arrays to be
	// merged with the root log file.
	Addresses []ExternalPropertyFileReference `json:"addresses,omitempty" yaml:"addresses,omitempty" mapstructure:"addresses,omitempty"`

	// An array of external property files containing run.artifacts arrays to be
	// merged with the root log file.
	Artifacts []ExternalPropertyFileReference `json:"artifacts,omitempty" yaml:"artifacts,omitempty" mapstructure:"artifacts,omitempty"`

	// An external property file containing a run.conversion object to be merged with
	// the root log file.
	Conversion ExternalPropertyFileReference `json:"conversion,omitempty" yaml:"conversion,omitempty" mapstructure:"conversion,omitempty"`

	// An external property file containing a run.driver object to be merged with the
	// root log file.
	Driver ExternalPropertyFileReference `json:"driver,omitempty" yaml:"driver,omitempty" mapstructure:"driver,omitempty"`

	// An array of external property files containing run.extensions arrays to be
	// merged with the root log file.
	Extensions []ExternalPropertyFileReference `json:"extensions,omitempty" yaml:"extensions,omitempty" mapstructure:"extensions,omitempty"`

	// An external property file containing a run.properties object to be merged with
	// the root log file.
	ExternalizedProperties ExternalPropertyFileReference `json:"externalizedProperties,omitempty" yaml:"externalizedProperties,omitempty" mapstructure:"externalizedProperties,omitempty"`

	// An array of external property files containing a run.graphs object to be merged
	// with the root log file.
	Graphs []ExternalPropertyFileReference `json:"graphs,omitempty" yaml:"graphs,omitempty" mapstructure:"graphs,omitempty"`

	// An array of external property files containing run.invocations arrays to be
	// merged with the root log file.
	Invocations []ExternalPropertyFileReference `json:"invocations,omitempty" yaml:"invocations,omitempty" mapstructure:"invocations,omitempty"`

	// An array of external property files containing run.logicalLocations arrays to
	// be merged with the root log file.
	LogicalLocations []ExternalPropertyFileReference `json:"logicalLocations,omitempty" yaml:"logicalLocations,omitempty" mapstructure:"logicalLocations,omitempty"`

	// An array of external property files containing run.policies arrays to be merged
	// with the root log file.
	Policies []ExternalPropertyFileReference `json:"policies,omitempty" yaml:"policies,omitempty" mapstructure:"policies,omitempty"`

	// Key/value pairs that provide additional information about the external property
	// files.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// An array of external property files containing run.results arrays to be merged
	// with the root log file.
	Results []ExternalPropertyFileReference `json:"results,omitempty" yaml:"results,omitempty" mapstructure:"results,omitempty"`

	// An array of external property files containing run.taxonomies arrays to be
	// merged with the root log file.
	Taxonomies []ExternalPropertyFileReference `json:"taxonomies,omitempty" yaml:"taxonomies,omitempty" mapstructure:"taxonomies,omitempty"`

	// An array of external property files containing run.threadFlowLocations arrays
	// to be merged with the root log file.
	ThreadFlowLocations []ExternalPropertyFileReference `json:"threadFlowLocations,omitempty" yaml:"threadFlowLocations,omitempty" mapstructure:"threadFlowLocations,omitempty"`

	// An array of external property files containing run.translations arrays to be
	// merged with the root log file.
	Translations []ExternalPropertyFileReference `json:"translations,omitempty" yaml:"translations,omitempty" mapstructure:"translations,omitempty"`

	// An array of external property files containing run.requests arrays to be merged
	// with the root log file.
	WebRequests []ExternalPropertyFileReference `json:"webRequests,omitempty" yaml:"webRequests,omitempty" mapstructure:"webRequests,omitempty"`

	// An array of external property files containing run.responses arrays to be
	// merged with the root log file.
	WebResponses []ExternalPropertyFileReference `json:"webResponses,omitempty" yaml:"webResponses,omitempty" mapstructure:"webResponses,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ExternalPropertyFileReferences) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ExternalPropertyFileReferences
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["addresses"]; !ok || v == nil {
		plain.Addresses = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["artifacts"]; !ok || v == nil {
		plain.Artifacts = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["extensions"]; !ok || v == nil {
		plain.Extensions = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["graphs"]; !ok || v == nil {
		plain.Graphs = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["invocations"]; !ok || v == nil {
		plain.Invocations = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["logicalLocations"]; !ok || v == nil {
		plain.LogicalLocations = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["policies"]; !ok || v == nil {
		plain.Policies = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["results"]; !ok || v == nil {
		plain.Results = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["taxonomies"]; !ok || v == nil {
		plain.Taxonomies = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["threadFlowLocations"]; !ok || v == nil {
		plain.ThreadFlowLocations = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["translations"]; !ok || v == nil {
		plain.Translations = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["webRequests"]; !ok || v == nil {
		plain.WebRequests = []ExternalPropertyFileReference{}
	}
	if v, ok := raw["webResponses"]; !ok || v == nil {
		plain.WebResponses = []ExternalPropertyFileReference{}
	}
	*j = ExternalPropertyFileReferences(plain)
	return nil
}

// A proposed fix for the problem represented by a result object. A fix specifies a
// set of artifacts to modify. For each artifact, it specifies a set of bytes to
// remove, and provides a set of new bytes to replace them.
type Fix struct {
	// One or more artifact changes that comprise a fix for a result.
	ArtifactChanges []ArtifactChange `json:"artifactChanges" yaml:"artifactChanges" mapstructure:"artifactChanges"`

	// A message that describes the proposed fix, enabling viewers to present the
	// proposed change to an end user.
	Description Message `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// Key/value pairs that provide additional information about the fix.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Fix) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["artifactChanges"]; raw != nil && !ok {
		return fmt.Errorf("field artifactChanges in Fix: required")
	}
	type Plain Fix
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if plain.ArtifactChanges != nil && len(plain.ArtifactChanges) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "artifactChanges", 1)
	}
	*j = Fix(plain)
	return nil
}

// A network of nodes and directed edges that describes some aspect of the
// structure of the code (for example, a call graph).
type Graph struct {
	// A description of the graph.
	Description Message `json:"description,omitempty" yaml:"description,omitempty" mapstructure:"description,omitempty"`

	// An array of edge objects representing the edges of the graph.
	Edges []Edge `json:"edges,omitempty" yaml:"edges,omitempty" mapstructure:"edges,omitempty"`

	// An array of node objects representing the nodes of the graph.
	Nodes []Node `json:"nodes,omitempty" yaml:"nodes,omitempty" mapstructure:"nodes,omitempty"`

	// Key/value pairs that provide additional information about the graph.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`
}

// Represents a path through a graph.
type GraphTraversal struct {
	// A description of this graph traversal.
	Description Message `json:"description,omitempty" yaml:"description,omitemp…15325 tokens truncated…k := raw["tool"]; raw != nil && !ok {
		return fmt.Errorf("field tool in Run: required")
	}
	type Plain Run
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["addresses"]; !ok || v == nil {
		plain.Addresses = []Address{}
	}
	if plain.BaselineGuid != nil {
		if matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`, string(*plain.BaselineGuid)); !matched {
			return fmt.Errorf("field %s pattern match: must match %s", "BaselineGuid", `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
		}
	}
	if v, ok := raw["graphs"]; !ok || v == nil {
		plain.Graphs = []Graph{}
	}
	if v, ok := raw["invocations"]; !ok || v == nil {
		plain.Invocations = []Invocation{}
	}
	if v, ok := raw["language"]; !ok || v == nil {
		plain.Language = "en-US"
	}
	if matched, _ := regexp.MatchString(`^[a-zA-Z]{2}(-[a-zA-Z]{2})?$`, string(plain.Language)); !matched {
		return fmt.Errorf("field %s pattern match: must match %s", "Language", `^[a-zA-Z]{2}(-[a-zA-Z]{2})?$`)
	}
	if v, ok := raw["logicalLocations"]; !ok || v == nil {
		plain.LogicalLocations = []LogicalLocation{}
	}
	if v, ok := raw["newlineSequences"]; !ok || v == nil {
		plain.NewlineSequences = []string{
			"\r\n",
			"\n",
		}
	}
	if plain.NewlineSequences != nil && len(plain.NewlineSequences) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "newlineSequences", 1)
	}
	if v, ok := raw["policies"]; !ok || v == nil {
		plain.Policies = []ToolComponent{}
	}
	if v, ok := raw["redactionTokens"]; !ok || v == nil {
		plain.RedactionTokens = []string{}
	}
	if v, ok := raw["runAggregates"]; !ok || v == nil {
		plain.RunAggregates = []RunAutomationDetails{}
	}
	if v, ok := raw["taxonomies"]; !ok || v == nil {
		plain.Taxonomies = []ToolComponent{}
	}
	if v, ok := raw["threadFlowLocations"]; !ok || v == nil {
		plain.ThreadFlowLocations = []ThreadFlowLocation{}
	}
	if v, ok := raw["translations"]; !ok || v == nil {
		plain.Translations = []ToolComponent{}
	}
	if v, ok := raw["versionControlProvenance"]; !ok || v == nil {
		plain.VersionControlProvenance = []VersionControlDetails{}
	}
	if v, ok := raw["webRequests"]; !ok || v == nil {
		plain.WebRequests = []WebRequest{}
	}
	if v, ok := raw["webResponses"]; !ok || v == nil {
		plain.WebResponses = []WebResponse{}
	}
	*j = Run(plain)
	return nil
}

// Static Analysis Results Format (SARIF) Version 2.1.0 JSON Schema: a standard
// format for the output of static analysis tools.
type SarifSchema210Json struct {
	// The URI of the JSON schema corresponding to the version.
	Schema *string `json:"$schema,omitempty" yaml:"$schema,omitempty" mapstructure:"$schema,omitempty"`

	// References to external property files that share data between runs.
	InlineExternalProperties []ExternalProperties `json:"inlineExternalProperties,omitempty" yaml:"inlineExternalProperties,omitempty" mapstructure:"inlineExternalProperties,omitempty"`

	// Key/value pairs that provide additional information about the log file.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The set of runs contained in this log file.
	Runs []Run `json:"runs" yaml:"runs" mapstructure:"runs"`

	// The SARIF format version of this log file.
	Version SarifSchema210JsonVersion `json:"version" yaml:"version" mapstructure:"version"`
}

type SarifSchema210JsonVersion string

const SarifSchema210JsonVersionA210 SarifSchema210JsonVersion = "2.1.0"

var enumValues_SarifSchema210JsonVersion = []interface{}{
	"2.1.0",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SarifSchema210JsonVersion) UnmarshalJSON(value []byte) error {
	var v string
	if err := json.Unmarshal(value, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_SarifSchema210JsonVersion {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_SarifSchema210JsonVersion, v)
	}
	*j = SarifSchema210JsonVersion(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SarifSchema210Json) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["runs"]; raw != nil && !ok {
		return fmt.Errorf("field runs in SarifSchema210Json: required")
	}
	if _, ok := raw["version"]; raw != nil && !ok {
		return fmt.Errorf("field version in SarifSchema210Json: required")
	}
	type Plain SarifSchema210Json
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = SarifSchema210Json(plain)
	return nil
}

// Defines locations of special significance to SARIF consumers.
type SpecialLocations struct {
	// Provides a suggestion to SARIF consumers to display file paths relative to the
	// specified location.
	DisplayBase *ArtifactLocation `json:"displayBase,omitempty" yaml:"displayBase,omitempty" mapstructure:"displayBase,omitempty"`

	// Key/value pairs that provide additional information about the special
	// locations.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`
}

// A call stack that is relevant to a result.
type Stack struct {
	// An array of stack frames that represents a sequence of calls, rendered in
	// reverse chronological order, that comprise the call stack.
	Frames []StackFrame `json:"frames" yaml:"frames" mapstructure:"frames"`

	// A message relevant to this call stack.
	Message Message `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message,omitempty"`

	// Key/value pairs that provide additional information about the stack.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`
}

// A function call within a stack trace.
type StackFrame struct {
	// The location to which this stack frame refers.
	Location *Location `json:"location,omitempty" yaml:"location,omitempty" mapstructure:"location,omitempty"`

	// The name of the module that contains the code of this stack frame.
	Module *string `json:"module,omitempty" yaml:"module,omitempty" mapstructure:"module,omitempty"`

	// The parameters of the call that is executing.
	Parameters []string `json:"parameters,omitempty" yaml:"parameters,omitempty" mapstructure:"parameters,omitempty"`

	// Key/value pairs that provide additional information about the stack frame.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The thread identifier of the stack frame.
	ThreadID *int `json:"threadID,omitempty" yaml:"threadID,omitempty" mapstructure:"threadID,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *StackFrame) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain StackFrame
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["parameters"]; !ok || v == nil {
		plain.Parameters = []string{}
	}
	*j = StackFrame(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Stack) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["frames"]; raw != nil && !ok {
		return fmt.Errorf("field frames in Stack: required")
	}
	type Plain Stack
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = Stack(plain)
	return nil
}

// A suppression that is relevant to a result.
type Suppression struct {
	// A stable, unique identifier for the suprression in the form of a GUID.
	Guid *string `json:"guid,omitempty" yaml:"guid,omitempty" mapstructure:"guid,omitempty"`

	// A string representing the justification for the suppression.
	Justification *string `json:"justification,omitempty" yaml:"justification,omitempty" mapstructure:"justification,omitempty"`

	// A string that indicates where the suppression is persisted.
	Kind SuppressionKind `json:"kind" yaml:"kind" mapstructure:"kind"`

	// IDentifies the location associated with the suppression.
	Location *Location `json:"location,omitempty" yaml:"location,omitempty" mapstructure:"location,omitempty"`

	// Key/value pairs that provide additional information about the suppression.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// A string that indicates the review status of the suppression.
	Status *SuppressionStatus `json:"status,omitempty" yaml:"status,omitempty" mapstructure:"status,omitempty"`
}

type SuppressionKind string

const SuppressionKindExternal SuppressionKind = "external"
const SuppressionKindInSource SuppressionKind = "inSource"

var enumValues_SuppressionKind = []interface{}{
	"inSource",
	"external",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SuppressionKind) UnmarshalJSON(value []byte) error {
	var v string
	if err := json.Unmarshal(value, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_SuppressionKind {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_SuppressionKind, v)
	}
	*j = SuppressionKind(v)
	return nil
}

type SuppressionStatus string

const SuppressionStatusAccepted SuppressionStatus = "accepted"
const SuppressionStatusRejected SuppressionStatus = "rejected"
const SuppressionStatusUnderReview SuppressionStatus = "underReview"

var enumValues_SuppressionStatus = []interface{}{
	"accepted",
	"underReview",
	"rejected",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SuppressionStatus) UnmarshalJSON(value []byte) error {
	var v string
	if err := json.Unmarshal(value, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_SuppressionStatus {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_SuppressionStatus, v)
	}
	*j = SuppressionStatus(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Suppression) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["kind"]; raw != nil && !ok {
		return fmt.Errorf("field kind in Suppression: required")
	}
	type Plain Suppression
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if plain.Guid != nil {
		if matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`, string(*plain.Guid)); !matched {
			return fmt.Errorf("field %s pattern match: must match %s", "Guid", `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
		}
	}
	*j = Suppression(plain)
	return nil
}

// Describes a sequence of code locations that specify a path through a single
// thread of execution such as an operating system or fiber.
type ThreadFlow struct {
	// An string that uniquely identifies the threadFlow within the codeFlow in which
	// it occurs.
	ID *string `json:"id,omitempty" yaml:"id,omitempty" mapstructure:"id,omitempty"`

	// Values of relevant expressions at the start of the thread flow that remain
	// constant.
	ImmutableState map[string]MultiformatMessageString `json:"immutableState,omitempty" yaml:"immutableState,omitempty" mapstructure:"immutableState,omitempty"`

	// Values of relevant expressions at the start of the thread flow that may change
	// during thread flow execution.
	InitialState map[string]MultiformatMessageString `json:"initialState,omitempty" yaml:"initialState,omitempty" mapstructure:"initialState,omitempty"`

	// A temporally ordered array of 'threadFlowLocation' objects, each of which
	// describes a location visited by the tool while producing the result.
	Locations []ThreadFlowLocation `json:"locations" yaml:"locations" mapstructure:"locations"`

	// A message relevant to the thread flow.
	Message Message `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message,omitempty"`

	// Key/value pairs that provide additional information about the thread flow.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`
}

// A location visited by an analysis tool while simulating or monitoring the
// execution of a program.
type ThreadFlowLocation struct {
	// An integer representing the temporal order in which execution reached this
	// location.
	ExecutionOrder int `json:"executionOrder,omitempty" yaml:"executionOrder,omitempty" mapstructure:"executionOrder,omitempty"`

	// The Coordinated Universal Time (UTC) date and time at which this location was
	// executed.
	ExecutionTimeUtc *time.Time `json:"executionTimeUtc,omitempty" yaml:"executionTimeUtc,omitempty" mapstructure:"executionTimeUtc,omitempty"`

	// Specifies the importance of this location in understanding the code flow in
	// which it occurs. The order from most to least important is "essential",
	// "important", "unimportant". Default: "important".
	Importance ThreadFlowLocationImportance `json:"importance,omitempty" yaml:"importance,omitempty" mapstructure:"importance,omitempty"`

	// The index within the run threadFlowLocations array.
	Index int `json:"index,omitempty" yaml:"index,omitempty" mapstructure:"index,omitempty"`

	// A set of distinct strings that categorize the thread flow location. Well-known
	// kinds include 'acquire', 'release', 'enter', 'exit', 'call', 'return',
	// 'branch', 'implicit', 'false', 'true', 'caution', 'danger', 'unknown',
	// 'unreachable', 'taint', 'function', 'handler', 'lock', 'memory', 'resource',
	// 'scope' and 'value'.
	Kinds []string `json:"kinds,omitempty" yaml:"kinds,omitempty" mapstructure:"kinds,omitempty"`

	// The code location.
	Location *Location `json:"location,omitempty" yaml:"location,omitempty" mapstructure:"location,omitempty"`

	// The name of the module that contains the code that is executing.
	Module *string `json:"module,omitempty" yaml:"module,omitempty" mapstructure:"module,omitempty"`

	// An integer representing a containment hierarchy within the thread flow.
	NestingLevel *int `json:"nestingLevel,omitempty" yaml:"nestingLevel,omitempty" mapstructure:"nestingLevel,omitempty"`

	// Key/value pairs that provide additional information about the threadflow
	// location.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The call stack leading to this location.
	Stack *Stack `json:"stack,omitempty" yaml:"stack,omitempty" mapstructure:"stack,omitempty"`

	// A dictionary, each of whose keys specifies a variable or expression, the
	// associated value of which represents the variable or expression value. For an
	// annotation of kind 'continuation', for example, this dictionary might hold the
	// current assumed values of a set of global variables.
	State map[string]MultiformatMessageString `json:"state,omitempty" yaml:"state,omitempty" mapstructure:"state,omitempty"`

	// An array of references to rule or taxonomy reporting descriptors that are
	// applicable to the thread flow location.
	Taxa []ReportingDescriptorReference `json:"taxa,omitempty" yaml:"taxa,omitempty" mapstructure:"taxa,omitempty"`

	// A web request associated with this thread flow location.
	WebRequest *WebRequest `json:"webRequest,omitempty" yaml:"webRequest,omitempty" mapstructure:"webRequest,omitempty"`

	// A web response associated with this thread flow location.
	WebResponse *WebResponse `json:"webResponse,omitempty" yaml:"webResponse,omitempty" mapstructure:"webResponse,omitempty"`
}

type ThreadFlowLocationImportance string

const ThreadFlowLocationImportanceEssential ThreadFlowLocationImportance = "essential"
const ThreadFlowLocationImportanceImportant ThreadFlowLocationImportance = "important"
const ThreadFlowLocationImportanceUnimportant ThreadFlowLocationImportance = "unimportant"

var enumValues_ThreadFlowLocationImportance = []interface{}{
	"important",
	"essential",
	"unimportant",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ThreadFlowLocationImportance) UnmarshalJSON(value []byte) error {
	var v string
	if err := json.Unmarshal(value, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_ThreadFlowLocationImportance {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_ThreadFlowLocationImportance, v)
	}
	*j = ThreadFlowLocationImportance(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ThreadFlowLocation) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ThreadFlowLocation
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["executionOrder"]; !ok || v == nil {
		plain.ExecutionOrder = -1.0
	}
	if -1 > plain.ExecutionOrder {
		return fmt.Errorf("field %s: must be >= %v", "executionOrder", -1)
	}
	if v, ok := raw["importance"]; !ok || v == nil {
		plain.Importance = "important"
	}
	if v, ok := raw["index"]; !ok || v == nil {
		plain.Index = -1.0
	}
	if -1 > plain.Index {
		return fmt.Errorf("field %s: must be >= %v", "index", -1)
	}
	if v, ok := raw["kinds"]; !ok || v == nil {
		plain.Kinds = []string{}
	}
	if plain.NestingLevel != nil && 0 > *plain.NestingLevel {
		return fmt.Errorf("field %s: must be >= %v", "nestingLevel", 0)
	}
	if v, ok := raw["taxa"]; !ok || v == nil {
		plain.Taxa = []ReportingDescriptorReference{}
	}
	*j = ThreadFlowLocation(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ThreadFlow) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["locations"]; raw != nil && !ok {
		return fmt.Errorf("field locations in ThreadFlow: required")
	}
	type Plain ThreadFlow
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if plain.Locations != nil && len(plain.Locations) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "locations", 1)
	}
	*j = ThreadFlow(plain)
	return nil
}

// The analysis tool that was run.
type Tool struct {
	// The analysis tool that was run.
	Driver ToolComponent `json:"driver" yaml:"driver" mapstructure:"driver"`

	// Tool extensions that contributed to or reconfigured the analysis tool that was
	// run.
	Extensions []ToolComponent `json:"extensions,omitempty" yaml:"extensions,omitempty" mapstructure:"extensions,omitempty"`

	// Key/value pairs that provide additional information about the tool.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`
}

// A component, such as a plug-in or the driver, of the analysis tool that was run.
type ToolComponent struct {
	// The component which is strongly associated with this component. For a
	// translation, this refers to the component which has been translated. For an
	// extension, this is the driver that provides the extension's plugin model.
	AssociatedComponent *ToolComponentReference `json:"associatedComponent,omitempty" yaml:"associatedComponent,omitempty" mapstructure:"associatedComponent,omitempty"`

	// The kinds of data contained in this object.
	Contents []ToolComponentContentsElem `json:"contents,omitempty" yaml:"contents,omitempty" mapstructure:"contents,omitempty"`

	// The binary version of the tool component's primary executable file expressed as
	// four non-negative integers separated by a period (for operating systems that
	// express file versions in this way).
	DottedQuadFileVersion *string `json:"dottedQuadFileVersion,omitempty" yaml:"dottedQuadFileVersion,omitempty" mapstructure:"dottedQuadFileVersion,omitempty"`

	// The absolute URI from which the tool component can be downloaded.
	DownloadURI *string `json:"downloadURI,omitempty" yaml:"downloadURI,omitempty" mapstructure:"downloadURI,omitempty"`

	// A comprehensive description of the tool component.
	FullDescription *MultiformatMessageString `json:"fullDescription,omitempty" yaml:"fullDescription,omitempty" mapstructure:"fullDescription,omitempty"`

	// The name of the tool component along with its version and any other useful
	// identifying information, such as its locale.
	FullName *string `json:"fullName,omitempty" yaml:"fullName,omitempty" mapstructure:"fullName,omitempty"`

	// A dictionary, each of whose keys is a resource identifier and each of whose
	// values is a multiformatMessageString object, which holds message strings in
	// plain text and (optionally) Markdown format. The strings can include
	// placeholders, which can be used to construct a message in combination with an
	// arbitrary number of additional string arguments.
	GlobalMessageStrings map[string]MultiformatMessageString `json:"globalMessageStrings,omitempty" yaml:"globalMessageStrings,omitempty" mapstructure:"globalMessageStrings,omitempty"`

	// A unique identifier for the tool component in the form of a GUID.
	Guid *string `json:"guid,omitempty" yaml:"guid,omitempty" mapstructure:"guid,omitempty"`

	// The absolute URI at which information about this version of the tool component
	// can be found.
	InformationURI *string `json:"informationURI,omitempty" yaml:"informationURI,omitempty" mapstructure:"informationURI,omitempty"`

	// Specifies whether this object contains a complete definition of the localizable
	// and/or non-localizable data for this component, as opposed to including only
	// data that is relevant to the results persisted to this log file.
	IsComprehensive bool `json:"isComprehensive,omitempty" yaml:"isComprehensive,omitempty" mapstructure:"isComprehensive,omitempty"`

	// The language of the messages emitted into the log file during this run
	// (expressed as an ISO 639-1 two-letter lowercase language code) and an optional
	// region (expressed as an ISO 3166-1 two-letter uppercase subculture code
	// associated with a country or region). The casing is recommended but not
	// required (in order for this data to conform to RFC5646).
	Language string `json:"language,omitempty" yaml:"language,omitempty" mapstructure:"language,omitempty"`

	// The semantic version of the localized strings defined in this component;
	// maintained by components that provide translations.
	LocalizedDataSemanticVersion *string `json:"localizedDataSemanticVersion,omitempty" yaml:"localizedDataSemanticVersion,omitempty" mapstructure:"localizedDataSemanticVersion,omitempty"`

	// An array of the artifactLocation objects associated with the tool component.
	Locations []ArtifactLocation `json:"locations,omitempty" yaml:"locations,omitempty" mapstructure:"locations,omitempty"`

	// The minimum value of localizedDataSemanticVersion required in translations
	// consumed by this component; used by components that consume translations.
	MinimumRequiredLocalizedDataSemanticVersion *string `json:"minimumRequiredLocalizedDataSemanticVersion,omitempty" yaml:"minimumRequiredLocalizedDataSemanticVersion,omitempty" mapstructure:"minimumRequiredLocalizedDataSemanticVersion,omitempty"`

	// The name of the tool component.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// An array of reportingDescriptor objects relevant to the notifications related
	// to the configuration and runtime execution of the tool component.
	Notifications []ReportingDescriptor `json:"notifications,omitempty" yaml:"notifications,omitempty" mapstructure:"notifications,omitempty"`

	// The organization or company that produced the tool component.
	Organization *string `json:"organization,omitempty" yaml:"organization,omitempty" mapstructure:"organization,omitempty"`

	// A product suite to which the tool component belongs.
	Product *string `json:"product,omitempty" yaml:"product,omitempty" mapstructure:"product,omitempty"`

	// A localizable string containing the name of the suite of products to which the
	// tool component belongs.
	ProductSuite *string `json:"productSuite,omitempty" yaml:"productSuite,omitempty" mapstructure:"productSuite,omitempty"`

	// Key/value pairs that provide additional information about the tool component.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// A string specifying the UTC date (and optionally, the time) of the component's
	// release.
	ReleaseDateUtc *string `json:"releaseDateUtc,omitempty" yaml:"releaseDateUtc,omitempty" mapstructure:"releaseDateUtc,omitempty"`

	// An array of reportingDescriptor objects relevant to the analysis performed by
	// the tool component.
	Rules []ReportingDescriptor `json:"rules,omitempty" yaml:"rules,omitempty" mapstructure:"rules,omitempty"`

	// The tool component version in the format specified by Semantic Versioning 2.0.
	SemanticVersion *string `json:"semanticVersion,omitempty" yaml:"semanticVersion,omitempty" mapstructure:"semanticVersion,omitempty"`

	// A brief description of the tool component.
	ShortDescription *MultiformatMessageString `json:"shortDescription,omitempty" yaml:"shortDescription,omitempty" mapstructure:"shortDescription,omitempty"`

	// An array of toolComponentReference objects to declare the taxonomies supported
	// by the tool component.
	SupportedTaxonomies []ToolComponentReference `json:"supportedTaxonomies,omitempty" yaml:"supportedTaxonomies,omitempty" mapstructure:"supportedTaxonomies,omitempty"`

	// An array of reportingDescriptor objects relevant to the definitions of both
	// standalone and tool-defined taxonomies.
	Taxa []ReportingDescriptor `json:"taxa,omitempty" yaml:"taxa,omitempty" mapstructure:"taxa,omitempty"`

	// Translation metadata, required for a translation, not populated by other
	// component types.
	TranslationMetadata *TranslationMetadata `json:"translationMetadata,omitempty" yaml:"translationMetadata,omitempty" mapstructure:"translationMetadata,omitempty"`

	// The tool component version, in whatever format the component natively provides.
	Version *string `json:"version,omitempty" yaml:"version,omitempty" mapstructure:"version,omitempty"`
}

type ToolComponentContentsElem string

const ToolComponentContentsElemLocalizedData ToolComponentContentsElem = "localizedData"
const ToolComponentContentsElemNonLocalizedData ToolComponentContentsElem = "nonLocalizedData"

var enumValues_ToolComponentContentsElem = []interface{}{
	"localizedData",
	"nonLocalizedData",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ToolComponentContentsElem) UnmarshalJSON(value []byte) error {
	var v string
	if err := json.Unmarshal(value, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_ToolComponentContentsElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_ToolComponentContentsElem, v)
	}
	*j = ToolComponentContentsElem(v)
	return nil
}

// IDentifies a particular toolComponent object, either the driver or an extension.
type ToolComponentReference struct {
	// The 'guid' property of the referenced toolComponent.
	Guid *string `json:"guid,omitempty" yaml:"guid,omitempty" mapstructure:"guid,omitempty"`

	// An index into the referenced toolComponent in tool.extensions.
	Index int `json:"index,omitempty" yaml:"index,omitempty" mapstructure:"index,omitempty"`

	// The 'name' property of the referenced toolComponent.
	Name *string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`

	// Key/value pairs that provide additional information about the
	// toolComponentReference.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ToolComponentReference) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ToolComponentReference
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if plain.Guid != nil {
		if matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`, string(*plain.Guid)); !matched {
			return fmt.Errorf("field %s pattern match: must match %s", "Guid", `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
		}
	}
	if v, ok := raw["index"]; !ok || v == nil {
		plain.Index = -1.0
	}
	if -1 > plain.Index {
		return fmt.Errorf("field %s: must be >= %v", "index", -1)
	}
	*j = ToolComponentReference(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ToolComponent) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in ToolComponent: required")
	}
	type Plain ToolComponent
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["contents"]; !ok || v == nil {
		plain.Contents = []ToolComponentContentsElem{
			"localizedData",
			"nonLocalizedData",
		}
	}
	if plain.DottedQuadFileVersion != nil {
		if matched, _ := regexp.MatchString(`[0-9]+(\.[0-9]+){3}`, string(*plain.DottedQuadFileVersion)); !matched {
			return fmt.Errorf("field %s pattern match: must match %s", "DottedQuadFileVersion", `[0-9]+(\.[0-9]+){3}`)
		}
	}
	if plain.Guid != nil {
		if matched, _ := regexp.MatchString(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`, string(*plain.Guid)); !matched {
			return fmt.Errorf("field %s pattern match: must match %s", "Guid", `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
		}
	}
	if v, ok := raw["isComprehensive"]; !ok || v == nil {
		plain.IsComprehensive = false
	}
	if v, ok := raw["language"]; !ok || v == nil {
		plain.Language = "en-US"
	}
	if matched, _ := regexp.MatchString(`^[a-zA-Z]{2}(-[a-zA-Z]{2})?$`, string(plain.Language)); !matched {
		return fmt.Errorf("field %s pattern match: must match %s", "Language", `^[a-zA-Z]{2}(-[a-zA-Z]{2})?$`)
	}
	if v, ok := raw["locations"]; !ok || v == nil {
		plain.Locations = []ArtifactLocation{}
	}
	if v, ok := raw["notifications"]; !ok || v == nil {
		plain.Notifications = []ReportingDescriptor{}
	}
	if v, ok := raw["rules"]; !ok || v == nil {
		plain.Rules = []ReportingDescriptor{}
	}
	if v, ok := raw["supportedTaxonomies"]; !ok || v == nil {
		plain.SupportedTaxonomies = []ToolComponentReference{}
	}
	if v, ok := raw["taxa"]; !ok || v == nil {
		plain.Taxa = []ReportingDescriptor{}
	}
	*j = ToolComponent(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Tool) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["driver"]; raw != nil && !ok {
		return fmt.Errorf("field driver in Tool: required")
	}
	type Plain Tool
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["extensions"]; !ok || v == nil {
		plain.Extensions = []ToolComponent{}
	}
	*j = Tool(plain)
	return nil
}

// Provides additional metadata related to translation.
type TranslationMetadata struct {
	// The absolute URI from which the translation metadata can be downloaded.
	DownloadURI *string `json:"downloadURI,omitempty" yaml:"downloadURI,omitempty" mapstructure:"downloadURI,omitempty"`

	// A comprehensive description of the translation metadata.
	FullDescription *MultiformatMessageString `json:"fullDescription,omitempty" yaml:"fullDescription,omitempty" mapstructure:"fullDescription,omitempty"`

	// The full name associated with the translation metadata.
	FullName *string `json:"fullName,omitempty" yaml:"fullName,omitempty" mapstructure:"fullName,omitempty"`

	// The absolute URI from which information related to the translation metadata can
	// be downloaded.
	InformationURI *string `json:"informationURI,omitempty" yaml:"informationURI,omitempty" mapstructure:"informationURI,omitempty"`

	// The name associated with the translation metadata.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// Key/value pairs that provide additional information about the translation
	// metadata.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// A brief description of the translation metadata.
	ShortDescription *MultiformatMessageString `json:"shortDescription,omitempty" yaml:"shortDescription,omitempty" mapstructure:"shortDescription,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *TranslationMetadata) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["name"]; raw != nil && !ok {
		return fmt.Errorf("field name in TranslationMetadata: required")
	}
	type Plain TranslationMetadata
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = TranslationMetadata(plain)
	return nil
}

// Specifies the information necessary to retrieve a desired revision from a
// version control system.
type VersionControlDetails struct {
	// A Coordinated Universal Time (UTC) date and time that can be used to
	// synchronize an enlistment to the state of the repository at that time.
	AsOfTimeUtc *time.Time `json:"asOfTimeUtc,omitempty" yaml:"asOfTimeUtc,omitempty" mapstructure:"asOfTimeUtc,omitempty"`

	// The name of a branch containing the revision.
	Branch *string `json:"branch,omitempty" yaml:"branch,omitempty" mapstructure:"branch,omitempty"`

	// The location in the local file system to which the root of the repository was
	// mapped at the time of the analysis.
	MappedTo *ArtifactLocation `json:"mappedTo,omitempty" yaml:"mappedTo,omitempty" mapstructure:"mappedTo,omitempty"`

	// Key/value pairs that provide additional information about the version control
	// details.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The absolute URI of the repository.
	RepositoryURI string `json:"repositoryURI" yaml:"repositoryURI" mapstructure:"repositoryURI"`

	// A string that uniquely and permanently identifies the revision within the
	// repository.
	RevisionID *string `json:"revisionID,omitempty" yaml:"revisionID,omitempty" mapstructure:"revisionID,omitempty"`

	// A tag that has been applied to the revision.
	RevisionTag *string `json:"revisionTag,omitempty" yaml:"revisionTag,omitempty" mapstructure:"revisionTag,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *VersionControlDetails) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["repositoryURI"]; raw != nil && !ok {
		return fmt.Errorf("field repositoryURI in VersionControlDetails: required")
	}
	type Plain VersionControlDetails
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = VersionControlDetails(plain)
	return nil
}

// Describes an HTTP request.
type WebRequest struct {
	// The body of the request.
	Body *ArtifactContent `json:"body,omitempty" yaml:"body,omitempty" mapstructure:"body,omitempty"`

	// The request headers.
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" mapstructure:"headers,omitempty"`

	// The index within the run.webRequests array of the request object associated
	// with this result.
	Index int `json:"index,omitempty" yaml:"index,omitempty" mapstructure:"index,omitempty"`

	// The HTTP method. Well-known values are 'GET', 'PUT', 'POST', 'DELETE', 'PATCH',
	// 'HEAD', 'OPTIONS', 'TRACE', 'CONNECT'.
	Method *string `json:"method,omitempty" yaml:"method,omitempty" mapstructure:"method,omitempty"`

	// The request parameters.
	Parameters map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty" mapstructure:"parameters,omitempty"`

	// Key/value pairs that provide additional information about the request.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The request protocol. Example: 'http'.
	Protocol *string `json:"protocol,omitempty" yaml:"protocol,omitempty" mapstructure:"protocol,omitempty"`

	// The target of the request.
	Target *string `json:"target,omitempty" yaml:"target,omitempty" mapstructure:"target,omitempty"`

	// The request version. Example: '1.1'.
	Version *string `json:"version,omitempty" yaml:"version,omitempty" mapstructure:"version,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *WebRequest) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain WebRequest
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["index"]; !ok || v == nil {
		plain.Index = -1.0
	}
	if -1 > plain.Index {
		return fmt.Errorf("field %s: must be >= %v", "index", -1)
	}
	*j = WebRequest(plain)
	return nil
}

// Describes the response to an HTTP request.
type WebResponse struct {
	// The body of the response.
	Body *ArtifactContent `json:"body,omitempty" yaml:"body,omitempty" mapstructure:"body,omitempty"`

	// The response headers.
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty" mapstructure:"headers,omitempty"`

	// The index within the run.webResponses array of the response object associated
	// with this result.
	Index int `json:"index,omitempty" yaml:"index,omitempty" mapstructure:"index,omitempty"`

	// Specifies whether a response was received from the server.
	NoResponseReceived bool `json:"noResponseReceived,omitempty" yaml:"noResponseReceived,omitempty" mapstructure:"noResponseReceived,omitempty"`

	// Key/value pairs that provide additional information about the response.
	Properties *PropertyBag `json:"properties,omitempty" yaml:"properties,omitempty" mapstructure:"properties,omitempty"`

	// The response protocol. Example: 'http'.
	Protocol *string `json:"protocol,omitempty" yaml:"protocol,omitempty" mapstructure:"protocol,omitempty"`

	// The response reason. Example: 'Not found'.
	ReasonPhrase *string `json:"reasonPhrase,omitempty" yaml:"reasonPhrase,omitempty" mapstructure:"reasonPhrase,omitempty"`

	// The response status code. Example: 451.
	StatusCode *int `json:"statusCode,omitempty" yaml:"statusCode,omitempty" mapstructure:"statusCode,omitempty"`

	// The response version. Example: '1.1'.
	Version *string `json:"version,omitempty" yaml:"version,omitempty" mapstructure:"version,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *WebResponse) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain WebResponse
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["index"]; !ok || v == nil {
		plain.Index = -1.0
	}
	if -1 > plain.Index {
		return fmt.Errorf("field %s: must be >= %v", "index", -1)
	}
	if v, ok := raw["noResponseReceived"]; !ok || v == nil {
		plain.NoResponseReceived = false
	}
	*j = WebResponse(plain)
	return nil
}
