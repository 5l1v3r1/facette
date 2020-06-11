<template>
    <div class="v-markdown" v-html="html"></div>
</template>

<script lang="ts">
import DOMPurify from "dompurify";
import marked from "marked";
import {Component, Prop, Vue} from "vue-property-decorator";

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
    @Prop({required: true, type: String})
    public source!: string;

    public get html(): string {
        return DOMPurify.sanitize(marked(this.source.trim()));
    }
}
</script>
