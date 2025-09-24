# Tasks: File Upload Web Application

**Input**: Design documents from `/specs/001-file-upload/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/http-api.md, quickstart.md, CLAUDE.md

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Extract: Go 1.21+, standard library only, Docker/Helm deployment
2. Load design documents:
   → data-model.md: File upload entity → storage implementation
   → contracts/http-api.md: 3 endpoints → contract tests
   → research.md: Technology decisions → setup tasks
   → quickstart.md: Test scenarios → integration tests
3. Generate tasks by category:
   → Setup: Go project init, test harness
   → Tests: Contract tests for all endpoints
   → Core: HTTP server with embedded HTML
   → Docker: Multi-stage build, GitHub Actions
   → Kubernetes: Helm chart creation
   → Polish: Documentation, manual testing
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file (main.go) = sequential
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001-T023)
6. Validate task completeness
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Phase 3.1: Setup
- [ ] T001 Create Go project structure with go.mod (module: github.com/stianfro/file-upload-web)
- [ ] T002 [P] Create test harness script test.sh for manual testing
- [ ] T003 [P] Create .gitignore for Go project (binaries, uploads/)

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T004 [P] Contract test GET / (HTML form) in tests/contract_index_test.go
- [ ] T005 [P] Contract test POST /upload (file upload) in tests/contract_upload_test.go
- [ ] T006 [P] Contract test GET /health (health check) in tests/contract_health_test.go
- [ ] T007 [P] Integration test browser upload flow in tests/integration_browser_test.go
- [ ] T008 [P] Integration test curl upload scenarios in tests/integration_curl_test.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T009 Create main.go with embedded HTML template (structure per CLAUDE.md)
- [ ] T010 Implement GET / handler serving embedded HTML form in main.go
- [ ] T011 Implement POST /upload handler with file streaming in main.go
- [ ] T012 Implement GET /health handler for K8s probes in main.go
- [ ] T013 Add filename sanitization function in main.go
- [ ] T014 Add environment variable configuration (PORT, UPLOAD_DIR, MAX_SIZE) in main.go
- [ ] T015 Embed index.html template with form and basic CSS in main.go

## Phase 3.4: Containerization
- [ ] T016 Create multi-stage Dockerfile with Alpine base (<20MB target)
- [ ] T017 [P] Create .dockerignore to exclude unnecessary files
- [ ] T018 [P] Create docker-compose.yml for local testing with volume mount

## Phase 3.5: CI/CD Pipeline
- [ ] T019 Create .github/workflows/build.yml for automated Docker builds
- [ ] T020 Configure multi-arch build (amd64, arm64) in GitHub Actions
- [ ] T021 Add GHCR push with version tagging in workflow

## Phase 3.6: Kubernetes Deployment
- [ ] T022 [P] Create charts/file-upload-web/Chart.yaml with OCI support
- [ ] T023 [P] Create charts/file-upload-web/values.yaml with sensible defaults
- [ ] T024 [P] Create charts/file-upload-web/templates/deployment.yaml
- [ ] T025 [P] Create charts/file-upload-web/templates/service.yaml
- [ ] T026 [P] Create charts/file-upload-web/templates/pvc.yaml for persistence

## Phase 3.7: Polish
- [ ] T027 Create README.md with purpose, quick start, curl examples
- [ ] T028 [P] Add troubleshooting section to README.md
- [ ] T029 Run manual test script against all endpoints
- [ ] T030 Verify Docker image size < 20MB
- [ ] T031 Test Helm chart deployment locally with kind/k3s

## Dependencies
- Tests (T004-T008) must complete before implementation (T009-T015)
- main.go implementation (T009-T015) is sequential - same file
- T009 blocks T010-T015 (all modify main.go)
- T016 requires T009-T015 (needs complete app to build)
- T019-T021 require T016 (needs Dockerfile)
- T022-T026 can run parallel (different Helm files)
- Polish tasks depend on all implementation

## Parallel Execution Examples

### Test Phase (T004-T008)
```bash
# Launch all contract and integration tests together:
Task subagent_type=general-purpose prompt="Create contract test for GET / endpoint verifying HTML form response in tests/contract_index_test.go"
Task subagent_type=general-purpose prompt="Create contract test for POST /upload endpoint in tests/contract_upload_test.go"
Task subagent_type=general-purpose prompt="Create contract test for GET /health endpoint in tests/contract_health_test.go"
Task subagent_type=general-purpose prompt="Create integration test for browser upload flow in tests/integration_browser_test.go"
Task subagent_type=general-purpose prompt="Create integration test for curl upload scenarios in tests/integration_curl_test.go"
```

### Helm Chart Creation (T022-T026)
```bash
# Create all Helm chart files in parallel:
Task subagent_type=general-purpose prompt="Create Helm Chart.yaml with OCI support at charts/file-upload-web/Chart.yaml"
Task subagent_type=general-purpose prompt="Create Helm values.yaml with config options at charts/file-upload-web/values.yaml"
Task subagent_type=general-purpose prompt="Create Kubernetes deployment template at charts/file-upload-web/templates/deployment.yaml"
Task subagent_type=general-purpose prompt="Create Kubernetes service template at charts/file-upload-web/templates/service.yaml"
Task subagent_type=general-purpose prompt="Create PVC template for persistence at charts/file-upload-web/templates/pvc.yaml"
```

## Task-Specific Notes

### Core Implementation (T009-T015)
- All modifications to main.go must be sequential
- Use embed directive for HTML: `//go:embed index.html`
- Implement streaming for large files per CLAUDE.md guidance
- Maximum 200 lines total to maintain simplicity

### Docker Build (T016)
- Multi-stage build required
- Build stage: `golang:alpine`
- Runtime stage: `alpine:latest` or `scratch`
- Non-root user execution
- Target size: < 20MB

### GitHub Actions (T019-T021)
- Trigger on: push to main, version tags
- Use GITHUB_TOKEN for GHCR authentication
- Tag pattern: `latest` and semantic version

### Helm Chart (T022-T026)
- No complex templating (KISS principle)
- Support both NodePort and LoadBalancer
- Optional PVC for production use
- Health/liveness probes configured

## Validation Checklist
- [x] All 3 endpoints from contracts have tests (T004-T006)
- [x] File entity from data-model has storage implementation (T011, T013)
- [x] All quickstart scenarios covered (browser T007, curl T008)
- [x] Tests come before implementation (Phase 3.2 before 3.3)
- [x] Parallel tasks use different files
- [x] Sequential tasks (T009-T015) all modify main.go
- [x] Each task specifies exact file path
- [x] Docker and Helm deployment included per plan.md
- [x] Manual testing and documentation tasks included

## Success Criteria
- Browser file upload works with drag-and-drop
- curl upload works as documented
- Docker image under 20MB
- Helm chart deploys successfully
- All tests pass
- Zero external dependencies in go.mod