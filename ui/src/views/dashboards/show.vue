<template>
    <v-content>
        <v-modal-prompt></v-modal-prompt>
        <v-modal-time-range></v-modal-time-range>

        <v-toolbar clip="content">
            <v-button
                icon="sync"
                :icon-badge="!erred && options.refresh ? formatValue(options.refresh, {type: 'duration'}, 0) : null"
                :disabled="loading || erred"
                v-tooltip="$t(`labels.${params.type}.refresh`)"
            >
                <template slot="dropdown">
                    <template v-if="options.refresh && modifiers.alt">
                        <v-button disabled icon="">
                            {{ $t("labels.refresh.next", [formatValue(refresh.countdown, {type: "duration"}, 0)]) }}
                        </v-button>

                        <v-divider></v-divider>
                    </template>

                    <v-button icon="" :shortcut="['r', $t(`labels.${params.type}.refresh`)]" @click="refreshDashboard">
                        {{ $t(`labels.${params.type}.refresh`) }}
                    </v-button>

                    <v-divider></v-divider>

                    <v-columns :count="2">
                        <v-button
                            :icon="value === options.refresh ? 'check' : ''"
                            :key="index"
                            @click="setRefreshInterval(value)"
                            v-for="(value, index) in intervals"
                        >
                            {{ formatValue(value, {type: "duration"}) }}
                        </v-button>
                    </v-columns>

                    <v-divider></v-divider>

                    <v-button
                        icon="clock"
                        :shortcut="['shift+r', $t('labels.refresh.setInterval')]"
                        @click="setRefreshInterval()"
                    >
                        {{ $t("labels.custom") }}
                    </v-button>
                </template>
            </v-button>

            <v-button icon="calendar-alt" :disabled="loading || erred" v-tooltip="$t('labels.timeRange.set')">
                <v-flex class="timerange" v-if="timeRangeSynced">
                    <span>{{ $t("labels.timeRange.from") }}</span>
                    {{ absoluteRange ? formatDate(options.timeRange.from, undefined, false) : options.timeRange.from }}
                    <span>{{ $t("labels.timeRange.to") }}</span>
                    {{ absoluteRange ? formatDate(options.timeRange.to, undefined, false) : options.timeRange.to }}
                    <span class="utc" v-if="timezoneUTC">UTC</span>
                </v-flex>

                <template v-else>
                    {{ $t("labels.timeRange.multiple") }}
                </template>

                <template slot="dropdown">
                    <v-button icon="history" :disabled="!canResetTimeRange" @click="resetTimeRange">
                        {{ $t("labels.timeRange.reset") }}
                    </v-button>

                    <v-divider></v-divider>

                    <v-columns :count="2">
                        <v-button
                            :icon="
                                range.value === options.timeRange.from && options.timeRange.to === 'now' ? 'check' : ''
                            "
                            :key="index"
                            @click="setTimeRange({from: range.value, to: 'now'})"
                            v-for="(range, index) in ranges"
                        >
                            {{ $tc(`labels.timeRange.units.${range.unit}`, range.amount) }}
                        </v-button>
                    </v-columns>

                    <v-divider></v-divider>

                    <v-button
                        icon="calendar"
                        :shortcut="['alt+shift+r', $t('labels.timeRange.set')]"
                        @click="setTimeRange()"
                    >
                        {{ $t("labels.custom") }}
                    </v-button>

                    <v-divider></v-divider>

                    <v-label>{{ $t("labels.options") }}</v-label>

                    <v-checkbox toggle v-model="autoPropagate">{{ $t("labels.timeRange.autoPropagate") }}</v-checkbox>
                </template>
            </v-button>

            <v-spacer></v-spacer>

            <template v-if="params.type !== 'basket' && basket.length > 0">
                <v-button dropdown-anchor="bottom-right" icon="shopping-basket" :icon-badge="basket.length">
                    <template slot="dropdown">
                        <template v-if="params.type !== 'basket'">
                            <v-button icon="eye" :to="{name: 'basket-show'}">
                                {{ $t("labels.basket.preview") }}
                            </v-button>

                            <v-divider></v-divider>
                        </template>

                        <v-button icon="broom" @click="clearBasket">{{ $t("labels.basket.clear") }}</v-button>
                    </template>
                </v-button>

                <v-divider vertical></v-divider>
            </template>

            <v-button
                class="icon"
                dropdown-anchor="bottom-right"
                icon="angle-down"
                v-tooltip.nowrap="$t('labels.moreActions')"
                v-if="dashboard"
            >
                <template slot="dropdown">
                    <template v-if="params.type === 'basket'">
                        <v-button
                            icon="save"
                            :disabled="basket.length === 0"
                            :shortcut="['s', $t('labels.dashboards.save')]"
                            @click="saveBasket()"
                        >
                            {{ $t("labels.basket.saveDashboard") }}
                        </v-button>

                        <v-divider></v-divider>

                        <v-button icon="broom" :disabled="basket.length === 0" @click="clearBasket">
                            {{ $t("labels.basket.clear") }}
                        </v-button>
                    </template>

                    <v-button
                        icon="pencil-alt"
                        :shortcut="['e', $t(`labels.${params.type}.edit`)]"
                        :to="{
                            name: `admin-${params.type}-edit`,
                            params: {id: dashboard.id},
                            hash: !dashboard.template ? '#layout' : undefined,
                        }"
                        v-else
                    >
                        {{ $t(`labels.${params.type}.edit`) }}
                    </v-button>
                </template>
            </v-button>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <v-message-error @retry="getDashboard" v-else-if="erred"></v-message-error>

        <template v-else>
            <v-message type="info" v-if="!dashboard.items || dashboard.items.length === 0">
                {{ $t(`messages.${this.params.type}.empty`) }}
            </v-message>

            <v-form v-if="dynamicVariables.length > 0">
                <v-select
                    :key="index"
                    :label="variable.name"
                    :options="dynamicOptions[variable.name] || []"
                    @input="emitDashboardLoaded"
                    v-model="options.data[variable.name]"
                    v-for="(variable, index) in dynamicVariables"
                >
                </v-select>
            </v-form>

            <v-grid
                readonly
                ref="grid"
                :highlight-index="highlightIndex"
                :layout="dashboard.layout"
                v-model="dashboard.items"
            >
                <template slot-scope="item">
                    <v-chart
                        controls
                        emit-cursor
                        tooltip
                        :legend="item.value.options.legend"
                        :range="Object.assign({}, options.timeRange)"
                        v-model="resolvedRefs[`chart|${item.value.options.id}`]"
                        v-if="item.value.type === 'chart'"
                    >
                        <template slot="more">
                            <v-button
                                icon="minus"
                                @click="removeBasketItem(item.index)"
                                v-if="params.type === 'basket'"
                            >
                                {{ $t("labels.basket.remove") }}
                            </v-button>

                            <v-button icon="shopping-basket" @click="addBasketItem(item.index)" v-else>
                                {{ $t("labels.basket.add") }}
                            </v-button>
                        </template>
                    </v-chart>
                </template>
            </v-grid>
        </template>
    </v-content>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import isEqual from "lodash/isEqual";
