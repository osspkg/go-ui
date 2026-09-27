import { copyFile } from "node:fs/promises";

await copyFile(
  new URL("../src/styles.css.d.ts", import.meta.url),
  new URL("../dist/styles.css.d.ts", import.meta.url),
);
