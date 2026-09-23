# React examples

These TypeScript entrypoints are intended for a React application that has
`@osspkg/ui-react` and `@osspkg/ui-core` installed. Each imports the packaged
Tailwind v4 stylesheet and uses only the public API.

The directory is a private workspace package, so React and its JSX runtime are
installed for the examples themselves. Check it with:

```bash
cd react
pnpm --filter @osspkg/ui-examples typecheck
```

To browse the examples locally, run this from the repository root:

```bash
make examples-dev
```

Vite prints the local address (normally http://localhost:5173). Use the
buttons at the top of the page to switch between the examples.

- [`basic.tsx`](basic.tsx) renders a static schema.
- [`html-form.tsx`](html-form.tsx) demonstrates native HTML, local state, and a
  submit action.
- [`dashboard.tsx`](dashboard.tsx) demonstrates shadcn defaults and the
  12-column `row`/`cols`/`offset` layout.

Render one exported component from an application entrypoint, for example:

```tsx
import { createRoot } from "react-dom/client";
import { DashboardExample } from "./dashboard";

createRoot(document.getElementById("root")!).render(<DashboardExample />);
```