import {Component, Mixins, Watch} from "vue-property-decorator";
import {Dictionary} from "vue-router/types/router";

import {SelectOption} from "@/types/components";

import {dateFormatDisplay, dateFormatRFC3339, defaultTimeRange, ranges} from "@/src/components/chart/chart.vue";
import GridComponent from "@/src/components/grid/grid.vue";
import {ModalPromptParams} from "@/src/components/modal/prompt.vue";
import {conflictCustomValidity} from "@/src/helpers/api";
import {renderChart} from "@/src/helpers/chart";
import {mapReferences, resolveVariables} from "@/src/helpers/dashboard";
import {CustomMixins} from "@/src/mixins";
import {updateRouteQuery} from "@/src/router";
import {ModalTimeRangeParams} from "@/src/views/dashboards/components/modal/time-range.vue";

interface Options {
    data: Record<string, string>;
    timeRange: TimeRange;
    refresh: number;
}

const defaultOptions = {
    data: {},
    timeRange: defaultTimeRange,
    refresh: 0,
};

export const itemTypes: Array<string> = ["chart"];

const intervals: Array<number> = [
    5, // 5s
    10, // 10s
    30, // 30s
    60, // 1m
    300, // 5m
    900, // 15m
    1800, // 30m
    3600, // 1h
    10800, // 3h
    21600, // 6h
    43200, // 12h
    86400, // 1d
];

