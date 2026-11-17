# React + TypeScript + Vite

## Browser tests

From `frontend`, install the dependencies and Playwright's Chromium browser once:

```sh
pnpm install --frozen-lockfile
pnpm exec playwright install chromium
```

Run the browser tests with `pnpm test`. Two browser workers keep resource use
bounded; override with `pnpm test --workers=4` when appropriate.
The tests start their own Vite server at
`http://127.0.0.1:4173`, so that port must be free. They mock the authentication,
profile, and media HTTP endpoints; no running backend or real credentials are needed.
The test server overrides `VITE_API_URL` to use the same origin's `/api` prefix.

Failure screenshots and traces are saved in `messenger-playwright-results` inside the
system temporary directory (normally `/tmp`), outside the repository. Use
`npm test -- --headed` to watch the tests or `npm test -- --debug` to step through
them.

## Profile onboarding

After sign-up, `/onboarding/profile` collects a required display name and username.
The router checks `GET /me/profile` on authenticated visits; a 404 requires
onboarding, and other errors keep the messenger inaccessible until retried.
Successful creation is confirmed by reading the profile before navigating home.

Avatars are initialized through `POST /media/uploads`, uploaded directly to the
returned presigned URL, and finalized through the upload completion endpoint.
Profile creation then sends the resulting `avatarMediaId` as JSON. PNG, JPEG, and
WebP uploads are limited to 2 MiB by the profile form. The object-storage endpoint
must be reachable by the browser and allow `PUT` requests from the frontend origin
through its CORS configuration.

Run `pnpm test` for browser coverage of onboarding, session guards, media uploads,
validation, retry behavior, and the saved profile in the account menu. `make test`
also runs backend integration tests and requires Docker for their dependencies.
The route gate restricts frontend access; it does not add a server-side profile
authorization policy to the existing messaging endpoints.

This template provides a minimal setup to get React working in Vite with HMR and some ESLint rules.

Currently, two official plugins are available:

- [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react) uses [Oxc](https://oxc.rs)
- [@vitejs/plugin-react-swc](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react-swc) uses [SWC](https://swc.rs/)

## React Compiler

The React Compiler is enabled on this template. See [this documentation](https://react.dev/learn/react-compiler) for more information.

Note: This will impact Vite dev & build performances.
You can also try [the experimental native React Compiler support in plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md#rust-react-compiler) by using `compiler: true` in the plugin options instead of using the Babel plugin.

## Expanding the ESLint configuration

If you are developing a production application, we recommend updating the configuration to enable type-aware lint rules:

```js
export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...

      // Remove tseslint.configs.recommended and replace with this
      tseslint.configs.recommendedTypeChecked,
      // Alternatively, use this for stricter rules
      tseslint.configs.strictTypeChecked,
      // Optionally, add this for stylistic rules
      tseslint.configs.stylisticTypeChecked,

      // Other configs...
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])

```

You can also install [eslint-plugin-react-x](https://npmx.dev/package/eslint-plugin-react-x) and [eslint-plugin-react-dom](https://npmx.dev/package/eslint-plugin-react-dom) for React-specific lint rules:

```js
// eslint.config.js
import reactX from 'eslint-plugin-react-x'
import reactDom from 'eslint-plugin-react-dom'

export default defineConfig([
  globalIgnores(['dist']),
  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      // Other configs...
      // Enable lint rules for React
      reactX.configs['recommended-typescript'],
      // Enable lint rules for React DOM
      reactDom.configs.recommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ['./tsconfig.node.json', './tsconfig.app.json'],
        tsconfigRootDir: import.meta.dirname,
      },
      // other options...
    },
  },
])

```
