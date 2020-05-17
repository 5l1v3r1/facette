<template>
    <v-content :class="{compact: !erred && section === 'layout'}">
        <v-modal-chart-series></v-modal-chart-series>
        <v-modal-confirm></v-modal-confirm>
        <v-modal-template-variable :allow-fixed="link"></v-modal-template-variable>

        <v-toolbar clip="content">
            <v-button
                icon="save"
                :disabled="erred || saving || !validity"
                @click="save(true)"
                v-if="prevRoute.name !== 'charts-show' && !template && modifiers.alt"
            >
                {{ $t("labels.saveAndGo") }}
            </v-button>

            <v-button icon="save" :disabled="erred || saving || !validity" @click="save(false)" v-else>
                {{ $t(`labels.${template ? "templates" : "charts"}.save`) }}
            </v-button>

            <v-button icon="trash" @click="deleteChart()" v-if="!erred && edit && modifiers.alt">
                {{ $t("labels.delete") }}
            </v-button>

            <v-button :disabled="erred" :to="prevRoute || {name: 'admin-charts-list'}" v-else>
                {{ $t("labels.cancel") }}
            </v-button>

            <v-divider vertical></v-divider>

            <v-button icon="undo" :disabled="erred || loading || !guarded" @click="reset()">
                {{ $t("labels.reset") }}
            </v-button>

            <template v-if="chart && edit">
                <v-spacer></v-spacer>

                <v-label class="modified" v-if="chart.modified">
                    {{ $t("messages.lastModified", [formatDate(chart.modified, $t("date.long"))]) }}
                </v-label>
            </template>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <v-message-error @retry="reset(true)" v-else-if="erred"></v-message-error>

        <template v-else>
            <div class="preview">
                <v-form v-if="dynamicVariables.length > 0">
                    <v-select
                        :key="index"
                        :label="variable.name"
                        :options="dynamicOptions[variable.name] || []"
                        v-model="resolveData[variable.name]"
                        v-for="(variable, index) in dynamicVariables"
                    >
                    </v-select>
                </v-form>

                <v-chart tooltip v-model="chartModel"></v-chart>
            </div>

            <h1 v-if="section === 'general'">{{ $t("labels.general") }}</h1>

            <v-form class="third" @validity="onValidity" v-show="section === 'general'">
                <v-label>{{ $t("labels.name") }}</v-label>
                <v-input
                    ref="name"
                    required
                    :delay="200"
                    :pattern="namePattern"
                    v-autofocus.select
                    v-model="chart.name"
                ></v-input>

                <template v-if="link">
                    <v-label>{{ $tc("labels.templates._", 1) }}</v-label>
                    <v-flex>
                        <v-select
                            :options="templates"
                            :placeholder="$t('labels.templates.select')"
                            v-model="chart.link"
                        >
                            <template slot="dropdown-placeholder" v-if="templates.length === 0">
                                <v-label>{{ $t("messages.templates.none") }}</v-label>
                            </template>
                        </v-select>

                        <v-button
                            icon="pencil-alt"
                            :to="{name: 'admin-charts-edit', params: {id: chart.link}}"
                            :style="{visibility: chart.link ? 'visible' : 'hidden'}"
                        >
                            {{ $t("labels.templates.edit") }}
                        </v-button>
                    </v-flex>
                </template>

                <template v-else>
                    <v-label>{{ $t("labels.title") }}</v-label>
                    <v-input :delay="200" :help="$t('help.charts.title')" v-model="chart.options.title"></v-input>

                    <v-flex class="half" direction="column">
                        <v-label>{{ $t("labels.charts.type._") }}</v-label>
                        <v-select
                            :options="types"
                            :placeholder="$t('labels.charts.type.select')"
                            v-model="chart.options.type"
                        ></v-select>
                    </v-flex>
                </template>
            </v-form>

            <template v-if="link">
                <template v-if="section === 'variables' && variables.length > 0">
                    <h1>{{ $t("labels.variables._") }}</h1>

                    <v-form-template-variables
                        :defined="chart.options.variables"
                        :parsed="variables"
                    ></v-form-template-variables>
                </template>
            </template>

            <template v-else>
                <template v-if="section === 'series'">
                    <h1>{{ $tc("labels.series._", 2) }}</h1>

                    <v-form>
                        <v-message type="info" v-if="chart.series.length === 0">
                            {{ $t("messages.series.none") }}
                        </v-message>

                        <v-table class="fixed" draggable v-model="chart.series" v-else>
                            <template slot="header">
                                <v-table-cell>{{ $tc("labels.series._", 2) }}</v-table-cell>
                                <v-table-cell class="more"></v-table-cell>
                            </template>

                            <template slot-scope="series">
                                <v-table-cell grow>
                                    <span class="color" :style="`background-color: ${getColor(series.index)}`"></span>
                                    <span class="monospace truncate" v-tooltip="series.value.expr">
                                        {{ series.value.expr }}
                                    </span>
                                </v-table-cell>

                                <v-table-cell class="more">
                                    <v-button
                                        class="reveal"
                                        icon="pencil-alt"
                                        @click="editSeries(series.index)"
                                        v-tooltip="$t('labels.series.edit')"
                                    ></v-button>
                                    <v-button
                                        class="reveal"
                                        icon="times"
                                        @click="removeSeries(series.index)"
                                        v-tooltip="$t('labels.series.remove')"
                                    ></v-button>
                                </v-table-cell>
                            </template>
                        </v-table>

                        <v-toolbar>
                            <v-button icon="plus" @click="addSeries">{{ $t("labels.series.add") }}</v-button>
                        </v-toolbar>
                    </v-form>
                </template>

                <template v-else-if="section === 'axes'">
                    <h1>{{ $t("labels.charts.axes._") }}</h1>

                    <v-form>
                        <v-flex class="columns">
                            <v-flex direction="column">
                                <h2>{{ $t("labels.charts.axes.yLeft") }}</h2>
                                <v-form-chart-yaxis :axis="chart.options.axes.y.left"></v-form-chart-yaxis>
                            </v-flex>

                            <v-flex direction="column">
                                <h2>{{ $t("labels.charts.axes.yRight") }}</h2>
                                <v-form-chart-yaxis :axis="chart.options.axes.y.right"></v-form-chart-yaxis>
                            </v-flex>

                            <v-flex direction="column">
                                <h2>{{ $t("labels.charts.axes.x") }}</h2>
                                <v-checkbox toggle v-model="chart.options.axes.x.show">
                                    {{ $t("labels.show") }}
                                </v-checkbox>
                            </v-flex>
                        </v-flex>
                    </v-form>
                </template>
            </template>
        </template>
    </v-content>
