# AGENTS.md

## Purpose
This file defines how you should work within this repo.
The goal is to help you make changes that allign with what the maintainers consider to be "good code"
## Rules
Fundamentally, your work should follow the principles outlined by the article "The Grug Brained Developer A layman's guide to thinking like the self-aware smol brained (https://grugbrain.dev/)". Essentially, avoid complexity and keep things simple!
- Minimize the number of new data structures. Prefer primitives, native data structures, and existing datastructures to creating new ones in that order.
- Strive towards creating pure functions. Prefer deterministic methods with no I/O, no hidden global state, and no mutation of receiver state.
- Method size is not relevant. What is relevant is purpose fit, testability, maintainability, and readability in that order.