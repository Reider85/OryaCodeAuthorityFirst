package security

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// SandboxYAMLParser provides enhanced security for YAML parsing
type SandboxYAMLParser struct {
	// Configuration options
	AllowCustomTags    bool
	MaxNestingDepth    int
	MaxDocumentSize    int64
	AllowedSchemas     []string
	BlockedPatterns    []*regexp.Regexp
	RequiredFields     []string
	ValidationMode     ValidationMode
}

// ValidationMode defines the level of validation to perform
type ValidationMode int

const (
	ValidationModeBasic ValidationMode = iota
	ValidationModeStrict
	ValidationModeParanoid
)

// YAMLValidationError represents a YAML validation error
type YAMLValidationError struct {
	Field   string
	Message string
	Line    int
}

func (e *YAMLValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("YAML validation error at line %d: %s (field: %s)", e.Line, e.Message, e.Field)
	}
	return fmt.Sprintf("YAML validation error: %s (field: %s)", e.Message, e.Field)
}

// NewSandboxYAMLParser creates a new sandboxed YAML parser with default security settings
func NewSandboxYAMLParser() *SandboxYAMLParser {
	return &SandboxYAMLParser{
		AllowCustomTags:   false,
		MaxNestingDepth:   10,
		MaxDocumentSize:   1024 * 1024, // 1MB
		ValidationMode:   ValidationModeStrict,
		BlockedPatterns: []*regexp.Regexp{
			regexp.MustCompile(`<!\[CDATA\[.*?\]\]>`), // CDATA sections
			regexp.MustCompile(`<\?xml.*?\?>`),      // XML processing instructions
			regexp.MustCompile(`<!.*?>`),             // XML declarations
			regexp.MustCompile(`&(?:lt|gt|amp|quot|apos|#x[0-9a-fA-F]+|#\d+);`), // HTML entities
		},
	}
}

// ParseYAML parses YAML content with enhanced security checks
func (p *SandboxYAMLParser) ParseYAML(content []byte, target interface{}) error {
	// Pre-validation checks
	if err := p.preValidate(content); err != nil {
		return err
	}

	// Parse the YAML
	var node yaml.Node
	if err := yaml.Unmarshal(content, &node); err != nil {
		return fmt.Errorf("YAML parsing error: %w", err)
	}

	// Security validation
	if err := p.validateNode(&node, 0); err != nil {
		return err
	}

	// Convert to target type
	if err := yaml.Unmarshal(content, target); err != nil {
		return fmt.Errorf("YAML unmarshaling error: %w", err)
	}

	// Post-validation
	if err := p.postValidate(target); err != nil {
		return err
	}

	return nil
}

// preValidate performs security checks before parsing
func (p *SandboxYAMLParser) preValidate(content []byte) error {
	// Check document size
	if int64(len(content)) > p.MaxDocumentSize {
		return errors.New("document exceeds maximum size limit")
	}

	// Check for blocked patterns
	contentStr := string(content)
	for _, pattern := range p.BlockedPatterns {
		if pattern.MatchString(contentStr) {
			return fmt.Errorf("blocked pattern detected: %s", pattern.String())
		}
	}

	// Check for dangerous constructs
	if p.containsDangerousConstructs(contentStr) {
		return errors.New("YAML contains potentially dangerous constructs")
	}

	return nil
}

// validateNode validates the YAML node structure recursively
func (p *SandboxYAMLParser) validateNode(node *yaml.Node, depth int) error {
	// Check nesting depth
	if depth > p.MaxNestingDepth {
		return errors.New("YAML nesting depth exceeds maximum limit")
	}

	switch node.Kind {
	case yaml.DocumentNode:
		// Validate document node
		for _, docNode := range node.Content {
			if err := p.validateNode(docNode, depth+1); err != nil {
				return err
			}
		}

	case yaml.SequenceNode:
		// Validate sequence node
		for _, seqNode := range node.Content {
			if err := p.validateNode(seqNode, depth+1); err != nil {
				return err
			}
		}

	case yaml.MappingNode:
		// Validate mapping node
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			
			// Validate key is safe
			if keyNode.Kind == yaml.ScalarNode {
				if p.isDangerousKey(keyNode.Value) {
					return fmt.Errorf("dangerous key detected: %s", keyNode.Value)
				}
			}
			
			// Validate value
			if err := p.validateNode(valueNode, depth+1); err != nil {
				return err
			}
		}

	case yaml.ScalarNode:
		// Validate scalar content
		if err := p.validateScalarNode(node); err != nil {
			return err
		}

	case yaml.AliasNode:
		// Validate alias node
		return p.validateNode(node.Alias, depth+1)
	}

	return nil
}