</template>

<script lang="ts">
import {defaultConfig as BoulaDefault} from "@facette/boula/src/config";
import cloneDeep from "lodash/cloneDeep";
import merge from "lodash/merge";
import {Component, Mixins, Watch} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

import {hash} from "@/src/helpers/hash";
import {beforeRoute} from "@/src/helpers/route";
import {resolveOption} from "@/src/helpers/select";
import {parseVariables, renderTemplate} from "@/src/helpers/template";
import {CustomMixins} from "@/src/mixins";

import {namePattern} from "..";

type DynamicData = Record<string, {hash: string; entries: Array<string>}>;

const defaultXAxis: ChartXAxis = {
    show: true,
};

const defaultYAxis: ChartYAxis = {
    show: true,
    constants: [],
    label: undefined,
    max: undefined,
    min: undefined,
    stack: undefined,
    unit: {
        type: undefined,
        base: undefined,
    },
};

const defaultChart: Chart = {
    id: "",
    name: "",
    options: {
        axes: {
            x: cloneDeep(defaultXAxis),
            y: {
                center: false,
                left: cloneDeep(defaultYAxis),
                right: merge({}, defaultYAxis, {show: false}),
            },
        },
        title: "",
        type: "area",
        variables: [],
    },
    series: [],
};

const defaultChartLinked: Chart = {
    id: "",
    name: "",
    options: {
        variables: [],
    },
};

