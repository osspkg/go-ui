import { createComponentRegistry, registerDefaultHTMLComponents, registerDefaultShadcnComponents, UIProvider, ViewRenderer } from "@osspkg/ui-react";
import type { ViewSchema } from "@osspkg/ui-core";
import "@osspkg/ui-react/styles.css";

const components = createComponentRegistry();
registerDefaultHTMLComponents(components);
registerDefaultShadcnComponents(components);

const schema: ViewSchema = {
  protocolVersion: "1.0",
  id: "profile.edit",
  state: { submitted: false },
  actions: {
    save: { type: "set-state", path: "submitted", value: true },
  },
  regions: {
    "top-header": [],
    "left-panel": [],
    "right-panel": [],
    bottom: [],
    content: [
      {
        id: "form-card",
        component: "article",
        layout: { row: 1, cols: 8, offset: 2 },
        children: [
          { id: "heading", component: "h1", layout: { row: 1, cols: 12 }, props: { text: "Edit profile" } },
          {
            id: "profile-form",
            component: "form",
            events: { submit: { action: "save" } },
            children: [
              { id: "name-label", component: "label", layout: { row: 1, cols: 12 }, props: { text: "Name" } },
              { id: "name", component: "input", layout: { row: 2, cols: 8, offset: 2 }, props: { name: "name", placeholder: "Ada Lovelace" } },
              { id: "save", component: "button", layout: { row: 3, cols: 8, offset: 2 }, props: { type: "submit", text: "Save profile" } },
            ],
          },
          {
            id: "confirmation",
            component: "p",
            when: { $eq: [{ $state: "submitted" }, true] },
            layout: { row: 4, cols: 8, offset: 2 },
            props: { text: "Profile saved." },
          },
        ],
      },
    ],
  },
};

export function HTMLFormExample() {
  return <UIProvider components={components}><ViewRenderer schema={schema} /></UIProvider>;
}
