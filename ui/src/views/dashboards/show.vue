<template>
    <v-content>
        <v-modal-prompt></v-modal-prompt>
        <v-modal-time-range></v-modal-time-range>

        <v-toolbar clip="content">
            <v-button
                icon="sync"
                :icon-badge="!erred && options.refresh ? `${options.refresh}s` : null"
                :disabled="loading || erred"
                :tooltip="$t(`labels.${params.type}.refresh`)"
            >
                <template slot="dropdown">
                    <template v-if="options.refresh && modifiers.alt">
                        <v-button disabled>
                            {{ $t("labels.refresh.next", [formatValue(refresh.countdown, "duration")]) }}
                        </v-button>
                        <v-divider></v-divider>
                    </template>

                    <v-button :shortcut="['r', $t(`labels.${params.type}.refresh`)]" @click="refreshDashboard">
                        {{ $t(`labels.${params.type}.refresh`) }}
                    </v-button>

                    <v-divider></v-divider>

                    <v-button :shortcut="['shift+r', $t('labels.refresh.setInterval')]" @click="setRefreshInterval()">
                        {{ $t("labels.refresh.setIntervalAlt") }}
                    </v-button>
                </template>
            </v-button>

            <v-button icon="clock" :disabled="loading || erred" :tooltip="$t('labels.timeRange.set')">
                <v-flex class="timerange">
                    <span>{{ $t("labels.timeRange.from") }}</span>
                    {{ absoluteRange ? formatDate(options.range.from, undefined, false) : options.range.from }}
                    <span>{{ $t("labels.timeRange.to") }}</span>
                    {{ absoluteRange ? formatDate(options.range.to, undefined, false) : options.range.to }}
                    <span class="utc" v-if="timezoneUTC">UTC</span>
                </v-flex>

                <template slot="dropdown">
                    <v-button icon="history" :disabled="defaultRange" @click="resetTimeRange">
                        {{ $t("labels.timeRange.reset") }}
                    </v-button>

                    <v-divider></v-divider>

                    <v-columns :count="2">
                        <v-button
                            :icon="range.value === options.range.from && options.range.to === 'now' ? 'check' : ''"
                            :key="index"
                            @click="setTimeRange({from: range.value, to: 'now'})"
                            v-for="(range, index) in ranges"
                        >
                            {{ $tc(`labels.timeRange.units.${range.unit}`, range.amount) }}
                        </v-button>
                    </v-columns>

                    <v-divider></v-divider>

                    <v-button
                        icon="calendar-alt"
                        :shortcut="['alt+shift+r', $t('labels.timeRange.set')]"
                        @click="setTimeRange()"
                    >
                        {{ $t("labels.timeRange.setCustom") }}
                    </v-button>

                    <v-divider></v-divider>

                    <v-label>{{ $t("labels.options") }}</v-label>

                    <v-checkbox toggle v-model="autoPropagate">{{ $t("labels.timeRange.autoPropagate") }}</v-checkbox>
                </template>
            </v-button>

            <v-spacer></v-spacer>

            <v-button
                class="icon"
                dropdown-anchor="bottom-right"
                icon="angle-down"
                :tooltip="$t('labels.moreActions')"
                v-if="dashboard"
            >
                <template slot="dropdown">
                    <v-button
                        icon="pencil-alt"
                        :shortcut="['e', $t(`labels.${params.type}.edit`)]"
                        :to="{
                            name: `admin-${params.type}-edit`,
                            params: {id: dashboard.id},
                            hash: !dashboard.template ? '#layout' : undefined,
                        }"
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

            <v-grid readonly :layout="dashboard.layout" v-model="dashboard.items">
                <template slot-scope="item">
                    <v-chart
                        controls
                        emit-cursor
                        tooltip
                        :legend="item.value.options.legend"
                        v-model="dashboardRefs[`chart|${item.value.options.id}`]"
                        v-if="item.value.type === 'chart'"
                    >
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

import {dateFormatRFC3339} from "@/src/components/chart/chart.vue";
import {CustomMixins} from "@/src/mixins";
import {updateRouteQuery} from "@/src/router";

interface Options {
    data: Record<string, string>;
    range: TimeRange;
    refresh: number;
}

interface Range {
    unit: string;
    amount: number;
    value: string;
}

const datetimeFormat = "YYYY-MM-DD HH:mm:ss";

const defaultOptions: Options = {
    data: {},
    range: {
        from: "-1h",
        to: "now",
    },
    refresh: 0,
};

export const itemTypes: Array<string> = ["chart"];

const ranges: Array<Range> = [
    {unit: "minutes", amount: 5, value: "-5m"},
    {unit: "minutes", amount: 15, value: "-15m"},
    {unit: "minutes", amount: 30, value: "-30m"},
    {unit: "hours", amount: 1, value: "-1h"},
    {unit: "hours", amount: 3, value: "-3h"},
    {unit: "hours", amount: 6, value: "-6h"},
    {unit: "hours", amount: 12, value: "-12h"},
    {unit: "days", amount: 1, value: "-1d"},
    {unit: "days", amount: 3, value: "-3d"},
    {unit: "days", amount: 7, value: "-7d"},
    {unit: "months", amount: 1, value: "-1M"},
    {unit: "months", amount: 3, value: "-3M"},
    {unit: "months", amount: 6, value: "-6M"},
    {unit: "years", amount: 1, value: "-1y"},
];