const defaultSection = "general";

const types: Array<SelectOption> = [
    {label: "labels.charts.type.area", value: "area"},
    {label: "labels.charts.type.bar", value: "bar"},
    {label: "labels.charts.type.line", value: "line"},
];

function parseChartVariables(chart: Chart): Array<TemplateVariable> {
    let data = "";

    if (chart.options?.axes?.y?.left?.label) {
        data += `\xff${chart.options.axes.y.left.label}`;
    }

    if (chart.options?.axes?.y?.right?.label) {
        data += `\xff${chart.options.axes.y.right.label}`;
    }

    if (chart.options?.title) {
        data += `\xff${chart.options.title}`;
    }

    if (chart.series) {
        chart.series.forEach(series => {
            data += `\xff${series.expr}`;
        });
    }

    return parseVariables(data).map(name => ({name, dynamic: false}));
}

@Component({
    beforeRouteLeave: beforeRoute,
    beforeRouteUpdate: beforeRoute,
})
export default class Edit extends Mixins<CustomMixins>(CustomMixins) {
    public chart: Chart | null = null;

    public chartModel: Chart | null = null;

    public dynamicData: DynamicData = {};

    public dynamicOptions: Record<string, Array<SelectOption>> = {};

    public loading = true;

    public linked: Chart | null = null;

    public namePattern = namePattern;

    public resolveData: Record<string, string> = {};

    public saving = false;

    public section: string = defaultSection;

    public template = false;

    public templates: Array<SelectOption> = [];

    public validity = false;

    public variables: Array<TemplateVariable> = [];

    private unwatchChart: (() => void) | null = null;

    public mounted(): void {
        this.reset(true);
    }

    public beforeDestroy(): void {
        if (this.unwatchChart) {
            this.unwatchChart();
        }
    }

    public addSeries(): void {
        if (this.chart === null) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        }

        const series = this.chart.series as Array<ChartSeries>;

