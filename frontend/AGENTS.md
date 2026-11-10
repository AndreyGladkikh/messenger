# Frontend engineering guidelines

These instructions apply to all work under `frontend/` and supplement the repository's `AGENTS.md`. Follow Feature-Sliced Design (FSD) for new code and changes to existing code.

## Before making changes

- Identify the business domain, bounded context, and owning slice for each concept. In this messenger, distinguish authentication/identity concerns from messaging concerns. FSD slices organize frontend responsibilities; they do not need to mirror backend packages or bounded contexts one-to-one.
- Inspect the affected slice, its consumers, and existing public APIs before creating abstractions. Apply SOLID within these boundaries.
- Research relevant approaches using current official documentation before implementation. Prefer the existing React, TypeScript, Vite, React Router, Zustand, and TanStack Query stack unless the task warrants a change.
- Start with the page that owns the behavior. Keep page-specific UI, state, and requests together; extract to lower layers when there is a concrete reuse or ownership need. Do not create a feature for every component or interaction. See the [FSD pages-first guidance](https://feature-sliced.design/docs/guides/migration/from-v2-0).

## Layers and ownership

The FSD root is `src/`. Dependencies flow downward in this order; layers may be skipped:

```text
app → pages → widgets → features → entities → shared
```

| Layer | Responsibility in this frontend |
| --- | --- |
| `app` | Application assembly: routing, providers, startup, global styles, and application-wide orchestration. |
| `pages` | Route screens and their local behavior. |
| `widgets` | Optional independent UI blocks combining lower layers, when page-local composition is insufficient. |
| `features` | Reusable user actions and use cases. |
| `entities` | Reusable domain concepts such as a user, chat, or message, including their models and representations. |
| `shared` | Infrastructure and reusable UI without product business rules: transport, configuration, UI primitives, focused libraries. |

Only introduce a missing layer when needed. Do not add custom layers or the deprecated `processes` layer. Keep `src/main.tsx` as the bootstrap entry; place new application assembly in `app`. These responsibilities follow the [FSD layer definitions](https://feature-sliced.design/docs/reference/layers).

## Slices and segments

- Use `src/<layer>/<slice>/<segment>/` for `pages`, `widgets`, `features`, and `entities`. Name slices after their domain or use case, using lowercase kebab-case for new names.
- `app` and `shared` contain segments directly, without business slices.
- Use purpose-based segments: `ui` for rendering/styles, `model` for state/validation/business rules, `api` for requests/DTOs/mapping, `lib` for focused internal helpers, and `config` for configuration.
- Create only the segments needed. Keep related hooks, types, tests, and assets with their owner; do not introduce top-level `components`, `hooks`, `types`, `services`, or `utils` buckets.
- Slice groups, if needed, are organizational folders only. They must not contain shared implementation or a group-wide public API.

Follow the [FSD slice and segment rules](https://feature-sliced.design/docs/reference/slices-segments).

## Dependency boundaries

- A slice may import its own internals and public APIs from any lower layer. It must not import an upper layer or another slice on its own layer. This includes type-only imports, re-exports, and dynamic imports.
- Segments within `app` may depend on one another; the same applies within `shared`. This does not permit upward dependencies or cycles.
- Compose independent slices in a higher layer. Pass data, callbacks, or UI slots downward instead of making a lower layer import its consumers. Never move business rules into `shared` merely to bypass isolation.
- For unavoidable entity relationships only, expose a minimal consumer-specific API at `entities/<provider>/@x/<consumer>.ts`. Document the relationship and keep it acyclic. Do not use `@x` for features, widgets, or pages.

## Public APIs and imports

- Give each slice a root `index.ts` with explicit named exports of its supported interface. Export only what consumers need; avoid `export *` and layer-wide barrels such as `features/index.ts`.
- Cross-boundary imports use `@/` (mapped to `src/`) and public entry points. Within a slice, use relative paths to implementation files, never its own root barrel.

These rules and the entity-only exception follow the [FSD public API guidance](https://feature-sliced.design/docs/reference/public-api).

```ts
// A page consumes the feature's public API (add the entry point if missing).
import { CredentialsForm } from '@/features/auth';

// Inside features/auth/ui/CredentialsForm.tsx, use the slice's internal module.
import type { Credentials } from '../model/credentials';

// Do not couple external consumers to feature internals.
// import { CredentialsForm } from '@/features/auth/ui/CredentialsForm';
```

## Applying the rules to this codebase

- Existing pages lack `ui` segments and slice-root public APIs, and several consumers use deep imports. Treat these as migration work, not examples to copy. When adding a consumer, expose the required API first. Correct boundaries in the area being changed and update affected consumers together; avoid unrelated wholesale moves.
- Keep generic HTTP mechanics in `shared/api`. Put use-case orchestration, domain validation, and state in their owning page, feature, or entity. Low-level endpoint wrappers may remain in `shared/api` when they contain no business policy; keep domain interpretation with its owner.
- Keep query definitions, cache keys, and mutations with the behavior/data they describe. Shared query-client infrastructure must not import page or feature code. Wire router/provider dependencies in `app` through lower-layer public APIs.
- Reuse the existing UI primitives and respect `components.json` when generating components. Inspect generated code: its default folders and imports do not override these boundaries. Authentication or messaging behavior belongs in the owning slice, even when built from shared UI.
- Record any unresolved boundary violation touched by the task and the reason it remains. Existing violations do not authorize new ones.

## Verification

- Review every added or changed import for layer direction, slice isolation, public API use, and cycles. The current ESLint configuration does not enforce FSD; passing lint alone does not establish architectural correctness.
- For frontend code changes, run `npm run lint` and `npm run build` from `frontend/`, plus tests covering the changed behavior. Add or update tests as required by the root instructions; keep unit/component tests with the owning code and integration tests focused on public contracts.
- There is currently no frontend test script or runner. When changing executable behavior, provide appropriate runnable test coverage and document its command/setup; do not claim lint or compilation replaces tests.
- The repository's `make test` command runs backend Go tests; backend end-to-end tests are in `backend/test/e2e`. It does not validate frontend behavior. Use it when backend/API integration is affected.
- For documentation-only changes, check accuracy, links, and the diff; application tests are unnecessary. Report checks performed and any limitations.
