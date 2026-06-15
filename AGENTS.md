# AGENTS.md

## Purpose
This file defines how you should work within this repo.
The goal is to help you make changes that align with what the maintainers consider to be "good code".

## Philosophy
The article "The Grug Brained Developer A layman's guide to thinking like the self-aware smol brained" (https://grugbrain.dev/) is the underlying inspiration. In this repo, the guiding principle is simple: avoid unnecessary complexity and prefer code that is easy to read, test, and change.

## Rules
- Minimize the number of new data structures. Prefer primitives, native data structures, and existing data structures to creating new ones in that order.
- Strive towards creating pure functions. Prefer deterministic methods with no I/O, no hidden global state, and no mutation of receiver state.
- Prefer cohesive methods that are purpose fit, testable, maintainable, and readable in that order. Method size is secondary to those qualities.
- Avoid recursion unless it is clearly simpler and bounded.
- Keep call chains shallow.
