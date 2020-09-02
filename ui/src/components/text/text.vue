<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div class="v-text">
        <v-toolbar v-if="controls">
            <v-button
                class="icon"
                dropdown-anchor="bottom-right"
                icon="angle-double-down"
                v-tooltip.nowrap="i18n.t('labels.moreActions')"
            >
                <template v-slot:dropdown>
                    <v-button dropdown-anchor="right" icon="file-download">
                        {{ i18n.t("labels.export._") }}
                        <template v-slot:dropdown>
                            <v-button @click="downloadExport()">
                                {{ i18n.t("labels.export.textMarkdown") }}
                            </v-button>
                        </template>
                    </v-button>

                    <template v-if="hasMore">
                        <v-divider></v-divider>

                        <slot name="more"></slot>
                    </template>
                </template>
            </v-button>
        </v-toolbar>

        <v-markdown :content="value.content"></v-markdown>
    </div>
</template>

<script lang="ts">
import slugify from "slugify";
import {SetupContext, onMounted, ref} from "vue";
import {useI18n} from "vue-i18n";

import {dateFormatFilename} from "@/components/chart";
import {parseDate} from "@/helpers/date";

export default {
    props: {
        controls: {
            default: false,
            type: Boolean,
        },
        value: {
            required: true,
            validator: (prop: unknown): boolean => typeof prop === "object" || prop === null,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const i18n = useI18n();

        const hasMore = ref(false);

        const downloadExport = (): void => {
            const el: HTMLAnchorElement = document.createElement("a");

            Object.assign(el, {
                download: `${slugify(props.value.title)}_${parseDate().toFormat(dateFormatFilename)}.md`,
                href: `data:text/markdown,${encodeURIComponent(props.value.content)}`,
            });

            document.body.appendChild(el);
            el.click();
            document.body.removeChild(el);
        };

        onMounted(() => {
            hasMore.value = Boolean(ctx.slots.more);
        });

        return {
            downloadExport,
            hasMore,
            i18n,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../../views/mixins";

.v-text {
    height: 100%;
    overflow: auto;
    padding: 1rem;

    .v-toolbar {
        @include item-toolbar;

        position: absolute;
        right: 0;
        top: 0;
        visibility: hidden;
    }

    &:hover .v-toolbar {
        visibility: visible;
    }

    ::v-deep(img) {
        max-width: 100%;
    }
}
</style>
