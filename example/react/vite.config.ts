import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";

const source = (path: string) => fileURLToPath(new URL(path, import.meta.url));

export default defineConfig({
  resolve: {
    alias: [
      { find: "@osspkg/ui-react/styles.css", replacement: source("../packages/ui-react/src/styles.css") },
      { find: "@osspkg/ui-core", replacement: source("../packages/ui-core/src/index.ts") },
      { find: "@osspkg/ui-react", replacement: source("../packages/ui-react/src/index.ts") },
      { find: "@osspkg/ui-transport", replacement: source("../packages/ui-transport/src/index.ts") },
    ],
  },
});
