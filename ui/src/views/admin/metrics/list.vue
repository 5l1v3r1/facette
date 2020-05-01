<template>
    <v-content>
        <v-modal-chart-preview></v-modal-chart-preview>

        <v-toolbar clip="content">
            <v-button
                icon="sync"
                :tooltip="$t('labels.refresh.list')"
                :shortcut="['r', $t('labels.refresh.list')]"
                @click="getMetrics"
            ></v-button>

            <v-divider vertical></v-divider>

            <v-button
                :disabled="loading || polling"
                :icon="polling ? 'circle-notch' : 'arrow-alt-circle-down'"
                :icon-spin="polling"
                @click="pollProviders"
            >
                {{ $t("labels.providers.pollAlt") }}
            </v-button>

            <template v-if="selection.length > 0">
                <v-divider vertical></v-divider>
                <v-button
                    icon="eye"
                    :icon-badge="selection.length"
                    @click="previewChart(Object.values(selection.items))"
                ></v-button>
            </template>

            <v-spacer></v-spacer>

            <v-input
                icon="filter"
                type="search"
                :placeholder="$t('labels.metrics.filter')"
                :shortcut="['f', $t('labels.metrics.filter')]"
                @clear="applyFilter"
                @focusout="applyFilter"
                @keypress.enter.native="applyFilter"
                v-model="tmpFilter"
            ></v-input>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <v-message-error @retry="getMetrics" v-else-if="erred"></v-message-error>

        <template v-else>
            <h1>
                {{ $tc("labels.metrics._", 2) }}
                <span class="count" v-if="total">{{ total }}</span>
            </h1>

            <v-message class="selection" icon="clipboard-check" type="info" v-if="selection.length > 0">
                {{ $tc("messages.metrics.selected", selection.length, [selection.length]) }}
                <v-button icon="times-circle" :tooltip="$t('labels.clearSelection')" @click="clearSelection"></v-button>
            </v-message>

            <v-message type="info" v-if="metrics.length === 0">{{ $t("messages.metrics.none") }}</v-message>

            <template v-else>
                <v-table ref="table" selectable @selection="onSelection" v-model="metrics">
                    <template slot="header">
                        <v-table-cell>
                            {{ $t("labels.name") }}
                        </v-table-cell>

                        <v-table-cell grow>
                            {{ $tc("labels.labels", 2) }}
                        </v-table-cell>

                        <v-table-cell></v-table-cell>
                    </template>

                    <template slot-scope="metric">
                        <v-table-cell>
                            <span class="link" @click="setMatch('__name__', metric.value.name)">
                                {{ metric.value.__name__ }}
                            </span>
                        </v-table-cell>

                        <v-table-cell grow>
                            <v-labels :handler="setMatch" :labels="metric.value"></v-labels>
                        </v-table-cell>

                        <v-table-cell>
                            <v-button class="reveal icon" icon="eye" @click="previewChart([metric.value])"></v-button>
                        </v-table-cell>
                    </template>
                </v-table>

                <v-paging :handler="switchPage" :page="options.page" :page-size="limit" :total="total"></v-paging>
            </template>
        </template>
    </v-content>
</template>

<script lang="ts">
import {Component, Mixins, Watch} from "vue-property-decorator";
import {Dictionary} from "vue-router/types/router";

import {hash} from "@/src/helpers/hash";
import {labelsToString, parseLabels} from "@/src/helpers/labels";
import {CustomMixins} from "@/src/mixins";
import {updateRouteQuery} from "@/src/router";
import TableComponent from "@/src/components/vue/table/table.vue";

interface Metric extends Labels {
    _hash: string;
}

interface Options {
    filter: string;
    page: number;
}

@Component
export default class List extends Mixins<CustomMixins>(CustomMixins) {
    public metrics: Array<Metric> = [];

    public limit = 20;

    public loading = true;

    public options: Options = {
        filter: "",
        page: 1,
    };

    public polling = false;

    public selection: {
        items: Record<string, Metric>;
        length: number;
    } = {
        items: {},
        length: 0,
    };

    public tmpFilter = "";

    public total = 0;

    public applyFilter(): void {
        if (this.tmpFilter !== this.options.filter) {
            this.options.filter = this.tmpFilter;
        }
    }

