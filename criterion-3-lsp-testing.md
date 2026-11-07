# Criterion 3 Verification: LSP Fallback Symbol Graph Testing
**Testing the LSP fallback symbol graph generation for Go repositories with 1000+ files**

## Overview
This document details the verification process for Criterion 3: "LSP fallback builds symbol graph for Go repo up to 1000 files" - a real-world testing requirement for the MVP → ALPHA transition.

## Test Objective
Verify that the LSP fallback mechanism can successfully build a symbol graph for a Go repository containing 1000+ files within acceptable performance limits.

## Test Environment Setup

### Required Components
1. **Go Development Environment**: Go 1.22+ installed
2. **gopls**: Go language server for symbol graph generation
3. **Large Go Repository**: Test repository with 1000+ Go files
4. **caf Binary**: Compiled caf binary with LSP support

### Test Repository Preparation
```bash
# Create a large Go test repository
mkdir large-go-test-repo
cd large-go-test-repo
git init
go mod init test-repo

# Generate 1000+ Go files with various structures
# (This would be done using scripts to generate realistic test code)
```

## LSP Implementation Analysis

### Current LSP Implementation
```go
// LSP fallback implementation in internal/caf/lsp/
- Symbol graph generation: cache/ directory structure
- Fallback mechanism: gopls-based symbol graph
- Zero-config coverage: Tier 4 adapter implementation
```

### Capability Matrix Validation
```go
var contractCases = []struct {
  id   string
  want caf.CapabilityMatrix
}{
  {
    id: "universal-lsp",
    want: caf.CapabilityMatrix{
      PreSession: caf.TierNone,      // No pre-session for LSP
      PostEdit:   caf.TierWatcher,   // File watcher
      PreCommit:  caf.TierUniversal, // Universal hook
      MCP:        false,             // No MCP support
      Tier:       4,                 // Fallback tier
    },
  },
}
```

### LSP Adapter Features
- ✅ **Symbol Graph Generation**: Uses gopls for symbol information
- ✅ **File Watching**: Monitors file changes for updates
- ✅ **Universal Hook Integration**: Integrates with pre-commit hooks
- ✅ **Zero Configuration**: Works without additional setup

## Test Scenarios

### Scenario 1: Basic Symbol Graph Generation
**Test Objective**: Verify LSP can generate basic symbol graph for 1000+ files

#### Test Steps
1. Setup test repository with 1000+ Go files
2. Initialize caf in the repository
3. Trigger LSP symbol graph generation
4. Measure generation time and resource usage
5. Verify symbol graph completeness

#### Expected Results
- ✅ Symbol graph generation completes successfully
- ✅ All 1000+ files processed
- ✅ Symbol information extracted accurately
- ✅ Generation time < 5 minutes
- ✅ Memory usage < 500MB

#### Success Criteria
- Symbol graph contains all expected symbols
- No errors during generation
- Performance within acceptable limits

### Scenario 2: Performance Under Load
**Test Objective**: Verify LSP performance with 1000+ files under various conditions

#### Test Steps
1. Setup test repository with 1000+ Go files
2. Generate symbol graph multiple times
3. Measure performance consistency
4. Test with concurrent file modifications
5. Monitor memory usage patterns

#### Expected Results
- ✅ Consistent performance across multiple runs
- ✅ Memory usage remains stable
- ✅ No performance degradation over time
- ✅ Handles concurrent modifications gracefully

#### Success Criteria
- Generation time variation < 20%
- Memory usage remains < 500MB
- No crashes or memory leaks

### Scenario 3: Symbol Accuracy Validation
**Test Objective**: Verify symbol accuracy and completeness

#### Test Steps
1. Setup test repository with known symbol structures
2. Generate symbol graph
3. Compare extracted symbols with expected symbols
4. Validate symbol relationships and types
5. Check for missing or incorrect symbols

#### Expected Results
- ✅ Symbol accuracy > 95%
- ✅ All expected symbols present
- ✅ Correct symbol types and relationships
- ✅ No false positives or negatives

#### Success Criteria
- Accuracy meets >95% threshold
- Symbol relationships correctly identified
- No significant data loss or corruption

### Scenario 4: Error Handling and Recovery
**Test Objective**: Verify LSP error handling and recovery mechanisms

#### Test Steps
1. Setup test repository with problematic files
2. Generate symbol graph with various error conditions
3. Test recovery from partial failures
4. Verify graceful degradation
5. Monitor error logging and reporting

#### Expected Results
- ✅ Graceful handling of error conditions
- ✅ Recovery from partial failures
- ✅ Appropriate error messages
- ✅ System remains stable under errors

#### Success Criteria
- No system crashes from error conditions
- Meaningful error messages generated
- System recovers from temporary failures

### Scenario 5: Integration with caf Features
**Test Objective**: Verify LSP integration with other caf features

#### Test Steps
1. Setup test repository with 1000+ Go files
2. Initialize caf with LSP support
3. Test integration with session tracking
4. Test integration with pre-commit hooks
5. Test integration with dashboard features

