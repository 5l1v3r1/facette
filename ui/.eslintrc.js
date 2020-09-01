module.exports = {
    root: true,
    env: {
        node: true,
    },
    extends: [
        "plugin:vue/vue3-essential",
        "eslint:recommended",
        "@vue/typescript/recommended",
        "@vue/prettier",
        "@vue/prettier/@typescript-eslint",
    ],
    parserOptions: {
        ecmaVersion: 2020,
    },
    rules: {
        "@typescript-eslint/array-type": ["error", {default: "generic"}],
        "@typescript-eslint/consistent-type-definitions": ["error", "interface"],
        "@typescript-eslint/no-explicit-any": "off",
        "eqeqeq": ["error", "always"],
        "lines-between-class-members": ["error", "always"],
        "no-console": process.env.NODE_ENV === "production" ? "error" : "off",
        "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
        "no-empty": ["error", {allowEmptyCatch: true}],
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
