<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-content>
        <teleport to="body">
            <v-modal-chart-preview></v-modal-chart-preview>
        </teleport>

        <v-toolbar clip="content">
            <v-button
                icon="eye"
                :disabled="selection.length === 0"
                :icon-badge="selection.length"
                @click="previewChart(selection)"
            >
                {{ i18n.t("labels.charts.preview") }}
            </v-button>

            <v-divider vertical></v-divider>

            <v-button icon="sync-alt" @click="getMetrics" v-shortcut="{keys: 'r', help: i18n.t('labels.refresh.list')}">
                {{ i18n.t("labels.refresh._") }}
            </v-button>

            <v-spacer></v-spacer>

            <v-input
                icon="filter"
                type="search"
                :placeholder="i18n.t('labels.metrics.filter')"
                @clear="applyFilter"
                @focusout="applyFilter"
                @keypress.enter="applyFilter"
                v-model:value="filter"
                v-shortcut="{keys: 'f', help: i18n.t('labels.metrics.filter')}"
            ></v-input>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <v-message-error @retry="getMetrics" v-else-if="erred"></v-message-error>

        <template v-else>
            <h1>
                {{ i18n.t("labels.metrics._", 2) }}
                <span class="count" v-if="total">{{ total }}</span>
            </h1>

            <v-message class="selection" icon="clipboard-list" type="info" v-if="selection.length > 0">
                {{ i18n.t("messages.metrics.selected", [selection.length], selection.length) }}

                <v-button
                    icon="times-circle"
                    @click="clearSelection"
                    v-tooltip="i18n.t('labels.clearSelection')"
                ></v-button>
            </v-message>

            <v-message type="info" v-if="metrics.length === 0">
                {{ i18n.t("messages.metrics.none") }}
            </v-message>

            <template v-else>
                <v-table ref="table" selectable v-model:selection="selection" v-model:value="metrics">
                    <template v-slot:header>
                        <v-table-cell>
                            {{ i18n.t("labels.name._") }}
                        </v-table-cell>

                        <v-table-cell grow>
                            {{ i18n.t("labels.labels", 2) }}
                        </v-table-cell>

                        <v-table-cell></v-table-cell>
                    </template>

                    <template v-slot="metric">
                        <v-table-cell>
                            <span class="link" @click="setMatch(null, metric.value.name())">
                                {{ metric.value.name() }}
                            </span>
                        </v-table-cell>

                        <v-table-cell grow>
                            <v-labels selectable @select="setMatch" :labels="metric.value.entries(false)"></v-labels>
                        </v-table-cell>

                        <v-table-cell>
                            <v-button
                                class="reveal icon"
                                icon="far/copy"
                                @click="clipboardCopy(metric.value.toString())"
                                v-tooltip="i18n.t('labels.clipboard.copy')"
                            ></v-button>
                        </v-table-cell>
                    </template>
                </v-table>

                <v-paging :page-size="limit" :total="total" v-model:page="options.page"></v-paging>
            </template>
        </template>
    </v-content>
</template>

<script lang="ts">
import {onBeforeMount, ref, watch} from "vue";
import {useI18n} from "vue-i18n";
import {useRouter} from "vue-router";
import {useStore} from "vuex";

import api from "@/lib/api";
import common from "@/common";
import {useUI} from "@/components/ui";
import {State} from "@/store";

import ModalChartPreviewComponent, {ModalChartPreviewParams} from "./modal/chart-preview.vue";

interface Options {
    filter: string;
    page: number;
}

const defaultOptions: Options = {
    filter: "",
    page: 1,
};

const limit = 20;

export default {
    components: {
        "v-modal-chart-preview": ModalChartPreviewComponent,
    },
    setup(): Record<string, unknown> {
        const i18n = useI18n();
        const router = useRouter();
        const store = useStore<State>();
        const ui = useUI();

        const {erred, loading, onFetchRejected} = common;

        const filter = ref("");
        const metrics = ref<Array<Labels>>([]);
        const options = ref(Object.assign({}, defaultOptions));
        const selection = ref<Array<Labels>>([]);
        const table = ref<HTMLTableElement | null>(null);
        const total = ref(0);

        const applyFilter = (): void => {
            options.value.filter = filter.value;
        };

        const clearSelection = (): void => {
            selection.value = [];
        };

        const getMetrics = (): void => {
            const params: ListParams = {
                limit,
                offset: (options.value.page - 1) * limit,
            };

            if (options.value.filter) {
                params.match =
                    options.value.filter.indexOf("{") !== -1
                        ? options.value.filter
                        : `{__name__=~${JSON.stringify(".*" + options.value.filter + ".*")}}`;
            }

            store.commit("loading", true);

            api.metrics(params)
                .then(response => {
                    if (response.data === undefined) {
                        return Promise.reject("cannot get metrics");
                    }

                    const pagesCount: number | undefined = response.total
                        ? Math.ceil(response.total / limit)
                        : undefined;

                    // Switch back to first/last page if current empty
                    if (!response.data?.length && options.value.page > 1) {
                        options.value.page = pagesCount !== undefined ? pagesCount : 1;
                        return;
                    }

                    metrics.value = response.data;
                    total.value = response.total ?? 0;
                }, onFetchRejected)
                .finally(() => {
                    store.commit("loading", false);
                });
        };

        const previewChart = (): void => {
            ui.modal("chart-preview", {
                exprs: selection.value.map(labels => labels.toString()),
            } as ModalChartPreviewParams);
        };

        const setMatch = (key: string | null, value: string): void => {
            filter.value = key === null ? value : `{${key}=${JSON.stringify(value)}}`;
            applyFilter();
        };

        const clipboardCopy = (value: string): void => {
            navigator.clipboard.writeText(value).then(() => ui.notify(i18n.t("messages.copied"), "success"));
        };

        onBeforeMount(() => {
            const query = router.currentRoute.value.query as Record<string, string>;

            if (query.filter) {
                options.value.filter = query.filter;
                filter.value = query.filter;
            }
            if (query.page) {
                options.value.page = parseInt(query.page, 10) ?? 1;
            }

            watch(
                options,
                (to: Options): void => {
                    const query: Record<string, string> = {};

                    if (to.filter !== "") {
                        query.filter = to.filter;
                    }
                    if (to.page !== 1) {
                        query.page = to.page.toString();
                    }

                    router.replace({query});

                    getMetrics();
                },
                {deep: true, immediate: true},
            );
        });

        return {
            applyFilter,
            clearSelection,
            clipboardCopy,
            erred,
            filter,
            getMetrics,
            i18n,
            limit,
            loading,
            metrics,
            options,
            previewChart,
            selection,
            setMatch,
            table,
            total,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../../mixins";

.v-content {
    @include content;

    ::v-deep(.v-table-row.selected .v-labels-entry) {
        background-color: var(--dark-accent);
    }
}
</style>
