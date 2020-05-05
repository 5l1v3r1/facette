/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

module.exports = {
    root: true,
    env: {
        node: true,
    },
    extends: [
        "plugin:vue/essential",
        "eslint:recommended",
        "@vue/typescript/recommended",
        "@vue/prettier",
        "@vue/prettier/@typescript-eslint",
    ],
    overrides: [
        {
            files: ["*.js"],
            rules: {
                "@typescript-eslint/camelcase": "off",
                "@typescript-eslint/explicit-function-return-type": "off",
                "@typescript-eslint/no-var-requires": "off",
            },
        },
        {
            files: ["**/*.test.ts"],
            env: {
                mocha: true,
            },
        },
    ],
    parserOptions: {
        ecmaVersion: 2020,
    },
    rules: {
        "@typescript-eslint/array-type": ["error", {default: "generic"}],
        "@typescript-eslint/consistent-type-definitions": ["error", "interface"],
        "@typescript-eslint/explicit-function-return-type": ["error", {allowExpressions: true}],
        "eqeqeq": ["error", "always"],
        "no-console": process.env.NODE_ENV === "production" ? "error" : "off",
        "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
        "prettier/prettier": [
            process.env.NODE_ENV === "production" ? "error" : "warn",
            {
                arrowParens: "avoid",
                bracketSpacing: false,
                printWidth: 120,
                quoteProps: "consistent",
                tabWidth: 4,
                trailingComma: "all",
            },
        ],
        "sort-imports": ["error", {ignoreDeclarationSort: true}],
    },
};
