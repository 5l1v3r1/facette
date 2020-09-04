<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div class="v-theme">
        <svg width="400" height="300" viewBox="0 0 400 300" xmlns="http://www.w3.org/2000/svg" :style="style">
            <rect width="400" height="300" fill="var(--background)" />
            <rect x="115" y="95" width="70" height="60" rx="5" fill="var(--input-background)" />
            <rect x="200" y="95" width="185" height="135" rx="5" fill="var(--input-background)" />
            <rect x="115" y="170" width="70" height="60" rx="5" fill="var(--input-background)" />
            <rect x="100" y="40" width="300" height="40" fill="var(--toolbar-background)" />
            <rect x="370" y="55" width="15" height="10" rx="5" fill="var(--toolbar-color)" />
            <rect x="145" y="55" width="60" height="10" rx="5" fill="var(--toolbar-color)" />
            <rect x="115" y="55" width="15" height="10" rx="5" fill="var(--toolbar-color)" />
            <rect y="40" width="100" height="260" fill="var(--sidebar-background)" />
            <rect x="20" y="190" width="60" height="10" rx="5" fill="var(--sidebar-color)" />
            <rect x="20" y="160" width="60" height="10" rx="5" fill="var(--sidebar-color)" />
            <rect x="20" y="130" width="60" height="10" rx="5" fill="var(--sidebar-color)" />
            <rect x="20" y="100" width="30" height="10" rx="5" fill="var(--light-gray)" />
            <rect y="40" width="100" height="40" fill="var(--sidebar-toolbar-background)" />
            <rect x="20" y="55" width="60" height="10" rx="5" fill="var(--sidebar-color)" />
            <rect width="400" height="40" fill="var(--accent)" />
            <rect x="140" y="10" width="120" height="20" rx="5" fill="var(--dark-accent)" />
            <rect x="20" y="15" width="40" height="10" rx="5" fill="white" />
        </svg>

        <div class="v-theme-name">{{ theme.name }}</div>
    </div>
</template>

<script lang="ts">
import {computed} from "vue";

import themes from "@/components/ui/themes";

export default {
    props: {
        name: {
            required: true,
            type: String,
        },
    },
    setup(props: Record<string, any>): Record<string, unknown> {
        const style = computed(() => {
            return Object.keys(theme.value).reduce((style: Record<string, string>, key: string) => {
                style[`--${key}`] = theme.value[key];
                return style;
            }, {});
        });

        const theme = computed(() => themes[props.name]);

        return {
            style,
            theme,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-theme {
    align-items: center;
    background-color: var(--input-background);
    display: flex;
    flex-direction: column;

    svg {
        border-radius: 0.2rem;
        overflow: hidden;
        pointer-events: none;
    }

    .v-theme-name {
        color: var(--color);
        margin-top: 0.5rem;
    }
}
</style>
