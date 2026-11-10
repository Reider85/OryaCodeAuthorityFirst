# Prompt 6 Execution Guide

## Goal
Fill topics, description, and social preview to improve GitHub discoverability and social sharing.

## Steps to Execute

### 1. Update README Description ✅
- **Status**: Complete
- **Changes**: Added clear, concise description under title
- **Impact**: Better first impression for visitors

### 2. Create Social Preview SVG ✅
- **Status**: Complete
- **File**: `docs/social-preview.svg`
- **Specs**: 1280×640, dark theme, includes title, subtitle, and key features
- **Usage**: Upload to GitHub repository settings

### 3. Manual Setup Required 📋
These must be done through GitHub web interface:

#### A. Repository Description
1. Go to repository → **Settings** → **About**
2. Update **Description** field:
   ```
   CodeAuthorityFirst (caf) — single-binary Go tool that records provenance for AI coding sessions: which agent touched which file, from which prompt, through which iteration. Offline, no LLM calls, fail-safe pre-commit hooks.
   ```

#### B. Topics
1. Go to **Settings** → **About** → **Tags**
2. Add exactly 20 topics:
   ```
   go, golang, cli, devtool, developer-tools, ai-coding, ai-agents, provenance, traceability, audit, developer-experience, dx, single-binary, pre-commit, tui, source-code, code-review, observability, llm-tools, agent-ops
   ```

#### C. Social Preview
1. Go to **Settings** → **Social preview**
2. Upload `docs/social-preview.svg`
3. Verify preview looks correct

#### D. Website
1. Go to **Settings** → **About** → **Website**
2. Set URL: `https://github.com/Reider85/OryaCodeAuthorityFirst#readme`

### 4. Expected Impact
- **SEO**: Better ranking in GitHub topic searches
- **Social**: Professional appearance when shared on Twitter/LinkedIn
- **Discovery**: More visible in GitHub's topic pages
- **Trust**: Complete, professional repository profile

### 5. Files Modified
- `README.md` - Enhanced description
- `docs/social-preview.svg` - Social media card
- `docs/prompt6-checklist.md` - Setup instructions

### 6. Verification
After manual setup, verify:
- [ ] Description appears under repository title
- [ ] 20 topics are visible on repository page
- [ ] Social preview works when sharing link
- [ ] Website link is functional

## Next Steps
1. Complete manual setup in GitHub repository settings
2. Test social sharing to verify preview
3. Monitor repository traffic from topic searches