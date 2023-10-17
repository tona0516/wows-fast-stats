module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:svelte/recommended",
    "prettier",
  ],
  overrides: [
    {
      env: {
        node: true,
      },
      files: [".eslintrc.{js,cjs}"],
      parserOptions: {
        sourceType: "script",
      },
    },
    {
      files: ["*.svelte"],
      parser: "svelte-eslint-parser",
      // Parse the `<script>` in `.svelte` as TypeScript by adding the following configuration.
      parserOptions: {
        parser: "@typescript-eslint/parser",
      },
    },
  ],
  parser: "@typescript-eslint/parser",
  parserOptions: {
    ecmaVersion: "latest",
    sourceType: "module",
  },
  plugins: ["@typescript-eslint"],
  rules: {
    "@typescript-eslint/ban-ts-comment": [
      "error",
      {
        // @ts-ignoreと@ts-nocheckのチェックを無効にする
        "ts-ignore": false,
        "ts-nocheck": false,
      },
    ],
    "@typescript-eslint/no-unused-vars": [
      "error",
      {
        // アンダースコアで始まる変数を無視する
        argsIgnorePattern: "^_",
        varsIgnorePattern: "^_",
        caughtErrorsIgnorePattern: "^_",
      },
    ],
    "@typescript-eslint/no-namespace": "off"
  },
};
