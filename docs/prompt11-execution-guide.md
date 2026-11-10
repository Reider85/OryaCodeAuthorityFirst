# Prompt 11 Execution Guide

## Goal
Configure Issue/PR templates and enable Discussions to lower barrier for community contributions.

## Current Status
- ✅ All templates already created and configured
- ✅ FUNDING.yml for sponsor support
- ❌ Discussions need manual enablement

## Execution Steps

### Step 1: Enable Discussions (Manual Setup)
1. **Repository**: Settings → General
2. **Features**: Find "Discussions"
3. **Toggle**: "Set up discussions"
4. **Categories**:
   - **Q&A** - Questions, help requests
   - **Ideas** - Feature suggestions, RFCs
   - **Show & tell** - Project updates, demos

### Step 2: Verify Templates (Already Complete)
- **Issue Templates**: GitHub Forms configured
  - Bug report: Environment, steps, logs
  - Feature request: Problem/solution/use-case
  - Config: Links to Discussions, no blank issues
- **PR Template**: Standard sections included
  - Summary, motivation, changes, type, testing, checklist
- **FUNDING.yml**: GitHub sponsor integration

### Step 3: Test Functionality
1. **Create Test Issue**:
   - Use "New Issue" → "Bug report" template
   - Fill out form, submit
   - Verify correct labels and structure
2. **Create Test PR**:
   - Create branch, make changes
   - Use PR template
   - Verify checklist and sections appear
3. **Test Discussions**:
   - Navigate to "Discussions" tab
   - Verify categories appear
   - Create test post in Q&A

## Template Features

### Issue Templates
- **Bug Report**: Comprehensive with environment details, reproduction steps, logs
- **Feature Request**: Structured with problem statement, proposed solution, use cases
- **Configuration**: Directs users to Discussions, disables blank issues

### PR Template
- **Structure**: Summary → Motivation → Changes → Type → Testing → Checklist
- **Categories**: Bug fix, feature, breaking change, docs, perf, refactor, test, chore
- **Quality**: Self-review, testing, documentation, style checks

### Community Features
- **FUNDING.yml**: GitHub sponsor button
- **Discussions**: Three categories for community engagement
- **Issue Templates**: Lower barrier for bug reports and features

## Expected Impact
- **Lower Friction**: Easy issue/PR creation with templates
- **Consistency**: Standardized contribution formats
- **Community**: Discussion-based engagement
- **Sustainability**: Sponsor support option

## Verification Checklist
- [ ] Discussions enabled in repository settings
- [ ] Three discussion categories configured
- [ ] Issue templates work correctly
- [ ] PR template works correctly
- [ ] Sponsor button visible on repository page
- [ ] Test issues/PRs created successfully

## Files Created/Verified
- ✅ `.github/ISSUE_TEMPLATE/bug_report.yml`
- ✅ `.github/ISSUE_TEMPLATE/feature_request.yml`
- ✅ `.github/ISSUE_TEMPLATE/config.yml`
- ✅ `.github/PULL_REQUEST_TEMPLATE.md`
- ✅ `.github/FUNDING.yml`

## Next Steps
1. Complete Discussions enablement in GitHub settings
2. Test template functionality with sample issue/PR
3. Verify sponsor button appears
4. Monitor community engagement