@Component
export default class Show extends Mixins<CustomMixins>(CustomMixins) {
    public dashboard: Dashboard | null = null;

    public dashboardRefs: Record<string, unknown> = {};

    public dynamicData: Record<string, Array<string>> = {};

    public highlightIndex: number | null = null;

    public intervals: Array<number> = intervals;

    public loading = true;

    public options: Options = cloneDeep(defaultOptions);

    public ranges = ranges;

    public timeRangeSynced = true;

    public refresh: {
        countdown: number | null;
        interval: number | null;
        suspended: boolean;
    } = {
        countdown: null,
        interval: null,
        suspended: false,
    };

    public mounted(): void {
        document.addEventListener("visibilitychange", this.onVisibilityChange);
        window.addEventListener("hashchange", this.onHashChange);
        this.$root.$on("item-loaded", this.onItemLoaded);
        this.$root.$on("item-timerange", this.onItemTimeRange);
    }

    public beforeDestroy(): void {
        document.removeEventListener("visibilitychange", this.onVisibilityChange);
        window.removeEventListener("hashchange", this.onHashChange);
        this.$root.$off("item-loaded", this.onItemLoaded);
        this.$root.$off("item-timerange", this.onItemTimeRange);
    }

    public get absoluteRange(): boolean {
        return (
            this.parseDate(this.options.timeRange.from).isValid() && this.parseDate(this.options.timeRange.to).isValid()
        );
    }

    public addBasketItem(index: number): void {
        if (!this.dashboard?.items) {
            return;
        }

        const basket: Array<DashboardItem> = this.basket;
        const item: DashboardItem = this.dashboard.items[index];

        basket.push({
            type: item.type,
            layout: {x: 0, y: basket.length, w: 1, h: 1},
            options: Object.assign({}, item.options),
        });

        this.$store.commit("basket", basket);
    }

    public get autoPropagate(): boolean {
        return this.$store.getters.autoPropagate;
    }

    public set autoPropagate(value: boolean) {
        // Auto-propagate has been enabled, thus trigger initial resync if
        // out-of-sync.
        if (value && !this.timeRangeSynced) {
            this.$root.$emit("item-timerange", null, this.options.timeRange);
        }

        this.$store.commit("autoPropagate", value);
    }

    public get basket(): Array<DashboardItem> {
        return this.$store.getters.basket;
    }

    public get canResetTimeRange(): boolean {
        return !isEqual(this.options.timeRange, defaultTimeRange) || !this.timeRangeSynced;
    }

    public clearBasket(): void {
        if (this.dashboard !== null) {
            this.$set(this.dashboard, "items", []);
            this.dashboardRefs = {};
        }

        this.$store.commit("basket", []);
    }

    public get dynamicOptions(): Record<string, Array<SelectOption>> {
        return Object.keys(this.dynamicData).reduce((options: Record<string, Array<SelectOption>>, name: string) => {
            options[name] = this.dynamicData[name].map(value => ({label: value, value}));
            return options;
        }, {});
    }

    public get dynamicVariables(): Array<TemplateVariable> {
        const names: Array<string> = [];

        let variables =
            this.dashboard?.options?.variables?.filter(variable => {
                names.push(variable.name);
                return variable.dynamic;
            }) ?? [];

        Object.keys(this.dashboardRefs).forEach(key => {
            if (key.startsWith("chart|")) {
                const vars = (this.dashboardRefs[key] as Chart).options?.variables?.filter(variable => {
                    const keep = !names.includes(variable.name);
                    if (keep) {
                        names.push(variable.name);
                    }
                    return keep;
                });
                if (vars) {
                    variables = variables.concat(vars);
                }
            }
        });

        return variables;
    }

    public emitDashboardLoaded(): void {
        // Emit loaded event to trigger sidebar update
        this.$parent.$emit("dashboard-loaded", this.dashboard, this.resolvedRefs);
    }