        this.$components.modal("chart-series", {series: {}}, (value: ChartSeries) => {
            if (value) {
                series.push(value);
            }
        });
    }

    public deleteChart(apply = false): void {
        if (this.chart === null) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        }

        if (apply) {
            this.$http.delete(`/api/v1/charts/${this.chart.id}`).then(() => {
                this.unguard();
                this.$components.notify(this.$tc("messages.charts.deleted", 1) as string, "success");
                this.$router.push(
                    this.prevRoute?.name === "charts-show"
                        ? {name: "dashboards-home"}
                        : {name: "admin-charts-list", query: this.template ? {kind: "template"} : {}},
                );
            }, this.handleError());

            return;
        }

        this.$components.modal(
            "confirm",
            {
                button: {
                    icon: "trash",
                    label: this.$tc(`labels.${this.params.type}.delete`, 1),
                    danger: true,
                },
                message: this.$tc(`messages.${this.params.type}.delete`, 1, this.chart),
            },
            (value: boolean) => {
                if (value) {
                    this.deleteChart(true);
                }
            },
        );
    }

    public get dynamicVariables(): Array<TemplateVariable> {
        return this.chart?.options?.variables?.filter(variable => variable.dynamic) ?? [];
    }

    public get edit(): boolean {
        return this.params.id !== "link" && this.params.id !== "new";
    }

    public editSeries(index: number): void {
        if (this.chart === null || !this.chart.series) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        }

        const series = this.chart.series as Array<ChartSeries>;

        this.$components.modal("chart-series", {series: cloneDeep(series[index])}, (value: ChartSeries) => {
            if (value) {
                series.splice(index, 1, value);
            }
        });
    }

    public getColor(index: number): string {
        const colors = BoulaDefault.colors as Array<string>;
        return this.chart?.series?.[index].options?.color || colors[index % colors.length];
    }

    public get link(): boolean {
        return this.params.id === "link" || Boolean(this.chart?.link);
    }

    @Watch("chart.options.variables", {deep: true})
    public async onVariables(to: Array<TemplateVariable>): Promise<void> {
        if (this.linked === null) {
            return;
        }

        const req: Array<BulkRequest> = [];
        const variables: Array<string> = [];

        const data = to.reduce((data: Record<string, string>, variable: TemplateVariable) => {
            if (variable.dynamic) {
                const hashValue = hash(variable);

                if (this.dynamicData[variable.name]?.hash === hashValue) {
                    data[variable.name] = this.dynamicData[variable.name].entries[0];
                } else {
                    this.dynamicData[variable.name] = {hash: hashValue, entries: []};

                    req.push({
                        endpoint: `/labels/${variable.label}/values`,
                        method: "GET",
                        params: variable.filter ? {match: variable.filter} : undefined,
                    });

                    variables.push(variable.name);
                }
            } else {
                data[variable.name] = variable.value as string;
            }

            return data;
        }, {});

        if (req.length > 0) {
            await this.$http
                .post("/api/v1/bulk", req)
                .then(response => response.json())
                .then((response: APIResponse<Array<BulkResult>>) => {
                    if (response.data && response.data.filter(result => result.status >= 400).length > 0) {
                        this.$components.notify(this.$t("messages.error.bulk") as string, "error");
                        return;
                    }

                    response.data?.forEach((result: BulkResult, index: number) => {
                        const values = result.response.data as Array<string>;

                        if (values.length > 0) {
                            data[variables[index]] = values[0];
                        }

                        this.dynamicData[variables[index]].entries = values;
                    });

                    this.updateDynamicOptions();
                });
        } else {
            this.updateDynamicOptions();
        }

        this.resolveData = data;
    }

    @Watch("params.id")
    public onParamsID(to: string, from: string): void {
        if (to !== from) {
            this.section = defaultSection;
            this.reset(true);
        }
    }

    @Watch("resolveData", {deep: true})
    public onResolveData(to: Record<string, string>): void {
        if (this.linked === null) {
            return;
        }

        const chart = cloneDeep(this.linked);

        if (chart.options?.title) {
            chart.options.title = renderTemplate(chart.options.title, to);
        }

        chart.series?.forEach(series => {
            series.expr = renderTemplate(series.expr, to);
        });

        this.chartModel = chart;
    }

    @Watch("$route.hash", {immediate: true})
    public onRouteHash(to: string): void {
        this.section = to ? to.substr(1) : defaultSection;
    }

    public onValidity(to: boolean): void {
        this.validity = to;
    }

    public removeSeries(index: number): void {
        if (this.chart === null || !this.chart.series) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        }

        this.chart.series.splice(index, 1);
    }

    public async reset(apply = false): Promise<void> {
        if (!apply) {
            this.$components.modal(
                "confirm",
                {
                    button: {
                        label: this.$t("labels.charts.reset"),
                        danger: true,
                    },
                    message: this.$t("messages.confirmLeave"),
                },
                (value: boolean) => {
                    if (value) {
                        this.reset(true);
                    }
                },
            );

            return;
        }

        if (this.unwatchChart) {
            this.unwatchChart();
        }

        this.validity = false;

        if (!this.edit) {
            let chart: Chart;

            if (this.link) {
                await this.getTemplates();
                chart = cloneDeep(defaultChartLinked);
                this.unwatchChart = this.$watch("chart", this.onChartLinked, {deep: true});
            } else {
                chart = cloneDeep(defaultChart);
                this.unwatchChart = this.$watch("chart", this.onChart, {deep: true});

                // Retrieve prefilled data from store if any
                const prefill: Array<string> | null = this.$store.getters.data;
                if (prefill) {
                    chart.series = prefill.map(expr => ({expr}));
                    this.$store.commit("data", null);
                }
            }

            this.chart = chart;
            this.loading = false;
            this.guardWatch("chart");

            return;
        }

        this.loading = true;

        this.$http
            .get(`/api/v1/charts/${this.params.id}`)
            .then(response => response.json())
            .then(
                async (response: APIResponse<Chart>) => {
                    if (response?.data?.link) {
                        await this.getTemplates();
                        this.unwatchChart = this.$watch("chart", this.onChartLinked, {deep: true});
                        this.chart = merge({}, defaultChartLinked, response.data);
                    } else {
                        this.unwatchChart = this.$watch("chart", this.onChart, {deep: true});
                        this.chart = merge({}, defaultChart, response.data);
                    }

                    this.loading = false;
                    this.guardWatch("chart");
                },
                this.handleError(() => {
                    this.loading = false;
                }, true),
            );
    }

    public save(go: boolean): void {
        if (this.chart === null) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        } else if (this.saving) {
            return;
        }
        this.saving = true;

        const chart: Chart = cloneDeep(this.chart);
        chart.template = this.template;

        this.$http({
            url: `/api/v1/charts${this.edit ? `/${this.params.id}` : ""}`,
            method: this.edit ? "PUT" : "POST",
            body: chart,
        }).then(
            () => {
                this.unguard();
                this.$components.notify(this.$t("messages.charts.saved") as string, "success");
                this.$router.push(
                    go || this.prevRoute?.name === "charts-show"
                        ? {name: "charts-show", params: {id: chart.name}}
                        : {name: "admin-charts-list", query: this.template ? {kind: "template"} : {}},
                );
                this.saving = false;
            },
            this.handleError(() => {
                this.saving = false;
            }),
        );
    }

    public get types(): Array<SelectOption> {
        return types.map((option: SelectOption) => resolveOption(this, option));
    }

    private emitUpdate(): void {
        this.$parent.$emit(
            "chart-updated",
            this.link,
            this.chart?.series?.length ?? null,
            this.variables.length ?? null,
        );
    }

    private getTemplates(): PromiseLike<void> {
        return this.$http
            .get("/api/v1/charts", {params: {kind: "template"}})
            .then(response => response.json())
            .then((response: APIResponse<Array<Chart>>) => {
                if (response.data) {
                    this.templates = response.data.map((chart: Chart) => ({label: chart.name, value: chart.id}));
                }
            });
    }

    private onChart(to: Chart): void {
        const variables = parseChartVariables(to);

        Object.assign(this, {
            chartModel: to,
            template: variables.length > 0,
            variables,
        });

        this.emitUpdate();
    }

    private onChartLinked(to: Chart, from: Chart | undefined): void {
        if (!to.link || to.link === from?.link) {
            this.emitUpdate();
            return;
        }

        this.$http
            .get(`/api/v1/charts/${to.link}`)
            .then(response => response.json())
            .then(
                (response: APIResponse<Chart>) => {
                    Object.assign(this, {
                        linked: response.data,
                        variables: parseChartVariables(response.data as Chart),
                    });

                    // FIXME: avoid having to trigger this manually (needed by
                    // initial template rendering)
                    this.onVariables(this.chart?.options?.variables ?? []);

                    this.emitUpdate();
                },
                this.handleError(() => {
                    this.emitUpdate();
                }),
            );
    }

    private updateDynamicOptions(): void {
        this.dynamicOptions = Object.keys(this.dynamicData).reduce(
            (options: Record<string, Array<SelectOption>>, name: string) => {
                options[name] = this.dynamicData[name].entries.map(value => ({label: value, value}));
                return options;
            },
            {},
        );
    }
}
</script>

<style lang="scss" scoped>
@import "../mixins";

.v-content {
    ::v-deep {
        @include content;
    }

    .preview {
        background-color: var(--background);
        border-bottom: 1px solid var(--sidebar-background);
        margin: -2rem -2rem 2rem;
        padding: 1rem;
        position: sticky;
        top: calc(var(--toolbar-size) * 2);
        z-index: 1;

        .v-form {
            margin-bottom: 1rem;

            .v-select {
                background-color: var(--toolbar-background);
                border-color: transparent;
                width: auto;
            }
        }

        .v-chart {
            height: 16rem;
            width: 100%;
        }
    }

    .color {
        border-radius: 0.1rem;
        display: inline-block;
        height: 0.6rem;
        margin-right: 0.5rem;
        min-width: 0.6rem;
        width: 0.6rem;
    }
}
</style>
