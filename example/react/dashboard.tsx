import { createComponentRegistry, registerDefaultHTMLComponents, registerDefaultShadcnComponents, UIProvider, ViewRenderer } from "@osspkg/ui-react";
import type { ViewSchema } from "@osspkg/ui-core";
import "@osspkg/ui-react/styles.css";

const components = createComponentRegistry();
registerDefaultHTMLComponents(components);
registerDefaultShadcnComponents(components);

const schema: ViewSchema = {
  protocolVersion: "1.0",
  id: "metrics.dashboard",
  regions: {
    "top-header": [
      { id: "title", component: "h1", layout: { row: 1, cols: 12 }, props: { text: "Metrics dashboard" } },
    ],
    "left-panel": [],
    "right-panel": [],
    bottom: [],
    content: [
      { id: "requests", component: "card", layout: { row: 1, cols: 6 }, props: { text: "Requests: 12,480" } },
      { id: "latency", component: "card", layout: { row: 1, cols: 6, offset: 6 }, props: { text: "p95 latency: 184 ms" } },
      {
        id: "events",
        component: "data-table",
        layout: { row: 2, cols: 12 },
        props: { rows: [{ event: "user.created", status: "success" }, { event: "user.deleted", status: "success" }] },
      },
    ],
  },
};

export function DashboardExample() {
  return <UIProvider components={components}><ViewRenderer schema={schema} /></UIProvider>;
}
