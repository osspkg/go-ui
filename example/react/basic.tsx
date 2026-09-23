import { createComponentRegistry, registerDefaultHTMLComponents, registerDefaultShadcnComponents, UIProvider, ViewRenderer } from "@osspkg/ui-react";
import type { ViewSchema } from "@osspkg/ui-core";
import "@osspkg/ui-react/styles.css";

const components = createComponentRegistry();
registerDefaultHTMLComponents(components);
registerDefaultShadcnComponents(components);

const schema: ViewSchema = {
  protocolVersion: "1.0",
  id: "welcome",
  title: "Welcome",
  regions: {
    "top-header": [],
    "left-panel": [],
    "right-panel": [],
    bottom: [],
    content: [
      {
        id: "welcome-card",
        component: "card",
        layout: { row: 1, cols: 8, offset: 2 },
        children: [
          { id: "title", component: "h1", props: { text: "Welcome" } },
          { id: "body", component: "p", props: { text: "This view is rendered from a declarative schema." } },
        ],
      },
    ],
  },
};

export function BasicExample() {
  return <UIProvider components={components}><ViewRenderer schema={schema} /></UIProvider>;
}