@Component
export default class Show extends Mixins<CustomMixins>(CustomMixins) {
    public dashboard: Dashboard | null = null;

    public dashboardRefs: Record<string, unknown> = {};

    public loading = true;

    public options: Options = cloneDeep(defaultOptions);

    public ranges: Array<Range> = ranges;

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
    }

    public beforeDestroy(): void {
        document.removeEventListener("visibilitychange", this.onVisibilityChange);
    }

    public get absoluteRange(): boolean {
        return this.parseDate(this.options.range.from).isValid() && this.parseDate(this.options.range.to).isValid();
    }

    public get autoPropagate(): boolean {
        return this.$store.getters.autoPropagate;
    }

    public set autoPropagate(value: boolean) {
        this.$store.commit("autoPropagate", value);
    }

    public get defaultRange(): boolean {
        return isEqual(this.options.range, defaultOptions.range);
    }

    @Watch("options", {deep: true})
    public onOptions(to: Options): void {
        const q: Dictionary<string> = {};

        if (to.range.from !== defaultOptions.range.from) {
            q.from = this.absoluteRange
                ? this.parseDate(to.range.from).valueOf().toString()
                : (to.range.from as string);
        }

        if (to.range.to !== defaultOptions.range.to) {
            q.to = this.absoluteRange ? this.parseDate(to.range.to).valueOf().toString() : (to.range.to as string);
        }

        if (to.refresh > 0) {
            q.refresh = to.refresh.toString();
        }

        updateRouteQuery(this.$route, q);

        this.updateRefresh();
    }

    @Watch("$route.query", {immediate: true})
    public onRouteQuery(to: Dictionary<string>): void {
        const options: Options = cloneDeep(defaultOptions);
        let n: number;

        if (to.from) {
            n = Number(to.from);
            options.range.from = !isNaN(n) ? this.formatDate(n, dateFormatRFC3339) : to.from;
        }

        if (to.to) {
            n = Number(to.to);
            options.range.to = !isNaN(n) ? this.formatDate(n, dateFormatRFC3339) : to.to;
        }

        if (to.refresh) {
            options.refresh = parseInt(to.refresh, 10);
        }

        Object.assign(this.options, options);

        this.getDashboard();
    }

    public refreshDashboard(): void {
        this.$root.$emit("item-refresh", null);
    }

    public resetTimeRange(): void {
        this.setTimeRange(Object.assign({}, defaultOptions.range));
    }

    public get timezoneUTC(): boolean {
        return this.$store.getters.timezoneUTC;
    }

    private getDashboard(): void {
        switch (this.params.type) {
            default: {
                this.$http
                    .post(`/api/v1/dashboards/${this.params.id}/resolve`, this.options.data || undefined)
                    .then(response => response.json())
                    .then(
                        (response: APIResponse<Dashboard>) => {
                            this.dashboard = response.data as Dashboard;

                            this.dashboardRefs =
                                this.dashboard?.references?.reduce((refs: Record<string, unknown>, ref: Reference) => {
                                    switch (ref.type) {
                                        case "chart":
                                            refs[`chart|${(ref.value as Chart).id}`] = ref.value;
                                            break;
                                    }

                                    return refs;
                                }, {}) ?? {};

                            delete this.dashboard.references;

                            this.$parent.$emit("dashboard-loaded", this.dashboard, this.dashboardRefs);
                            this.loading = false;
                        },
                        this.handleError(() => {
                            this.$parent.$emit("dashboard-loaded", null);
                            this.loading = false;
                        }),
                    );
            }
        }
    }

    public setRefreshInterval(interval: number | null = null): void {
        if (interval !== null) {
            this.options.refresh = interval;
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
            },
            (value: string | false) => {
                if (value !== false) {
                    this.setRefreshInterval(parseInt(value, 10));
                }
            },
        );
    }

    public setTimeRange(range: TimeRange | null = null, emit = true): void {
        if (range !== null) {
            Object.assign(this.options.range, range);
            if (emit) {
                this.$root.$emit("item-timerange", null, range);
            }
            return;
        }
        this.$components.modal(
            "time-range",
            {
                timeRange: this.absoluteRange
                    ? {
                          from: this.parseDate(this.options.range.from).format(datetimeFormat),
                          to: this.parseDate(this.options.range.to).format(datetimeFormat),
                      }
                    : {
                          from: "",
                          to: "",
                      },
            },
            (value: TimeRange) => {
                if (value !== false) {
                    this.setTimeRange(
                        {
                            from: this.parseDate(value.from, datetimeFormat).format(dateFormatRFC3339),
                            to: this.parseDate(value.to, datetimeFormat).format(dateFormatRFC3339),
                        },
                        true,
                    );
                }
            },
        );
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

    ::v-deep .v-grid-item {
        border: 0.15rem solid transparent;
        background-color: #272727;

        &:hover {
            border-color: var(--toolbar-background);
        }
    }
}
</style>
