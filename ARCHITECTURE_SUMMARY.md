# Architecture Issues - Executive Summary

> For detailed analysis in Japanese, see [ARCHITECTURE_ISSUES.md](./ARCHITECTURE_ISSUES.md)

## Overview

This document identifies 7 critical architectural issues in the `wows-fast-stats` codebase and provides actionable improvement recommendations.

## Critical Issues Summary

### 1. God Object Anti-pattern 🔴 **HIGH PRIORITY**

**Problem:** `BattleFetcher` service class (631 lines) violates Single Responsibility Principle
- Handles data fetching, transformation, composition, concurrency, state management, error handling, and event emission
- Difficult to test, maintain, and extend

**Impact:** High maintenance cost, low reusability, testing difficulties

**Recommendation:** Decompose into:
- `BattleDataFetcher` (data fetching)
- `BattleDataCache` (cache management)  
- `BattleComposer` (data composition)
- `BattleService` (thin orchestration layer)

---

### 2. Manual Dependency Injection Complexity 🟡 **MEDIUM PRIORITY**

**Problem:** `DependencyContainer` manually wires all dependencies (104 lines)
- Hard to extend with new dependencies
- Difficult to partially mock for testing
- Manual initialization order management

**Impact:** Poor extensibility, testing friction

**Recommendation:** Introduce DI framework (Google Wire) or Builder pattern

---

### 3. Stateful Services & Concurrency Safety 🔴 **HIGH PRIORITY**

**Problem:** `BattleFetcher` maintains mutable state (`isFirstBattle`, cached data)
- Potential data races if called concurrently
- Order-dependent behavior
- Difficult to test

**Impact:** Concurrency bugs, unpredictable behavior

**Recommendation:** 
- Extract cache to separate thread-safe `BattleDataCache`
- Make services stateless
- Use `sync.RWMutex` for cache access

---

### 4. Poor Separation of Concerns 🟡 **MEDIUM PRIORITY**

**Problem:** `app.go` mixes UI and business logic
- Tightly coupled to Wails framework
- Cannot reuse for CLI/API versions
- Hard to unit test

**Impact:** Poor reusability, framework lock-in

**Recommendation:** Introduce Application Service layer (Use Cases) separate from UI adapters

---

### 5. Inconsistent Error Handling 🟡 **MEDIUM PRIORITY**

**Problem:** Multiple error patterns coexist
- `failure.New`, `failure.Wrap`, ignored errors, `apperr.Unwrap`
- Loss of context information
- Hard to debug

**Impact:** Poor debuggability, inconsistent error responses

**Recommendation:** Standardize on unified error handling strategy with context preservation

---

### 6. Insufficient Test Coverage 🟠 **MEDIUM-HIGH PRIORITY**

**Problem:** Critical logic lacks tests
- `BattleFetcher.compose` (120 lines) untested
- Main flow `Invoke` untested
- Concurrency edge cases untested

**Impact:** Regression risk, refactoring friction

**Recommendation:** 
- Target 80%+ coverage for business logic
- Table-driven tests for composition
- Integration tests with mock servers

---

### 7. Concurrency Complexity 🟡 **MEDIUM PRIORITY**

**Problem:** Complex goroutine and channel management
- Multiple channels for parallel fetching
- Error handling complexity
- Potential deadlocks and resource leaks

**Impact:** Hard to maintain, deadlock risk

**Recommendation:** Use `golang.org/x/sync/errgroup` for cleaner parallel execution with automatic error handling

---

## Improvement Roadmap

### Phase 1: Foundation (1-2 weeks) - **START HERE**
1. ✅ Document architecture issues (this document)
2. Set up test infrastructure
3. Standardize error handling

### Phase 2: Service Layer Refactoring (2-4 weeks)
1. Decompose `BattleFetcher`
2. Improve concurrency with `errgroup`
3. Make services stateless

### Phase 3: DI Improvement (1-2 weeks)
1. Introduce Builder pattern or Wire
2. Refactor `DependencyContainer`

### Phase 4: Layer Separation (3-4 weeks)
1. Add Application Service layer
2. Separate UI adapters
3. Clarify domain logic

### Phase 5: Quality (Ongoing)
1. Increase test coverage to 80%+
2. Documentation
3. Performance profiling

---

## Measurable Goals

### Code Quality Metrics
- Average function length: ≤50 lines
- Maximum function length: ≤200 lines
- Cyclomatic complexity: ≤15
- Test coverage: ≥80%

### Architecture Metrics
- Layer dependency violations: 0 (verified by arch-go)
- God Objects (≥300 lines): 0
- Implementations per interface: 1-2

### Maintainability Metrics
- Files changed per feature: ≤5 on average
- Build time: ≤30 seconds
- Test execution time: ≤10 seconds

---

## Benefits of Improvements

1. **Maintainability**: Easier to understand, modify, and extend
2. **Testability**: Comprehensive test coverage enables confident refactoring
3. **Scalability**: Clear architecture supports team growth
4. **Reliability**: Proper concurrency handling prevents bugs
5. **Flexibility**: Decoupled layers enable multiple UIs (CLI, Web, Desktop)

---

## Current State Assessment

✅ **Strengths:**
- Basic layer separation enforced by `arch-go`
- Interface-based design in repository layer
- Active development with clear domain model

⚠️ **Needs Improvement:**
- Service layer organization
- Testing practices
- Dependency management
- Concurrency patterns

---

## Next Steps

1. **Review this analysis** with the team
2. **Prioritize issues** based on project needs
3. **Start with Phase 1** (foundation)
4. **Implement incrementally** - small, safe changes
5. **Measure progress** against defined metrics

For detailed analysis, code examples, and Japanese version, see [ARCHITECTURE_ISSUES.md](./ARCHITECTURE_ISSUES.md).
