import { defineConfig } from "vitest/config";
import { fileURLToPath } from "node:url";

export default defineConfig({
  root: fileURLToPath(new URL(".", import.meta.url)),
  test: {
    include: ["react-sdk/*/src/**/*.test.ts", "react-sdk/*/src/**/*.test.tsx"],
    environment: "node",
  },
});
