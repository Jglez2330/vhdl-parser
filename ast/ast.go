package ast

import (
	"vhdl/token"
)

//VHDL AST

// AUX types
type CommentGroup struct {
	Comments []*string
}

type Comment struct {
	Slash token.Pos // position of "/" starting the comment
	Text  string    // comment text (excluding '\n' for //-style comments)
}

//We are following the VHDL 2019 LRM
//https://standards.ieee.org/standard/1076-2019.html

type Node struct {
	Pos token.Pos
	End token.Pos
}

type OptionalNode interface {
	Type() string
}

// 3 Design entities and configurations
type EntityDeclaration struct {
	Identifier            Identifier
	EntityHeader          EntityHeader
	EntityDeclarativePart EntityDeclarativePart
	EntityStatementPart   *EntityStatementPart
	EntitySimpleName      *SimpleName
	Node
	LibraryUnit
	PrimaryUnit
}

type EntitySimpleName struct {
	SimpleName
}

type EntityHeader struct {
	FormalGenericClause *FormalGenericClause
	FormalPortClause    *FormalPortClause
	Node
}

type FormalGenericClause struct {
}

type FormalPortClause struct {
}

type EntityDeclarativePart struct {
	EntityDeclarativeItems *[]EntityDeclarativeItem
	Node
}

type EntityDeclarativeItem interface {
}

type EntityStatementPart struct {
	EntityStatements *[]EntityStatement
	Node
}

type EntityStatement interface{}

type ArchitectureBody struct {
	Identifier                Identifier
	EntityName                Name
    ArchitectureDeclarativePart ArchitectureDeclarativePart
	ArchitectureStatementPart ArchitectureStatementPart
	ArchitectureSimpleName    *SimpleName
    SecondaryUnit
    Node
}

type EntityName struct {
	SimpleName
}

type ArchitectureDeclarativePart struct {
	BlockDeclarativeItems *[]BlockDeclarativeItem
}

type BlockDeclarativeItem interface{}

type ArchitectureStatementPart struct {
	ConcurrentStatements *[]ConcurrentStatement
}

type ConcurrentStatement interface{}

type ArchitectureSimpleName struct {
	SimpleName
}

/*
type ConfigurationDeclaration struct {
	Identifier                         Identifier
	EntityName                         EntityName
	ConfigurationDeclarativePart       ConfigurationDeclarativePart
	VerificationUnitBindingIndications *[]VerificationUnitBindingIndication
	BlockConfiguration                 BlockConfiguration
	ConfigurationSimpleName            *ConfigurationSimpleName
}

type ConfigurationDeclarativePart struct {
	ConfigurationDeclarativeItems *[]ConfigurationDeclarativeItem
}

type ConfigurationDeclarativeItem interface{}

type BlockConfiguration struct {
	BlockSpecification BlockSpecification
	UseClauses         *[]UseClause
	ConfigurationItems *[]ConfigurationItem
}

type BlockSpecification interface{}

type GennerateSpecification interface{}

type ConfigurationItem interface{}

type ComponentConfiguration struct {
	ComponentSpecification             ComponentSpecification
	BindingIndication                  *BindingIndication
	VerificationUnitBindingIndications *[]VerificationUnitBindingIndication
	BlockConfiguration                 *BlockConfiguration
}*/