    public clearSelection(): void {
        this.selection = {
            items: {},
            length: 0,
        };

        const table = this.$refs.table as TableComponent;
        if (table) {
            table.select([]);
        }
    }

    @Watch("metrics")
    public onMetrics(to: Array<Metric>): void {
        if (this.selection.length !== 0) {
            this.$nextTick(() => {
                const table = this.$refs.table as TableComponent;
                if (table) {
                    table.select(to.filter(metric => this.selection.items[metric._hash]));
                }
            });
        }
    }

    @Watch("options", {deep: true})
    public onOptions(to: Options): void {
        const q: Dictionary<string> = {};

        if (to.filter !== "") {
            q.filter = to.filter;
        }
        if (to.page !== 1) {
            q.page = to.page.toString();
        }

        updateRouteQuery(this.$route, q);

        this.getMetrics();
    }

    @Watch("$route.query", {immediate: true})
    public onRouteQuery(to: Dictionary<string>, from: Dictionary<string> | undefined): void {
        if (from === undefined) {
            const options: Options = {
                filter: "",
                page: 1,
            };

            if (to.filter) {
                options.filter = to.filter;
                this.tmpFilter = to.filter;
            }
            if (to.page) {
                options.page = parseInt(to.page, 10) || 1;
            }

            Object.assign(this.options, options);
        }

        // Get metrics if no query parameter (i.e. initial get)
        if (Object.keys(to).length === 0) {
            this.getMetrics();
        }
    }

    public onSelection(to: Array<Metric>): void {
        const hashes: Array<string> = this.metrics.map(metric => metric._hash);

        const items = Object.values(this.selection.items)
            .filter(metric => !hashes.includes(metric._hash))
            .concat(to)
            .reduce((items: Record<string, Metric>, metric: Metric) => {
                items[metric._hash] = metric;
                return items;
            }, {});

        this.selection = {
            items,
            length: Object.keys(items).length,
        };
    }

    public pollProviders(): void {
        this.polling = true;

        this.$http.post("/api/v1/providers/poll").then(
            () => {
                this.polling = false;
                this.getMetrics();
            },
            this.handleError(() => {
                this.polling = false;
            }),
        );
    }

    public previewChart(metrics: Array<Metric>): void {
        this.$components.modal("chart-preview", {exprs: metrics.map((metric: Metric) => labelsToString(metric))});
    }

    public setMatch(key: string, value: string): void {
        this.tmpFilter = key === "__name__" ? value : `{${key}=${JSON.stringify(value)}}`;
        this.applyFilter();
    }

    public switchPage(page: number): void {
        this.options.page = page;
    }

    private getMetrics(): void {
        const params: ListParams = {
            limit: this.limit,
            offset: (this.options.page - 1) * this.limit || undefined,
        };

        if (this.options.filter) {
            params.match =
                this.options.filter.indexOf("{") !== -1
                    ? this.options.filter
                    : `{__name__=~${JSON.stringify(".*" + this.options.filter + ".*")}}`;
        }

        this.loading = true;

        this.$http
            .get("/api/v1/metrics", {params})
            .then(response => response.json())
            .then(
                (response: APIResponse<Array<string>>) => {
                    const pagesCount: number | undefined = response.total
                        ? Math.ceil(response.total / this.limit)
                        : undefined;

                    // Switch back to first/last page if current empty
                    if (response.data && response.data.length === 0 && this.options.page > 1) {
                        this.options.page = pagesCount !== undefined ? pagesCount : 1;
                        return;
                    }

                    if (response.data) {
                        this.metrics = response.data.map(expr => {
                            const metric: Metric = parseLabels(expr) as Metric;
                            metric._hash = hash(metric);
                            return metric;
                        });
                    }

                    this.total = response.total || 0;
                    this.loading = false;
                },
                this.handleError(() => {
                    this.loading = false;
                }),
            );
    }
}
</script>

<style lang="scss" scoped>
@import "../mixins";

.v-content {
    @include admin;

    .v-toolbar .v-input {
        width: 20rem;
    }

    ::v-deep .v-table-row.selected .v-labels-item {
        background-color: var(--table-row-selected-form);
    }
}
</style>
