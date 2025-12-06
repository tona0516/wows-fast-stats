# Architecture Documentation

This directory contains architecture documentation for the wows-fast-stats project.

## Documents

- [Architecture Issues](../ARCHITECTURE_ISSUES.md) - Detailed analysis of architectural problems and improvement proposals (日本語)
- [Architecture Matrix](../ARCHITECTURE_MATRIX.md) - Priority matrix and ROI analysis (日本語)
- [Architecture Summary](../ARCHITECTURE_SUMMARY.md) - Executive summary (English)

## Overview

The architecture documentation provides:

1. **Comprehensive analysis** of 7 major architectural problems
2. **Current state** code examples from the actual codebase
3. **Concrete refactoring proposals** with code samples
4. **Prioritized roadmap** with 5 implementation phases
5. **Measurable improvement targets** and metrics
6. **ROI analysis** to guide implementation priority

## Quick Links

### Main Issues Identified

1. **God Object Anti-pattern** - `BattleFetcher` (631 lines) with multiple responsibilities
2. **Manual Dependency Injection Complexity** - Hard to extend and test
3. **Stateful Services** - Concurrent safety concerns
4. **Separation of Concerns** - Layer boundaries are unclear
5. **Inconsistent Error Handling** - Multiple error patterns mixed
6. **Test Coverage Gaps** - Critical logic untested
7. **Concurrency Complexity** - Goroutine and channel management

### Recommendations

See [ARCHITECTURE_ISSUES.md](../ARCHITECTURE_ISSUES.md) for detailed proposals including:
- Service layer decomposition strategy
- Dependency injection improvements
- Stateless architecture patterns
- Application layer introduction
- Error handling standardization
- Testing strategy
- Concurrency simplification with errgroup