    public getDashboard(): void {
        switch (this.params.type) {
            case "basket": {
                const charts: Array<string> = [];
                const types: Array<DashboardItemType> = [];

                this.$http
                    .post(
                        "/api/v1/bulk",
                        this.basket.reduce((req: Array<BulkRequest>, item: DashboardItem) => {
                            switch (item.type) {
                                case "chart":
                                    if (item.options && !charts.includes(item.options.id as string)) {
                                        req.push({
                                            endpoint: `/charts/${item.options.id}/resolve`,
                                            method: "POST",
                                        });

                                        charts.push(item.options.id as string);
                                        types.push("chart");
                                    }

                                    break;
                            }

                            return req;
                        }, []),
                    )
                    .then(response => response.json())
                    .then(
                        async (response: APIResponse<Array<BulkResult>>) => {
                            if (response.data && response.data.filter(result => result.status >= 400).length > 0) {
                                this.$components.notify(this.$t("messages.error.bulk") as string, "error");
                            }

                            this.dashboard = {
                                id: "",
                                name: "basket",
                                layout: {
                                    columns: 1,
                                    rowHeight: 260,
                                    rows: this.basket.length,
                                },
                                items: this.basket,
                            };

                            if (response.data) {
                                this.dashboardRefs = mapReferences(
                                    response.data.map((result, index) => ({
                                        type: types[index],
                                        value: result.response.data,
                                    })),
                                );

                                if (this.dynamicVariables.length > 0) {
                                    this.dynamicData = await resolveVariables(this, this.dynamicVariables);

                                    Object.keys(this.dynamicData).forEach(label => {
                                        if (!this.options.data[label]) {
                                            this.$set(this.options.data, label, this.dynamicData[label]?.[0]);
                                        }
                                    });
                                }
                            }

                            this.$parent.$emit("dashboard-loaded", this.dashboard, this.resolvedRefs);
                            this.loading = false;
                        },
                        this.handleError(() => {
                            this.$parent.$emit("dashboard-loaded", null);
                            this.loading = false;
                        }),
                    );

                break;
            }

            case "charts": {
                this.$http
                    .post(`/api/v1/charts/${this.params.id}/resolve`)
                    .then(response => response.json())
                    .then(
                        async (response: APIResponse<Chart>) => {
                            const chart = response.data as Chart;

                            this.dashboard = {
                                id: chart.id,
                                name: chart.name,
                                options: {
                                    variables: chart.options?.variables,
                                },
                                layout: {
                                    columns: 1,
                                    rowHeight: 260,
                                    rows: 1,
                                },
                                items: [
                                    {
                                        type: "chart",
                                        layout: {
                                            x: 0,
                                            y: 0,
                                            w: 1,
                                            h: 1,
                                        },
                                        options: {
                                            id: chart.id,
                                        },
                                    },
                                ],
                            };

                            if (this.dashboard?.options?.variables) {
                                this.dynamicData = await resolveVariables(this, this.dashboard.options.variables);

                                Object.keys(this.dynamicData).forEach(label => {
                                    if (!this.options.data[label]) {
                                        this.$set(this.options.data, label, this.dynamicData[label]?.[0]);
                                    }
                                });
                            }

                            this.dashboardRefs = {[`chart|${chart.id}`]: chart};

                            this.$parent.$emit("dashboard-loaded", this.dashboard, this.resolvedRefs);
                            this.loading = false;
                        },
                        this.handleError(() => {
                            this.$parent.$emit("dashboard-loaded", null);
                            this.loading = false;
                        }, true),
                    );
                break;
            }

            default: {
                this.$http
                    .post(`/api/v1/dashboards/${this.params.id}/resolve`)
                    .then(response => response.json())
                    .then(
                        async (response: APIResponse<Dashboard>) => {
                            this.dashboard = response.data as Dashboard;
                            this.dashboardRefs = mapReferences(this.dashboard.references as Array<Reference>);

                            delete this.dashboard.references;

                            if (this.dynamicVariables.length > 0) {
                                this.dynamicData = await resolveVariables(this, this.dynamicVariables);

                                Object.keys(this.dynamicData).forEach(label => {
                                    if (!this.options.data[label]) {
                                        this.$set(this.options.data, label, this.dynamicData[label]?.[0]);
                                    }
                                });
                            }

                            this.$parent.$emit("dashboard-loaded", this.dashboard, this.resolvedRefs);
                            this.loading = false;
                        },
                        this.handleError(() => {
                            this.$parent.$emit("dashboard-loaded", null);
                            this.loading = false;
                        }, true),
                    );
            }
        }
    }

