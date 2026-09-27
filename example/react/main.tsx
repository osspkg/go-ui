import { createRoot } from "react-dom/client";
import { useState, type ComponentType, type ReactElement } from "react";
import { BasicExample } from "./basic.js";
import { DashboardExample } from "./dashboard.js";
import { HTMLFormExample } from "./html-form.js";
import "./preview.css";

const examples = {
  Basic: BasicExample,
  "HTML form": HTMLFormExample,
  Dashboard: DashboardExample,
} satisfies Record<string, ComponentType>;

type ExampleName = keyof typeof examples;
const exampleNames = Object.keys(examples) as ExampleName[];

function ExamplesApp(): ReactElement {
  const [name, setName] = useState<ExampleName>("Basic");
  const Example = examples[name];

  return (
    <main className="min-h-screen bg-ui-muted px-4 py-8 text-ui-foreground sm:px-8 sm:py-12">
      <section className="mx-auto flex w-full max-w-6xl flex-col gap-8">
        <header className="flex flex-col gap-2">
          <p className="text-sm font-medium text-ui-muted-foreground">@osspkg/ui-react</p>
          <h1 className="text-3xl font-semibold tracking-tight sm:text-4xl">Interactive examples</h1>
          <p className="max-w-2xl text-ui-muted-foreground">Preview declarative schemas rendered with the default HTML and shadcn component registries.</p>
        </header>
        <nav aria-label="Examples" className="flex flex-wrap gap-2 rounded-xl border border-ui-border bg-ui-card p-2 shadow-sm">
          {exampleNames.map((exampleName) => (
            <button
              key={exampleName}
              type="button"
              aria-pressed={name === exampleName}
              className="min-h-11 min-w-11 rounded-md border border-ui-border bg-ui-card px-3 py-2 text-sm font-medium transition-colors hover:bg-ui-muted aria-pressed:border-ui-primary aria-pressed:bg-ui-primary aria-pressed:text-ui-primary-foreground"
              onClick={() => setName(exampleName)}
            >
              {exampleName}
            </button>
          ))}
        </nav>
        <section aria-label={`${name} preview`} className="rounded-xl border border-ui-border bg-ui-card p-4 shadow-sm sm:p-6">
          <Example />
        </section>
      </section>
    </main>
  );
}

const rootElement = document.getElementById("root");
if (!rootElement) throw new Error("Missing React root element");

createRoot(rootElement).render(<ExamplesApp />);
