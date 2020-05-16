<template>
    <v-content>
        <v-modal-prompt></v-modal-prompt>
        <v-modal-time-range></v-modal-time-range>

        <v-toolbar clip="content">
            <v-button
                icon="sync"
                :icon-badge="!erred && options.refresh ? formatValue(options.refresh, {type: 'duration'}, 0) : null"
                :disabled="loading || erred"
                :tooltip="$t(`labels.${params.type}.refresh`)"
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

            <v-button icon="calendar-alt" :disabled="loading || erred" :tooltip="$t('labels.timeRange.set')">
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

import {dateFormatDisplay, dateFormatRFC3339, defaultTimeRange, ranges} from "@/src/components/chart/chart.vue";
import {CustomMixins} from "@/src/mixins";
import {updateRouteQuery} from "@/src/router";

interface Options {
    timeRange: TimeRange;
    refresh: number;
}

const defaultOptions = {
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

    public intervals: Array<number> = intervals;

    public loading = true;

    public options = cloneDeep(defaultOptions);

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
        this.$root.$on("item-loaded", this.onItemLoaded);
        this.$root.$on("item-timerange", this.onItemTimeRange);
    }

    public beforeDestroy(): void {
        document.removeEventListener("visibilitychange", this.onVisibilityChange);
        this.$root.$off("item-loaded", this.onItemLoaded);
        this.$root.$off("item-timerange", this.onItemTimeRange);
    }

    public get absoluteRange(): boolean {
        return (
            this.parseDate(this.options.timeRange.from).isValid() && this.parseDate(this.options.timeRange.to).isValid()
        );
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

    public get canResetTimeRange(): boolean {
        return !isEqual(this.options.timeRange, defaultTimeRange) || !this.timeRangeSynced;
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

        updateRouteQuery(this.$route, q);

        this.updateRefresh();
    }

    @Watch("$route.query", {immediate: true})
    public onRouteQuery(to: Dictionary<string>): void {
        const options = cloneDeep(defaultOptions);
        let n: number;

        if (to.from) {
            n = Number(to.from);
            options.timeRange.from = !isNaN(n) ? this.formatDate(n, dateFormatRFC3339) : to.from;
        }

        if (to.to) {
            n = Number(to.to);
            options.timeRange.to = !isNaN(n) ? this.formatDate(n, dateFormatRFC3339) : to.to;
        }

        if (to.refresh) {
            options.refresh = parseInt(to.refresh, 10) || 0;
        }

        this.options = options;
        this.getDashboard();
    }

    public refreshDashboard(): void {
        this.$root.$emit("item-refresh", null);
    }

    public resetTimeRange(): void {
        this.setTimeRange(Object.assign({}, defaultTimeRange));
    }

    public get timezoneUTC(): boolean {
        return this.$store.getters.timezoneUTC;
    }

    private getDashboard(): void {
        switch (this.params.type) {
            default: {
                this.$http
                    .post(`/api/v1/dashboards/${this.params.id}/resolve`) // TODO: handle data
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
            },
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
            },
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

    ::v-deep .v-grid-item {
        border: 0.15rem solid transparent;
        background-color: #272727;

        &:hover {
            border-color: var(--toolbar-background);
        }
    }
}
</style>