    @Watch("options", {deep: true})
    public onOptions(to: Options): void {
        const q: Dictionary<string> = {};

        if (to.timeRange.from !== defaultTimeRange.from) {
            q.from = this.absoluteRange
                ? this.parseDate(to.timeRange.from).valueOf().toString()
                : (to.timeRange.from as string);
        }

        if (to.timeRange.to !== defaultTimeRange.to) {
            q.to = this.absoluteRange
                ? this.parseDate(to.timeRange.to).valueOf().toString()
                : (to.timeRange.to as string);
        }

        if (to.refresh > 0) {
            q.refresh = to.refresh.toString();
        }

        if (to.data) {
            Object.keys(to.data).forEach(label => {
                q[`var-${label}`] = encodeURIComponent(to.data[label]);
            });
        }

        updateRouteQuery(this.$route, q);

        this.updateRefresh();
    }

    @Watch("$route.path")
    public onRoutePath(): void {
        Object.assign(this, {
            dashboard: null,
            dashboardsRefs: {},
            loading: true,
        });
    }

    @Watch("$route.query", {immediate: true})
    public onRouteQuery(to: Dictionary<string>): void {
        const timeRange = cloneDeep(this.options.timeRange);
        let n: number;
        let update = false;

        if (to.from) {
            n = Number(to.from);
            timeRange.from = !isNaN(n) ? this.formatDate(n, dateFormatRFC3339, false) : to.from;
            update = true;
        }

        if (to.to) {
            n = Number(to.to);
            timeRange.to = !isNaN(n) ? this.formatDate(n, dateFormatRFC3339, false) : to.to;
            update = true;
        }

        this.options.data = Object.keys(to).reduce((data: Record<string, string>, key: string) => {
            if (key.startsWith("var-")) {
                data[key.substr(4)] = to[key];
            }
            return data;
        }, {});

        if (update) {
            this.setTimeRange(timeRange);
        }

        if (to.refresh) {
            this.setRefreshInterval(parseInt(to.refresh, 10) || 0);
        }

        this.getDashboard();
    }

    public refreshDashboard(): void {
        this.$root.$emit("item-refresh", null);
    }

    public removeBasketItem(index: number): void {
        if (!this.dashboard?.layout) {
            return;
        }

        const basket: Array<DashboardItem> = this.basket;
        basket.splice(index, 1);

        // Decrement following items Y position
        for (let i: number = index; i < basket.length; i++) {
            basket[i].layout.y--;
        }

        this.dashboard.layout.rows--;

        this.$store.commit("basket", basket);
    }

    public resetTimeRange(): void {
        this.setTimeRange(Object.assign({}, defaultTimeRange));
    }

    public get resolvedRefs(): Record<string, unknown> {
        return Object.keys(this.dashboardRefs).reduce((refs: Record<string, unknown>, key: string) => {
            if (key.startsWith("chart|")) {
                refs[key] = renderChart(this.dashboardRefs[key] as Chart, this.options.data);
            }
            return refs;
        }, {});
    }

    public saveBasket(name: string | null = null): void {
        if (name !== null) {
            this.$http.post("/api/v1/dashboards", Object.assign({}, this.dashboard, {name})).then(() => {
                this.clearBasket();
                this.$router.push({name: "dashboards-show", params: {id: name}});
            }, this.handleError());

            return;
        }

        this.$components.modal(
            "prompt",
            {
                button: {
                    label: this.$t("labels.dashboards.save") as string,
                    primary: true,
                },
                input: {
                    customValidity: conflictCustomValidity(this, "dashboards"),
                    required: true,
                    value: "",
                },
                message: this.$t("labels.dashboards.name") as string,
            } as ModalPromptParams,
            (value: string) => {
                if (value) {
                    this.saveBasket(value);
                }
            },
        );
    }

    public setRefreshInterval(value: number | null = null): void {
        if (value !== null) {
            this.options.refresh = value;
            return;
        }

        this.$components.modal(
            "prompt",
            {
                button: {
                    label: this.$t("labels.refresh.setInterval"),
                    primary: true,
                },
                input: {
                    help: this.$t("help.refresh.interval"),
                    type: "number",
                    value: this.options.refresh,
                },
                message: this.$t("labels.refresh.interval"),
            } as ModalPromptParams,
            (value: string | false) => {
                if (value !== false) {
                    this.setRefreshInterval(parseInt(value, 10) || 0);
                }
            },
        );
    }