/*
// 4 Subprogram and package
type SubprogramDeclaration struct {
	SubprogramSpecification SubprogramSpecification
}

type SubprogramSpecification interface{}

type ProcedureSpecification struct {
	Designator          Designator
	SubprogramHeader    SubprogramHeader
	FormalParameterList *FormalParameterList
}

type FunctionSpecification struct {
	Designator          Designator
	SubprogramHeader    SubprogramHeader
	FormalParameterList *FormalParameterList
	ReturnIdentifier    ReturnIdentifier
	TypeMark            TypeMark
}

type SubprogramHeader struct {
	GenericList      *GenericList
	GenericMapAspect *GenericMapAspect
}

type Designator interface {
	OptionalNode
}

type OperatorSymbol struct {
	StringLiteral
}

type FormalParameterList struct {
	ParameterInterfaceList
}

type SubprogramBody struct {
	SubprogramSpecification   SubprogramSpecification
	SubprogramDeclarativePart SubprogramDeclarativePart
	SubprogramStatementPart   SubprogramStatementPart
	SubprogramKind            *SubprogramKind
	Designator                *Designator
}

type SubprogramDeclarativePart struct {
	SubprogramDeclarativeItems *[]SubprogramDeclarativeItem
}

type SubprogramDeclarativeItem interface {
	OptionalNode
}

type SubprogramStatementPart struct {
	SequentialStatements *[]SequentialStatement
}

type SubprogramKind struct {
	Node
	token token.Token
}

type Signature struct {
	TypeMark       *TypeMark
	TypeMarkList   *[]TypeMark
	ReturnTypeMark *TypeMark
}

type PackageBody struct {
	PackageSimpleName          PackageSimpleName
	PackageBodyDeclarativePart PackageBodyDeclarativePart
	ClosingPackageSimpleName   *PackageSimpleName
}

type PackageBodyDeclarativePart struct {
	PackageBodyDeclarativeItems *[]PackageBodyDeclarativeItem
}

type PackageBodyDeclarativeItem interface{}

type PackageInstantiationDeclaration struct {
	Identifier                Identifier
	UninstantiatedPackageName UninstantiatedPackageName
	GenericMapAspect          *GenericMapAspect
}

// 5 Types

type ScalarTypeDefinition interface{}

type RangeConstraint struct {
	Range Range
}

type Range interface{}

type SimpleRange struct {
	SimpleExpression  SimpleExpression
	Direction         Direction
	SimpleExpression2 SimpleExpression
}

type Direction struct {
	Node
	token token.Token
}

type EnumerationTypeDefinition struct {
	EnumerationLiteral     EnumerationLiteral
	EnumerationLiteralList *[]EnumerationLiteral
}

type EnumerationLiteral interface{}

type IntegerTypeDefinition struct {
	RangeConstraint
}

type PhysicalTypeDefinition struct {
	RangeConstraint          RangeConstraint
	PrimaryUnitDeclaration   PrimaryUnitDeclaration
	SecondaryUnitDeclaration *[]SecondaryUnitDeclaration
	PhysicalTypeSimpleName   *PhysicalTypeSimpleName
}

type PrimaryUnitDeclaration struct {
	Identifier
}

type SecondaryUnitDeclaration struct {
	Identifier      Identifier
	PhysicalLiteral PhysicalLiteral
}

type PhysicalLiteral struct {
	AbstractLiteral *AbstractLiteral
	UnitName        UnitName
}

type FloatingTypeDefinition struct {
	RangeConstraint RangeConstraint
}

type CompositeTypeDefinition interface{}

type ArrayTypeDefinition interface{}

type UnboundedArrayDefinition struct {
	IndexSubtypeDefinition     IndexSubtypeDefinition
	IndexSubtypeDefinitionList *[]IndexSubtypeDefinition
	ElementSubtypeIndication   ElementSubtypeIndication
}

type ConstrainedArrayDefinition struct {
	IndexConstraint          IndexConstraint
	ElementSubtypeIndication ElementSubtypeIndication
}

type IndexSubtypeDefinition struct {
	TypeMark TypeMark
}

type ArrayConstraint interface {
	OptionalNode
}

type IndexConstraint struct {
	DiscreteRange     DiscreteRange
	DiscreteRangeList *[]DiscreteRange
}

type DiscreteRange interface{}

type RecordTypeDefinition struct {
	ElementDeclarations  *[]ElementDeclaration
	RecordTypeSimpleName *RecordTypeSimpleName
}

type ElementDeclaration struct {
	IdentifierList           IdentifierList
	ElementSubtypeDefinition ElementSubtypeDefinition
}

type IdentifierList struct {
	Identifier
	IdentifierList *[]Identifier
}

type ElementSubtypeDefinition struct {
	SubtypeIndication
}

type RecordConstraint struct {
	RecordElementConstraint     RecordElementConstraint
	RecordElementConstraintList *[]RecordElementConstraint
}

type RecordElementConstraint struct {
	RecordElementSimpleName RecordElementSimpleName
	ElementConstraint       ElementConstraint
}

type AccessTypeDefinition struct {
	SubtypeIndication SubtypeIndication
	GenericMapAspect  *GenericMapAspect
}

type ImcompleteTypeDeclaration struct {
	Identifier Identifier
}

type FileTypeDefinition struct {
	TypeMark TypeMark
}

type ProtectedTypeDefinition interface{}

type ProtectedTypeDeclaration struct {
	ProtectedTypeHeader          ProtectedTypeHeader
	ProtectedTypeDeclarativePart ProtectedTypeDeclarativePart
	ProtectedTypeSimpleName      *ProtectedTypeSimpleName
}

type ProtectedTypeHeader struct {
	GenericClause    *GenericClause
	GenericMapAspect *GenericMapAspect
}

type ProtectedTypeDeclarativePart struct {
	ProtectedTypeDeclarativeItem  ProtectedTypeDeclarativeItem
	ProtectedTypeDeclarativeItems *[]ProtectedTypeDeclarativeItem
}

type ProtectedTypeDeclarativeItem interface{}

type PrivateVariableDeclaration struct {
	VariableDeclaration VariableDeclaration
}

type ProtectedTypeBody struct {
	ProtectedTypeBodyDeclarativePart ProtectedTypeBodyDeclarativePart
	ProtectedTypeSimpleName          *ProtectedTypeSimpleName
}

type ProtectedTypeBodyDeclarativePart struct {
	ProtectedTypeBodyDeclarativeItems *[]ProtectedTypeBodyDeclarativeItem
}

type ProtectedTypeBodyDeclarativeItem interface{}

type ProtectedTypeInstantiationDefinition struct {
	SubtypeIndication SubtypeIndication
	GenericMapAspect  *GenericMapAspect
}

type UnspecifiedTypeIndication struct {
	IncompleteTypeDefintion IncompleteTypeDefintion
}

type IncompleteTypeDefintion interface{}

type IncompleteSubtypeIndication interface{}

type PrivateIncompleteTypeDefinition struct {
	token token.Token
}

type ScalarIncompleteTypeDefinition struct {
	token token.Token
}

type DiscreteIncompleteTypeDefinition struct {
	Lbracket   token.Token
	RangeToken token.Token
	Rbracket   token.Token
}

type IntegerIncompleteTypeDefinition struct {
	RangeKey   token.Token
	RangeToken token.Token
}

type PhysicalIncompleteTypeDefinition struct {
	UnitsKeyword token.Token
	RangeToken   token.Token
}

type FloatingIncompleteTypeDefinition struct {
	RangeKey     token.Token
	RangeToken_1 token.Token
	Dot          token.Token
	RangeToken_2 token.Token
}

type ArrayIncompleteTypeDefinition struct {
	ArrayIndexIncompleteTypeList       ArrayIndexIncompleteTypeList
	ElementIncompleteSubtypeIndication ElementIncompleteSubtypeIndication
}

type ArrayIndexIncompleteTypeList struct {
	ArrayIndexIncompleteType     ArrayIndexIncompleteType
	ArrayIndexIncompleteTypeList *[]ArrayIndexIncompleteType
}

type ArrayIndexIncompleteType interface{}

type AccessIncompleteTypeDefinition struct {
	AccessIncompleteSubtypeIndication AccessIncompleteSubtypeIndication
}

type FileIncompleteTypeDefinition struct {
	FileIncompleteTypeMark FileIncompleteTypeMark
}

//6 Declarations and types

type TypeDeclaration interface{}

type FullTypeDeclaration struct {
	Identifier     Identifier
	TypeDefinition TypeDefinition
}

type TypeDefinition interface{}

type SubtypeDeclaration struct {
	Identifier        Identifier
	SubtypeIndication SubtypeIndication
}

type SubtypeIndication struct {
	ResolutionIndication *ResolutionIndication
	TypeMark             TypeMark
	Constraint           *Constraint
}

type ResolutionIndication interface{}

type ElementResolution interface{}

type ArrayElementResolution struct {
	ResolutionIndication
}

type RecordResolution struct {
	RecordElementResolution     RecordElementResolution
	RecordElementResolutionList *[]RecordElementResolution
}

type RecordElementResolution struct {
	RecordElementSimpleName RecordElementSimpleName
	ReslutionIndication     ResolutionIndication
}

type TypeMark interface{}

type Constraint interface{}

type ElementConstraint interface{}

type ObjectDeclaration interface{}

type VariableDeclaration struct {
	IdentifierList        IdentifierList
	SubtypeIndication     SubtypeIndication
	GenericMapAspect      *GenericMapAspect
	ConditionalExpression *ConditionalExpression
}

type FileDeclaration struct {
	IdentifierList      IdentifierList
	SubtypeIndication   SubtypeIndication
	FileOpenInformation *FileOpenInformation
}

type FileOpenInformation struct {
	FileOpenKindExpression *FileOpenKindExpression
	FileLogicalName        FileLogicalName
}

type FileLogicalName struct {
	StringLiteral
}

type InterfaceDeclaration interface{}

type InterfaceObjectDeclaration interface{}

type InterfaceConstantDeclaration struct {
	IdentifierList              IdentifierList
	InterfaceTypeIndication     InterfaceTypeIndication
	StaticConditionalExpression *StaticConditionalExpression
}

type InterfaceSignalDeclaration struct {
	IdentifierList              IdentifierList
	Mode                        *Mode
	InterfaceTypeIndication     InterfaceTypeIndication
	StaticConditionalExpression *StaticConditionalExpression
}

type InterfaceFileDeclaration struct {
	IdentifierList    IdentifierList
	SubtypeIndication SubtypeIndication
}

type InterfaceTypeIndication interface{}

type ModeIndication interface{}

type SimpleModeIndication struct {
	Mode                        *Mode
	InterfaceTypeIndication     InterfaceTypeIndication
	StaticConditionalExpression *StaticConditionalExpression
}

type Mode struct {
	Token token.Token
}

type ModeViewIndication interface{}

type RecordModeViewIndication struct {
	ModeViewName                      ModeViewName
	UnresolvedRecordSubtypeIndication *UnresolvedRecordSubtypeIndication
}

type ArrayModeViewIndication struct {
	ModeViewName                     ModeViewName
	UnresolvedArraySubtypeIndication *UnresolvedArraySubtypeIndication
}

type ModeViewDeclaration struct {
	Identifier                        Identifier
	UnresolvedRecordSubtypeIndication UnresolvedRecordSubtypeIndication
	ModeViewElementDefinition         *[]ModeViewElementDefinition
	ModeViewSimpleName                *ModeViewSimpleName
}

type ModeViewElementDefinition struct {
	RecordElementList     RecordElementList
	ElementModeIndication ElementModeIndication
}

type RecordElementList struct {
	RecordElementSimpleName     RecordElementSimpleName
	RecordElementSimpleNameList *[]RecordElementSimpleName
}

type ElementModeIndication interface{}

type ElementModeViewIndication interface{}

type ElementRecordModeViewIndication struct {
	ModeViewName ModeViewName
}

type ElementArrayModeViewIndication struct {
    ModeViewName ModeViewName
}

type InterfaceTypeDeclaration struct {
    Type Type
    Identifier Identifier
    IncompleteTypeDefintion *IncompleteTypeDefintion
}


type InterfaceSubprogramDeclaration struct {
    InterfaceSubprogramSpecification InterfaceSubprogramSpecification
    InterfaceSubprogramDefault *InterfaceSubprogramDefault
}

type InterfaceSubprogramSpecification interface{}

type IntefaceProcedureSpecification struct {
    Designator Designator
    FormalParameterList *FormalParameterList
}

type InterfaceFunctionSpecification struct {
    Designator Designator
    FormalParameterList *FormalParameterList
    TypeMark TypeMark
}

type InterfaceSubprogramDefault interface{}

type InterfacePackageDeclaration struct {
    Identifier Identifier
    UnstantiatedPackageName UninstantiatedPackageName
    InterfaceGenericMapAspect InterfaceGenericMapAspect
}

type InterfaceGenericMapAspect interface{}

type InterfaceList struct {
    InterfaceElement InterfaceElement
    InterfaceElementList *[]InterfaceElement
}

type InterfaceElement struct {
    InterfaceDeclaration
}

type GenericClause struct {
    GenericList GenericList
}

type GenericList struct {
    GenericInterfaceList
}


type PortClause struct {
    PortList PortList
}

type PortList struct {
    PortInterfaceList
}

type AssociationList struct {
    AssociationElement AssociationElement
    AssociationElementList *[]AssociationElement
}

type AssociationElement struct {
    FormalPart *FormalPart
    ActualPart ActualPart
}

type FormalPart interface{}

type FormalDesignator interface{}

type ActualPart interface{}

type ActualDesignator interface{}

type GenericMapAspect struct {
    GenericAssociationList GenericAssociationList
}

type PortMapAspect struct {
    PortAssociationList PortAssociationList
}

type AliasDeclaration struct {
    AliasDesignator AliasDesignator
    SubtypeIndication *SubtypeIndication
    Name Name
    Signature *Signature
}

type AliasDesignator interface{}


type AttributeDeclaration struct {
    Identifier Identifier
    TypeMark TypeMark
}

type ComponentDeclaration struct{
    Identifier Identifier
    LocalGenericClause *LocalGenericClause
    LocalPortClause *LocalPortClause
    ComponentSimpleName *ComponentSimpleName
}

type GroupTemplateDeclaration struct{
    Identifier Identifier
    EntityClassEntryList EntityClassEntryList
}

type EntityClassEntryList struct{
    EntityClassEntry EntityClassEntry
    EntityClassEntries *[] EntityClassEntry
}

type EntityClassEntry struct{
    EntityClass EntityClass
}

type GroupDeclaration struct{
    Identifier  Identifier
    GroupTemplateName GroupTemplateName
    GroupConstituentList GroupConstituentList
}
type GroupConstituentList struct{
    GroupConstituent GroupConstituent
    GroupConstituents *[]GroupConstituent
}

type GroupConstituent interface{}
*/

type UseClause struct {
	SelectedName     SelectedName
	SelectedNameList *[]SelectedName
}

type SelectedName struct {
	Prefix Prefix
	Suffix Suffix
	Name
}

type Prefix interface{ OptionalNode }

type Name interface{}

type Suffix interface{ OptionalNode }

type SimpleName struct {
	Identifier Identifier
	Suffix
	Name
}

type CharacterLiteral struct {
	GraphicCharacter GraphicCharacter
	Suffix
	Name
}

type GraphicCharacter struct {
	Token token.Token
}

type OperatorSymbol struct {
	Token token.Token
	//StringLiteral
	Suffix
	Name
}

type SuffixKeyword struct {
	Token token.Token
	Suffix
}

type Identifier struct {
	Identifier string
	Node
}

type File struct {
	FileStart, FileEnd token.Pos // start and end of entire file
	DesignUnits        []DesignUnit
}

type DesignUnit struct {
	ContextClause ContextClause
	LibraryUnit   LibraryUnit
}

type ContextClause interface{}

type LibraryUnit interface{}

type PrimaryUnit interface {
	LibraryUnit
}

type SecondaryUnit interface {
	LibraryUnit
}