#### Expected Results
- ✅ Seamless integration with session tracking
- ✅ Pre-commit hooks work with LSP data
- ✅ Dashboard displays LSP information correctly
- ✅ No conflicts between features

#### Success Criteria
- All integrations work without conflicts
- Performance impact is minimal
- User experience is seamless

## Test Execution Plan

### Phase 1: Test Environment Setup
1. **Install Dependencies**
   - Install Go 1.22+
   - Install gopls
   - Build caf binary with LSP support

2. **Create Test Repository**
   - Generate 1000+ Go files with realistic structures
   - Include various Go features and patterns
   - Create realistic project structure

3. **Initialize caf**
   - Run `caf init` in test repository
   - Verify LSP adapter detection
   - Test basic functionality

### Phase 2: Symbol Graph Generation Testing
1. **Basic Generation Test**
   - Execute symbol graph generation
   - Measure performance metrics
   - Verify output completeness

2. **Performance Testing**
   - Run multiple generation cycles
   - Monitor resource usage
   - Test with varying file sizes

3. **Accuracy Testing**
   - Validate symbol extraction
   - Check symbol relationships
   - Verify type information

### Phase 3: Integration Testing
1. **Feature Integration**
   - Test with session tracking
   - Test with pre-commit hooks
   - Test with dashboard

2. **Error Handling**
   - Test with error conditions
   - Verify recovery mechanisms
   - Test edge cases

## Success Metrics

### Performance Metrics
- **Generation Time**: < 5 minutes for 1000+ files
- **Memory Usage**: < 500MB during generation
- **CPU Usage**: < 80% during generation
- **Disk Usage**: < 1GB for symbol graph cache

### Accuracy Metrics
- **Symbol Accuracy**: > 95% accuracy rate
- **Symbol Coverage**: > 95% of expected symbols found
- **Type Accuracy**: > 95% correct type identification
- **Relationship Accuracy**: > 95% correct relationships

### Reliability Metrics
- **Success Rate**: 100% successful generation attempts
- **Error Recovery**: Graceful recovery from errors
- **Consistency**: Consistent results across multiple runs

## Test Evidence Collection

### Automated Evidence
1. **Symbol Graph Files**: Capture all generated symbol graph data
2. **Performance Logs**: Record execution time and resource usage
3. **Error Logs**: Record any errors or warnings during generation
4. **Memory Dumps**: Capture memory usage patterns

### Manual Evidence
1. **Symbol Validation**: Manual verification of symbol accuracy
2. **Performance Testing**: Manual performance testing and validation
3. **Error Scenario Testing**: Manual testing of error conditions

### Documentation
1. **Test Results**: Detailed execution results for each test
2. **Performance Reports**: Performance benchmarking data
3. **Accuracy Reports**: Symbol accuracy validation results

## Risk Assessment

### High Risk Items
1. **Memory Usage**: Large repositories may exceed memory limits
2. **Generation Time**: 1000+ files may take longer than expected
3. **Symbol Accuracy**: Complex code may reduce accuracy

### Mitigation Strategies
1. **Memory Management**: Implement memory usage monitoring and limits
2. **Performance Optimization**: Optimize symbol extraction algorithms
3. **Accuracy Validation**: Implement comprehensive validation mechanisms

## Test Deliverables

### Immediate Deliverables
1. **Symbol Graph Analysis**: Analysis of generated symbol graphs
2. **Performance Report**: Detailed performance benchmarking results
3. **Accuracy Report**: Symbol accuracy validation results

### Supporting Documentation
1. **Test Setup Guide**: Instructions for reproducing test results
2. **Performance Tuning Guide**: Optimization recommendations
3. **Best Practices**: Guidelines for LSP usage with large repositories

## Success Criteria

### Mandatory Criteria
- [ ] Symbol graph generation succeeds for 1000+ files
- [ ] Generation time < 5 minutes
- [ ] Memory usage < 500MB
- [ ] Symbol accuracy > 95%

### Performance Criteria
- [ ] Consistent performance across multiple runs
- [ ] No memory leaks detected
- [ ] CPU usage remains reasonable

### Integration Criteria
- [ ] Seamless integration with other caf features
- [ ] No conflicts between features
- [ ] User experience remains positive

## Next Steps

1. **Setup Test Environment**: Prepare Go environment and test repository
2. **Execute Symbol Generation Tests**: Run all planned test scenarios
3. **Validate Results**: Verify success criteria are met
4. **Generate Reports**: Create comprehensive test reports
5. **Optimize Performance**: Implement improvements based on results

## Conclusion

This verification plan provides a comprehensive approach to testing the LSP fallback symbol graph generation for large Go repositories. The testing will ensure that the LSP functionality meets performance and accuracy requirements for the MVP → ALPHA transition.

The LSP fallback mechanism with zero-config coverage provides a solid foundation for symbol graph generation. The testing will validate this design and identify any areas for improvement before deployment to production environments.

The successful completion of this testing will demonstrate that the LSP adapter can handle large-scale Go repositories effectively, providing valuable symbol information for AI coding assistance.