<template>
    <body class="v-app" :class="{[theme]: true}" :style="style">
        <slot></slot>

        <v-tooltip v-model="tooltip"></v-tooltip>
    </body>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

import {AppTheme, TooltipState} from "@/types/components";

import themes from "./themes";

@Component
export default class AppComponent extends Vue {
    @Prop({default: "dark", type: String})
    public theme!: string;

    public get style(): AppTheme {
        const theme = themes[this.theme];
        if (!theme) {
            return {};
        }

        return Object.keys(theme).reduce((style: AppTheme, key: string) => {
            if (key !== "name") {
                style[`--${key}`] = theme[key];
            }
            return style;
        }, {});
    }

    public get tooltip(): TooltipState | null {
        return this.$components.state.tooltip;
    }
}
</script>

<style lang="scss">
@font-face {
    font-display: swap;
    font-family: "Roboto";
    font-style: normal;
    font-weight: 400;
    src: local("Roboto Regular"), local("Roboto-Regular"),
        url("~typeface-roboto/files/roboto-latin-400.woff2") format("woff2"),
        url("~typeface-roboto/files/roboto-latin-400.woff") format("woff");
}

@font-face {
    font-display: swap;
    font-family: "Roboto";
    font-style: normal;
    font-weight: 500;
    src: local("Roboto Regular"), local("Roboto-Regular"),
        url("~typeface-roboto/files/roboto-latin-500.woff2") format("woff2"),
        url("~typeface-roboto/files/roboto-latin-500.woff") format("woff");
}

@font-face {
    font-display: swap;
    font-family: "Roboto Mono";
    font-style: normal;
    font-weight: 400;
    src: local("Roboto Mono"), local("Roboto-Mono"),
        url("~typeface-roboto-mono/files/roboto-mono-latin-400.woff2") format("woff2"),
        url("~typeface-roboto-mono/files/roboto-mono-latin-400.woff") format("woff");
}

html {
    box-sizing: border-box;
    font-family: "Roboto", sans-serif;
}

body {
    background-color: var(--background);
    color: var(--color);
    font-size: 0.85rem;
    line-height: 1.5;
    margin: 0;
}

*,
*::before,
*::after {
    box-sizing: inherit;
}

a {
    color: inherit;
}

::placeholder {
    opacity: var(--placeholder-opacity);
}
</style>
