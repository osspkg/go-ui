import { FlatCompat } from "@eslint/eslintrc";
import js from "@eslint/js";
import globals from "globals";
import { dirname } from "node:path";
import { fileURLToPath } from "node:url";
import tseslint from "typescript-eslint";

const compat = new FlatCompat({
  baseDirectory: dirname(fileURLToPath(import.meta.url)),
  recommendedConfig: js.configs.recommended,
});

export default tseslint.config(
  {
    ignores: ["**/dist/**", "**/node_modules/**", "**/coverage/**", "**/*.d.ts"],
  },
  ...compat.extends("airbnb", "airbnb/hooks"),
  ...tseslint.configs.recommended,
  {
    files: ["**/*.{ts,tsx}"],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
    settings: {
      react: { version: "19.0" },
    },
    rules: {
      "@typescript-eslint/no-unused-vars": ["error", { argsIgnorePattern: "^_" }],
      "react/jsx-filename-extension": ["error", { extensions: [".jsx", ".tsx"] }],
      "react/react-in-jsx-scope": "off",
      "react/prop-types": "off",
      "react/require-default-props": "off",
      "react/destructuring-assignment": "off",
      "react/jsx-props-no-spreading": "off",
      "react/no-array-index-key": "off",
      "import/prefer-default-export": "off",
      "no-nested-ternary": "off",
      "consistent-return": "off",
      "import/no-extraneous-dependencies": "off",
      "import/no-unresolved": "off",
      "import/extensions": "off",
      "no-use-before-define": "off",
      "@typescript-eslint/no-use-before-define": "off",
      "no-restricted-syntax": "off",
      "no-await-in-loop": "off",
      "no-continue": "off",
      "no-void": "off",
      "lines-between-class-members": "off",
      "class-methods-use-this": "off",
      "no-bitwise": "off",
      "no-plusplus": "off",
      "no-shadow": "off",
      "no-param-reassign": "off",
    },
  },
  ...compat.extends("prettier"),
);