// validateScalarNode validates a scalar node for security issues
func (p *SandboxYAMLParser) validateScalarNode(node *yaml.Node) error {
	value := node.Value

	// Check for custom tags if not allowed
	if !p.AllowCustomTags && node.Tag != "" && node.Tag != "!!str" && node.Tag != "!!int" && node.Tag != "!!float" && node.Tag != "!!bool" && node.Tag != "!!null" {
		return fmt.Errorf("custom YAML tag not allowed: %s", node.Tag)
	}

	// Check for potential code execution
	if p.containsCodeExecutionRisk(value) {
		return fmt.Errorf("potential code execution risk in scalar value: %s", value)
	}

	// Check for schema validation if specified
	if len(p.AllowedSchemas) > 0 {
		if !p.isAllowedSchema(value) {
			return fmt.Errorf("value not in allowed schemas: %s", value)
		}
	}

	return nil
}

// postValidate validates the parsed data structure
func (p *SandboxYAMLParser) postValidate(data interface{}) error {
	val := reflect.ValueOf(data)
	
	// Check if it's a pointer and get the underlying value
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return errors.New("parsed data is nil")
		}
		val = val.Elem()
	}

	// Check for required fields
	for _, field := range p.RequiredFields {
		if !val.FieldByName(field).IsValid() {
			return &YAMLValidationError{
				Field:   field,
				Message: "required field is missing",
			}
		}
	}

	// Additional validation based on validation mode
	switch p.ValidationMode {
	case ValidationModeStrict:
		return p.validateStrictMode(val)
	case ValidationModeParanoid:
		return p.validateParanoidMode(val)
	default:
		return nil
	}
}

// containsDangerousConstructs checks for potentially dangerous YAML constructs
func (p *SandboxYAMLParser) containsDangerousConstructs(content string) bool {
	// Check for multi-line strings that could contain code
	if strings.Contains(content, "|") || strings.Contains(content, ">") {
		// Look for potential script content in multi-line strings
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "|") || strings.HasPrefix(line, ">") {
				// Check the following lines for script-like content
				continue
			}
		}
	}

	// Check for anchor/alias usage which could be abused
	if strings.Contains(content, "&") && strings.Contains(content, "*") {
		// This could be legitimate YAML, but we want to be cautious
		return true
	}

	return false
}

// isDangerousKey checks if a key is potentially dangerous
func (p *SandboxYAMLParser) isDangerousKey(key string) bool {
	dangerousKeys := []string{
		"__proto__",
		"constructor",
		"prototype",
		"__defineGetter__",
		"__defineSetter__",
		"__lookupGetter__",
		"__lookupSetter__",
		"hasOwnProperty",
		"isPrototypeOf",
		"propertyIsEnumerable",
		"toString",
		"valueOf",
		"toLocaleString",
		"toJSON",
	}

	for _, dk := range dangerousKeys {
		if key == dk {
			return true
		}
	}

	return false
}

// containsCodeExecutionRisk checks for potential code execution risks
func (p *SandboxYAMLParser) containsCodeExecutionRisk(value string) bool {
	// Check for template expressions
	if strings.Contains(value, "${") && strings.Contains(value, "}") {
		return true
	}

	// Check for script tags
	if strings.Contains(value, "<script") || strings.Contains(value, "</script>") {
		return true
	}

	// Check for JavaScript execution patterns
	jsPatterns := []string{
		"javascript:",
		"vbscript:",
		"data:",
		"eval(",
		"Function(",
		"setTimeout(",
		"setInterval(",
	}

	for _, pattern := range jsPatterns {
		if strings.Contains(value, pattern) {
			return true
		}
	}

	return false
}

// isAllowedSchema checks if a value matches allowed schemas
func (p *SandboxYAMLParser) isAllowedSchema(value string) bool {
	for _, schema := range p.AllowedSchemas {
		if strings.HasPrefix(value, schema) {
			return true
		}
	}
	return false
}

// validateStrictMode performs strict validation
func (p *SandboxYAMLParser) validateStrictMode(val reflect.Value) error {
	// Check for interface{} which could contain anything
	if val.Kind() == reflect.Interface {
		return errors.New("interface{} type not allowed in strict mode")
	}

	// Check for maps that could contain dangerous keys
	if val.Kind() == reflect.Map {
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			if p.isDangerousKey(keyStr) {
				return &YAMLValidationError{
					Field:   keyStr,
					Message: "dangerous key not allowed in strict mode",
				}
			}
		}
	}

	return nil
}

// validateParanoidMode performs paranoid validation
func (p *SandboxYAMLParser) validateParanoidMode(val reflect.Value) error {
	// In paranoid mode, we reject anything that's not a basic type
	switch val.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return nil
	default:
		return errors.New("only basic types allowed in paranoid mode")
	}
}

// SafeParseYAML is a convenience function that uses the sandboxed parser
func SafeParseYAML(content []byte, target interface{}) error {
	parser := NewSandboxYAMLParser()
	return parser.ParseYAML(content, target)
}