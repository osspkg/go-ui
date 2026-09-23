import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";

const source = (path: string) => fileURLToPath(new URL(path, import.meta.url));

export default defineConfig({
  resolve: {
    alias: [
      { find: "@osspkg/ui-react/styles.css", replacement: source("../../react-sdk/ui-react/src/styles.css") },
      { find: "@osspkg/ui-core", replacement: source("../../react-sdk/ui-core/src/index.ts") },
      { find: "@osspkg/ui-react", replacement: source("../../react-sdk/ui-react/src/index.ts") },
      { find: "@osspkg/ui-transport", replacement: source("../../react-sdk/ui-transport/src/index.ts") },
    ],
  },
});
