package ast

import (
	"vhdl/token"
)

//VHDL AST

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
	EntitySimpleName      *EntitySimpleName
	Node
}

type EntityHeader struct {
	FormalGenericClause *FormalGenericClause
	FormalPortClause    *FormalPortClause
	Node
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
	EntityName                EntityName
	ArchitectureStatementPart ArchitectureStatementPart
	ArchitectureSimpleName    ArchitectureSimpleName
}

type ArchitectureDeclarativePart struct {
	BlockDeclarativeItems *[]BlockDeclarativeItem
}

type BlockDeclarativeItem interface{}

type ArchitectureStatementPart struct {
	ConcurrentStatements *[]ConcurrentStatement
}

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
}


// 4 Subprogram and package
type SubprogramDeclaration struct{
    SubprogramSpecification SubprogramSpecification
}

type SubprogramSpecification interface{}

type ProcedureSpecification struct{
    Designator Designator
    SubprogramHeader SubprogramHeader
    FormalParameterList *FormalParameterList
}

type FunctionSpecification struct{
    Designator Designator
    SubprogramHeader SubprogramHeader
    FormalParameterList *FormalParameterList
    ReturnIdentifier ReturnIdentifier
    TypeMark TypeMark
}

type SubprogramHeader struct{
    GenericList *GenericList
    GenericMapAspect *GenericMapAspect
}

type Designator interface{
    OptionalNode
}

type OperatorSymbol struct{
    StringLiteral
}

type FormalParameterList struct{
    ParameterInterfaceList
}

type SubprogramBody struct{
    SubprogramSpecification SubprogramSpecification
    SubprogramDeclarativePart SubprogramDeclarativePart
    SubprogramStatementPart SubprogramStatementPart
    SubprogramKind *SubprogramKind
    Designator *Designator
}

type SubprogramDeclarativePart struct{
    SubprogramDeclarativeItems *[]SubprogramDeclarativeItem
}

type SubprogramDeclarativeItem interface{
    OptionalNode
}

type SubprogramStatementPart struct{
    SequentialStatements *[]SequentialStatement
}

type SubprogramKind struct{
    Node
    token token.Token
}

type Signature struct{
    TypeMark *TypeMark
    TypeMarkList *[]TypeMark
    ReturnTypeMark *TypeMark
}

type PackageBody struct{
    PackageSimpleName PackageSimpleName
    PackageBodyDeclarativePart PackageBodyDeclarativePart
    ClosingPackageSimpleName *PackageSimpleName
}

type PackageBodyDeclarativePart struct{
    PackageBodyDeclarativeItems *[]PackageBodyDeclarativeItem
}

type PackageBodyDeclarativeItem interface{}

type PackageInstantiationDeclaration struct{
    Identifier Identifier
    UninstantiatedPackageName UninstantiatedPackageName
    GenericMapAspect *GenericMapAspect
}

// 5 Types

type ScalarTypeDefinition interface{}

type RangeConstraint struct{
    Range Range
}

type Range interface{}

type SimpleRange struct{
    SimpleExpression SimpleExpression
    Direction Direction
    SimpleExpression2 SimpleExpression
}

type Direction struct{
    Node
    token token.Token
}

type EnumerationTypeDefinition struct{
    EnumerationLiteral EnumerationLiteral
    EnumerationLiteralList *[]EnumerationLiteral
}

type EnumerationLiteral interface{}

type IntegerTypeDefinition struct{
    RangeConstraint
}

type PhysicalTypeDefinition struct{
    RangeConstraint RangeConstraint
    PrimaryUnitDeclaration PrimaryUnitDeclaration
    SecondaryUnitDeclaration *[]SecondaryUnitDeclaration
    PhysicalTypeSimpleName *PhysicalTypeSimpleName
}

type PrimaryUnitDeclaration struct{
    Identifier
}

type SecondaryUnitDeclaration struct{
    Identifier Identifier
    PhysicalLiteral PhysicalLiteral
}

type PhysicalLiteral struct{
    AbstractLiteral *AbstractLiteral
    UnitName UnitName
}

type FloatingTypeDefinition struct{
    RangeConstraint RangeConstraint
}

type CompositeTypeDefinition interface{}

type ArrayTypeDefinition interface{}

type UnboundedArrayDefinition struct{
    IndexSubtypeDefinition IndexSubtypeDefinition
    IndexSubtypeDefinitionList *[]IndexSubtypeDefinition
    ElementSubtypeIndication ElementSubtypeIndication
}

type ConstrainedArrayDefinition struct{
    IndexConstraint IndexConstraint
    ElementSubtypeIndication ElementSubtypeIndication
}

type IndexSubtypeDefinition struct{
    TypeMark TypeMark
}

type ArrayConstraint interface{
    OptionalNode
}

type IndexConstraint struct{
    DiscreteRange DiscreteRange
    DiscreteRangeList *[]DiscreteRange  
}

type DiscreteRange interface{}

type RecordTypeDefinition struct{
    ElementDeclarations *[]ElementDeclaration
    RecordTypeSimpleName *RecordTypeSimpleName
}

type ElementDeclaration struct{
    IdentifierList IdentifierList
    ElementSubtypeDefinition ElementSubtypeDefinition
}

type IdentifierList struct{
    Identifier
    IdentifierList *[]Identifier
}

type ElementSubtypeDefinition struct{
    SubtypeIndication
}

type RecordConstraint struct{
    RecordElementConstraint RecordElementConstraint
    RecordElementConstraintList *[]RecordElementConstraint
}

type RecordElementConstraint struct{
    RecordElementSimpleName RecordElementSimpleName
    ElementConstraint ElementConstraint
}

type FileTypeDefinition struct{
    TypeMark TypeMark
}






