<template>
    <div class="v-markdown" v-html="rendered"></div>
</template>

<script lang="ts">
import DOMPurify from "dompurify";
import marked from "marked";
import {Component, Vue} from "vue-property-decorator";

DOMPurify.addHook("afterSanitizeAttributes", node => {
    if (node.tagName === "A") {
        try {
            if (new URL(node.getAttribute("href") as string, location.origin).origin !== location.origin) {
                node.setAttribute("target", "_blank");
            }
        } catch (e) {
            /* noop */
        }
    }
});

@Component
export default class MarkdownComponent extends Vue {
    public get rendered(): string {
        let text: string = this.$slots.default?.[0].text?.trim() ?? "";
        if (text) {
            text = DOMPurify.sanitize(marked(text));
        }

        return text;
    }
}
</script>
