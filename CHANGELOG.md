# CHANGELOG

All notable changes to GoTiny are documented in this file.

## [0.4.1] - 2026-09-25

### Added

- Add support for single-line comments using `#`. A comment can start anywhere on a line and extends to the end of that line.
- Add integration test scripts with comments.

## [0.4.0] - 2026-09-25

### Added

- Function types can now be expressed in source syntax.
  - Functions can return functions.
  - Functions can accept functions as arguments.
  - Function types support nested function types.

  For example:

  ```gotiny
  fn fibonacci() fn() Int {
      ...
  }

  fn apply(f fn(Int) Int, x Int) Int {
      return f(x);
  }
  ```

### Changed

- Compare function types structurally during static type checking.

## [0.3.1] - 2026-09-25

### Added

- Add a `--version` flag to display the current version.
- Add an MIT License and CHANGELOG file.

### Changed

- Improve file error diagnostics with standard formatted output.
- Add golangci-lint configuration and standardize code formatting.

## [0.3.0] - 2026-09-24

### Added

- Add function types with parameters and return types, including implicit `VoidType`.
- Add function declarations and function calls.
- Add return statements with propagation through nested blocks using a `ReturnValue`.
- Add closures with lexical environment capture.
- Allow recursive function calls.
- Add static checking for function declarations and calls.
- Add runtime evaluation of function declarations and calls.

### Changed

- Refactor the runtime environment into a separate package.
- Separate the type environment and static checker into their own package.
- Refactor the type system to use a Type interface.
- Avoid printing nil evaluation results.
- Polish and update the README.

## [0.2.0] - 2026-09-22

### Added

- Add scoped blocks and if statement evaluation.
- Add static checking for blocks and if statements.
- Add scoped type environments.

### Changed

- Add location token helper to expressions to support diagnostics.

## [0.1.0] - 2026-09-21

### Added

- Add a basic CLI with a REPL and file runner.
- Add static type checking.
- Add boolean operators with short-circuit evaluation.
- Add comparison and equality operators.
- Add variable declarations and assignments with type checking.
- Add int and bool runtime types and values.
- Add support for multi-statement programs.
- Add semicolon statement delimiters.Add variable expressions and runtime environment evaluation.
- Add unary plus and minus operators.
- Add if/else statements and code blocks.
- Add parse and evaluation error types.
- Add line and column information to tokens.

### Changed

- Organize the project into cmd and internal packages.
- Introduce a statement layer around expressions.
- Preserve name tokens in declaration and assignment statements.