    public setTimeRange(range: TimeRange | null = null): void {
        if (range !== null) {
            this.options.timeRange = range;
            this.$root.$emit("item-timerange", null, range);

            return;
        }

        this.$components.modal(
            "time-range",
            {
                timeRange: this.absoluteRange
                    ? {
                          from: this.parseDate(this.options.timeRange.from).format(dateFormatDisplay),
                          to: this.parseDate(this.options.timeRange.to).format(dateFormatDisplay),
                      }
                    : {
                          from: "",
                          to: "",
                      },
            } as ModalTimeRangeParams,
            (value: TimeRange) => {
                if (value !== false) {
                    this.setTimeRange({
                        from: this.parseDate(value.from, dateFormatDisplay).format(dateFormatRFC3339),
                        to: this.parseDate(value.to, dateFormatDisplay).format(dateFormatRFC3339),
                    });
                }
            },
        );
    }

    public get timezoneUTC(): boolean {
        return this.$store.getters.timezoneUTC;
    }

    private onHashChange(): void {
        if (window.location.hash) {
            const index = parseInt(window.location.hash.substr(5), 10);
            this.highlightIndex = !isNaN(index) ? index : null;

            this.$nextTick(() => {
                (this.$refs.grid as GridComponent).scrollTo(window.location.hash.substr(1));
            });
        } else {
            this.highlightIndex = null;
        }
    }

    private onItemLoaded(range: TimeRange): void {
        if (this.$store.getters.timeRange === null) {
            // We only keep the first received time range from items as a
            // reference.
            this.$store.commit("timeRange", range);
        }
    }

    private onItemTimeRange(target: Element | null, range: TimeRange): void {
        this.timeRangeSynced = target === null;
        if (this.timeRangeSynced) {
            this.options.timeRange = range;
        }
    }

    private onVisibilityChange(): void {
        this.refresh.suspended = document.visibilityState !== "visible";

        if (!this.refresh.suspended && this.options.refresh > 0) {
            this.refreshDashboard();
            this.updateRefresh();
        }
    }

    private updateRefresh(): void {
        // Trigger/Cancel refresh
        if (this.refresh.interval) {
            clearInterval(this.refresh.interval);
            this.refresh.interval = null;
        }

        if (this.options.refresh > 0) {
            this.refresh.countdown = this.options.refresh;

            this.refresh.interval = setInterval(() => {
                if (this.refresh.suspended) {
                    clearInterval(this.refresh.interval as number);
                    this.refresh.interval = null;
                    return;
                }

                (this.refresh.countdown as number)--;

                if (this.refresh.countdown === 0) {
                    this.refreshDashboard();
                    this.refresh.countdown = this.options.refresh;
                }
            }, 1000);
        } else {
            this.refresh.countdown = null;
        }
    }
}
</script>

<style lang="scss" scoped>
.v-content {
    padding: 1rem;

    .v-toolbar {
        .timerange {
            color: var(--color);

            span {
                color: var(--button-color);
                margin: 0 0.35rem;

                &:first-child {
                    margin-left: 0;
                }

                &:nth-child(2) {
                    text-transform: lowercase;
                }

                &.utc {
                    background-color: var(--gray);
                    border-radius: 0.2rem;
                    color: var(--toolbar-background-color);
                    display: inline-block;
                    font-size: 0.7rem;
                    height: 1.1rem;
                    line-height: 1.1rem;
                    margin: 0 0 0 0.65rem;
                    padding: 0 0.25rem;
                }
            }
        }
    }

    .v-form {
        margin-bottom: 1rem;

        .v-select {
            background-color: var(--toolbar-background);
            border-color: transparent;
            width: auto;

            & + .v-select {
                margin-left: 1rem;
            }
        }
    }

    ::v-deep .v-grid-item {
        background-color: #272727;

        &:hover {
            border-color: var(--toolbar-background);
        }

        .v-grid-anchor {
            display: block;
            height: 0;
            transform: translateY(calc(-2 * var(--toolbar-size) - 1rem));
            width: 0;
        }
    }
}
</style>
