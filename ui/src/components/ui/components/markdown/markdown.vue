<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div class="v-markdown" v-html="html"></div>
</template>

<script lang="ts">
import DOMPurify from "dompurify";
import marked from "marked";
import {computed} from "vue";

DOMPurify.addHook("afterSanitizeAttributes", node => {
    if (node.tagName === "A") {
        try {
            if (new URL(node.getAttribute("href") as string, location.origin).origin !== location.origin) {
                node.setAttribute("target", "_blank");
            }
        } catch (err) {}
    }
});

export default {
    props: {
        content: {
            required: true,
            type: String,
        },
    },
    setup(props: Record<string, any>): Record<string, unknown> {
        const html = computed(() => {
            return DOMPurify.sanitize(marked(props.content.trim()));
        });

        return {
            html,
        };
    },
};
</script>
