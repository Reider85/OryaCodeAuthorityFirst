# Prompt 10 Setup: Codecov + Coverage Badge

## Current Status
- ✅ Codecov integration in CI workflow
- ✅ Coverage badge in README
- ❌ Token needs to be updated (currently PLACEHOLDER_TOKEN)

## Steps to Complete

### 1. Get Codecov Token
1. Go to [codecov.io](https://codecov.io)
2. Sign in with GitHub
3. Select repository `Reider85/OryaCodeAuthorityFirst`
4. Go to **Settings** → **General**
5. Find **Repository Token** under "Token"
6. Copy the token (starts with `...`)

### 2. Add GitHub Secret
1. Go to repository → **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Name: `CODECOV_TOKEN`
4. Paste the token from step 1
5. Click **Add secret**

### 3. Update README Badge
Replace the placeholder token in README.md:

**Current:**
```markdown
[![coverage](https://codecov.io/gh/Reider85/OryaCodeAuthorityFirst/graph/badge.svg?token=PLACEHOLDER_TOKEN)](https://codecov.io/gh/Reider85/OryaCodeAuthorityFirst)
```

**After:**
```markdown
[![coverage](https://codecov.io/gh/Reider85/OryaCodeAuthorityFirst/graph/badge.svg)](https://codecov.io/gh/Reider85/OryaCodeAuthorityFirst)
```

## Verification Steps

After setup:
1. [ ] Codecov token is added as GitHub secret
2. [ ] README badge shows actual coverage percentage
3. [ ] CI workflow uploads coverage to Codecov
4. [ ] Codecov dashboard shows coverage trends

## Expected Impact
- **Transparency**: Public coverage tracking
- **Trust**: Visible quality metrics
- **Trend Analysis**: Historical coverage data
- **Badge**: Real-time coverage percentage in README

## Files to Update
- `README.md` - Remove placeholder from coverage badge

## Notes
- The token is automatically used by the existing CI workflow
- Coverage upload is already configured in `.github/workflows/ci.yml`
- No additional workflow changes needed