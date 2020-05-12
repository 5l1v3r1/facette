<template>
    <div class="v-chart" :class="{empty: !value || !series}">
        <v-spinner v-if="loading"></v-spinner>

        <v-message v-else-if="!value || !series">{{ $t("messages.data.none") }}</v-message>

        <div class="v-chart-controls" v-if="controls">
            <v-toolbar v-if="!series">
                <v-button icon="sync" @click="refresh" v-tooltip="$t('labels.charts.refresh')"></v-button>
            </v-toolbar>

            <v-toolbar v-else>
                <v-button icon="sync" @click="refresh" v-tooltip="$t('labels.charts.refresh')"></v-button>
                <v-button
                    icon="history"
                    @click="update"
                    v-tooltip="$t('labels.timeRange.reset')"
                    v-if="!autoPropagate"
                ></v-button>
                <v-divider vertical></v-divider>
                <v-button
                    icon="search-minus"
                    @click="updateRange('zoom-out')"
                    v-tooltip="$t('labels.charts.zoomOut')"
                ></v-button>
                <v-button
                    icon="search-plus"
                    @click="updateRange('zoom-in')"
                    v-tooltip="$t('labels.charts.zoomIn')"
                ></v-button>
                <v-button
                    icon="arrows-alt-h"
                    @click="updateRange('propagate')"
                    v-tooltip="$t('labels.timeRange.propagate')"
                    v-if="!autoPropagate"
                ></v-button>
                <v-divider vertical></v-divider>
                <v-button
                    class="icon"
                    dropdown-anchor="bottom-right"
                    icon="angle-double-down"
                    :tooltip="$t('labels.moreActions')"
                >
                    <template slot="dropdown">
                        <v-button dropdown-anchor="right" icon="file-download">
                            {{ $t("labels.charts.export._") }}
                            <template slot="dropdown">
                                <v-button @click="downloadExport('png')">
                                    {{ $t("labels.charts.export.imagePNG") }}</v-button
                                >
                                <v-divider></v-divider>
                                <v-button @click="downloadExport('csv')">
                                    {{ $t("labels.charts.export.summaryCSV") }}</v-button
                                >
                                <v-button @click="downloadExport('json')">
                                    {{ $t("labels.charts.export.summaryJSON") }}</v-button
                                >
                            </template>
                        </v-button>
                        <template v-if="more">
                            <v-divider></v-divider>
                            <slot name="more"></slot>
                        </template>
                    </template>
                </v-button>
            </v-toolbar>

            <div class="v-chart-sliders">
                <div class="range">
                    <v-button
                        class="backward"
                        icon="arrow-left"
                        :class="{active: this.sliders.backward}"
                        @click="updateRange('backward')"
                    ></v-button>
                    <v-button
                        class="forward"
                        icon="arrow-right"
                        :class="{active: this.sliders.forward}"
                        @click="updateRange('forward')"
                    ></v-button>
                </div>
                <v-button class="summary" icon="list" :class="{active: sliders.summary}" @click="toggleSummary(true)">
                    {{ $t("labels.charts.showSummary") }}
                </v-button>
            </div>
        </div>

        <div ref="chart" @mousemove="controls && onMouse($event)" v-show="series"></div>
    </div>
</template>

<script lang="ts">
import Boula, {Config as BoulaConfig, Marker as BoulaMarker, Series as BoulaSeries} from "@facette/boula";
import * as d3 from "d3";
import dayjs from "dayjs";
import debounce from "lodash/debounce";
import ResizeObserver from "resize-observer-polyfill";
import slugify from "slugify";
import {Component, Mixins, Prop, Watch} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

export const dateFormatFilename = "YYYYMMDDHHmmss";

export const dateFormatRFC3339 = "YYYY-MM-DDTHH:mm:ss.SSSZ";

const mouseRange = 40;

@Component
export default class ChartComponent extends Mixins<CustomMixins>(CustomMixins) {
    @Prop({default: false, type: Boolean})
    public controls!: boolean;

    @Prop({default: false, type: Boolean})
    public emitCursor!: boolean;

    @Prop({default: false, type: Boolean})
    public legend!: boolean;

    @Prop({default: null, type: Object})
    public range!: TimeRange | null;

    @Prop({default: false, type: Boolean})
    public tooltip!: boolean;

    @Prop({required: true, validator: v => typeof v === "object" || v === undefined})
    public value!: Chart;

    public chartColors: Array<string> = [];

    public data: SeriesResult | null = null;

    public loading = true;

    public more = false;

    private chart: Boula | null = null;

    private domRect: DOMRect | null = null;

    private markers: Record<string, BoulaMarker> = {};

    private intersection: IntersectionObserver | null = null;

    private resize: ResizeObserver | null = null;

    private sliders: Record<string, boolean> = {
        backward: false,
        forward: false,
        summary: false,
    };

    public mounted(): void {
        this.checkSlots();

        this.$root.$on("item-refresh", this.onRefresh);
        this.$root.$on("item-timerange", this.onTimeRange);

        this.observeIntersection();

        if (this.value) {
            this.update();
        }
    }

    public updated(): void {
        this.checkSlots();
    }

    public beforeDestroy(): void {
        this.$root.$off("item-refresh", this.onRefresh);
        this.$root.$off("item-timerange", this.onTimeRange);

        this.unobserveIntersection();
        this.unobserveResize();

        if (this.chart) {
            this.chart.destroy();
        }
    }

    public downloadExport(type: "csv" | "json" | "png"): void {
        if (!this.data) {
            return;
        }

        const el: HTMLAnchorElement = document.createElement("a");

        const baseName = `${slugify(this.value.options?.title || this.value.name)}_${dayjs(this.data.from).format(
            dateFormatFilename,
        )}_${dayjs(this.data.to).format(dateFormatFilename)}`;

        switch (type) {
            case "csv":
            case "json": {
                let data: string;

                if (type === "csv") {
                    const summary = this.data.series.reduce((out: string, series: Series, index: number) => {
                        const keys: Array<string> = Object.keys(series.summary);
                        if (index === 0) {
                            out += `name,${keys.join(",")}\n`;
                        }
                        return `${out}"${series.name}",${keys.map((k: string) => series.summary[k]).join(",")}\n`;
                    }, "");

                    data = `data:text/csv;charset=utf-8,${encodeURIComponent(summary)}`;
                } else {
                    const summaries = this.data.series.reduce((out: Record<string, SeriesSummary>, series: Series) => {
                        out[series.name] = series.summary;
                        return out;
                    }, {});

                    data = `data:application/json,${encodeURIComponent(JSON.stringify(summaries, null, "\t"))}`;
                }

                Object.assign(el, {download: `${baseName}.${type}`, href: data});

                document.body.appendChild(el);
                el.click();
                document.body.removeChild(el);

                break;
            }

            case "png": {
                if (!this.chart) {
                    this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
                    return;
                }

                const data: string = (this.chart.canvas as HTMLCanvasElement).toDataURL("image/png");

                Object.assign(el, {download: `${baseName}.png`, href: data.replace("image/png", "image/octet-stream")});

                document.body.appendChild(el);
                el.click();
                document.body.removeChild(el);

                URL.revokeObjectURL(data);

                break;
            }
        }
    }

    public get series(): boolean {
        return Boolean(!this.loading && this.data?.series);
    }

    public get autoPropagate(): boolean {
        return this.$store.getters.autoPropagate;
    }

    @Watch("value", {deep: true})
    public onChange(): void {
        this.update();
    }

    public onMouse(e: MouseEvent): void {
        switch (e.type) {
            case "mouseleave": {
                const related: Element = e.relatedTarget as Element;
                const closest: Element | null = related ? related.closest(".v-chart") : null;
                if (!related || closest !== this.$el) {
                    this.sliders = {backward: false, forward: false, summary: false};
                }

                break;
            }

            case "mousemove": {
                if (!this.domRect) {
                    return;
                }

                const relX: number = e.pageX - this.domRect.x;
                const relY: number = e.pageY - this.domRect.y;

                if (!this.sliders.backward && !this.sliders.forward) {
                    if (relX <= mouseRange) {
                        this.sliders.backward = true;
                    } else if (relX >= this.domRect.width - mouseRange) {
                        this.sliders.forward = true;
                    }
                } else if (relX > mouseRange * 1.65 && relX < this.domRect.width - mouseRange * 1.65) {
                    Object.assign(this.sliders, {backward: false, forward: false});
                }

                if (relY >= this.domRect.height - mouseRange) {
                    if (!this.sliders.summary) {
                        this.sliders.summary = true;
                    }
                } else if (this.sliders.summary) {
                    this.sliders.summary = false;
                }

                break;
            }
        }
    }

    public refresh(): void {
        this.$root.$emit("item-refresh", this.$el);
    }

    public updateRange(mode: "backward" | "forward" | "propagate" | "zoom-in" | "zoom-out"): void {
        if (!this.data) {
            return;
        }

        const from: dayjs.Dayjs = dayjs(this.data.from);
        const to: dayjs.Dayjs = dayjs(this.data.to);
        let delta: number;

        switch (mode) {
            case "backward":
                delta = to.diff(from, "second") * 0.25;
                from.add(-delta, "second");
                to.add(-delta, "second");
                break;

            case "forward":
                delta = to.diff(from, "second") * 0.25;
                from.add(delta, "second");
                to.add(delta, "second");
                break;

            case "zoom-in":
                delta = to.diff(from, "second") * 0.25;
                from.add(delta, "second");
                to.add(-delta, "second");
                break;

            case "zoom-out":
                delta = to.diff(from, "second") * 0.5;
                from.add(-delta, "second");
                to.add(delta, "second");
                break;
        }

        const range: TimeRange = {
            from: from.format(dateFormatRFC3339),
            to: to.format(dateFormatRFC3339),
        };

        if (this.autoPropagate || mode === "propagate") {
            this.$root.$emit("item-timerange", null, range);
        } else {
            this.update(range);
        }
    }

    private checkSlots(): void {
        this.more = Boolean(this.$slots?.more);
    }

    private draw(): void {
        const refChart: HTMLDivElement = this.$refs.chart as HTMLDivElement;
        if (!refChart) {
            return;
        }

        const markers: Array<BoulaMarker> = Object.values(this.markers);
        this.value.options?.axes?.y?.left?.constants?.forEach(constant => {
            markers.push({label: constant.label || true, y: constant.value});
        });

        let series: Array<BoulaSeries>;
        if (this.data?.series) {
            const disabledSeries: Array<string> = this.chart
                ? this.chart.config.series.reduce((names: Array<string>, s: BoulaSeries) => {
                      if (s.disabled && s.name) {
                          names.push(s.name);
                      }
                      return names;
                  }, [])
                : [];

            series = this.data.series.map((series: Series, index: number) => {
                const bs: BoulaSeries = {
                    name: series.name,
                    points: series.points.map((p: Point) => ({0: p[0] * 1000, 1: p[1]})),
                };

                if (disabledSeries.includes(series.name)) {
                    bs.disabled = true;
                }

                const options: ChartSeriesOptions | undefined = this.value.series
                    ? this.value.series[index]?.options
                    : undefined;
                if (options && options.color) {
                    bs.color = options.color;
                }

                return bs;
            });
        } else {
            series = [];
        }

        const config: BoulaConfig = {
            axes: {
                x: {
                    grid: false,
                    max: (this.data && dayjs(this.data.to).valueOf()) || undefined,
                    min: (this.data && dayjs(this.data.from).valueOf()) || undefined,
                    ticks: {
                        count: Math.max(Math.floor(refChart.clientWidth / 80), 2),
                        format: (date: Date): string => {
                            const format: (specifier: string) => (date: Date) => string = this.timezoneUTC
                                ? d3.utcFormat
                                : d3.timeFormat;

                            return (d3.timeSecond(date) < date
                                ? format(".%L")
                                : d3.timeMinute(date) < date
                                ? format(":%S")
                                : d3.timeHour(date) < date
                                ? format("%H:%M")
                                : d3.timeDay(date) < date
                                ? format("%H:00")
                                : d3.timeMonth(date) < date
                                ? format("%a %d")
                                : d3.timeYear(date) < date
                                ? format("%B")
                                : d3.timeFormat("%Y"))(date);
                        },
                    },
                },
                y: {
                    center: this.value.options?.axes?.y?.center ?? false,
                    stack: this.value.options?.axes?.y?.left?.stack || false,
                    label: {
                        text: this.value.options?.axes?.y?.left?.label,
                    },
                    max: this.value.options?.axes?.y?.left?.max,
                    min: this.value.options?.axes?.y?.left?.min,
                    ticks: {
                        draw: false,
                        format: v => this.formatValue(v, this.value.options?.axes?.y?.left?.unit),
                    },
                },
            },
            bindTo: refChart,
            cursor: {
                enabled: this.tooltip,
            },
            events: {
                afterDraw: () => {
                    if (this.chart) {
                        const defaultColors = this.chart.config.colors as Array<string>;
                        this.chartColors = this.chart.config.series.map(
                            (s: BoulaSeries, index: number) => s.color || defaultColors[index % defaultColors.length],
                        );
                    }
                },
                cursorMove: (date: Date | null) => {
                    this.$root.$emit("item-cursor", date, this.$el);
                },
                select: (from: Date, to: Date) => {
                    if (to > from) {
                        const range: TimeRange = {
                            from: dayjs(from).format(dateFormatRFC3339),
                            to: dayjs(to).format(dateFormatRFC3339),
                        };

                        if (this.autoPropagate) {
                            this.$root.$emit("item-timerange", null, range);
                        } else {
                            this.update(range);
                        }
                    }
                },
            },
            legend: {
                enabled: this.legend,
            },
            markers,
            selection: {
                enabled: this.controls,
            },
            series,
            tooltip: {
                enabled: this.tooltip,
                date: {
                    format: d => this.formatDate(d),
                },
            },
            title: {
                text: this.value.options?.title,
            },
            type: this.value.options?.type || "area",
        };

        if (!this.chart) {
            this.chart = new Boula(config);
        } else {
            this.chart.update(config);
        }

        this.chart.draw();

        // Set DOMRect for sliders area detection
        this.domRect = this.$el.getBoundingClientRect() as DOMRect;

        if (this.resize === null) {
            this.observeResize();
        }
    }

    private observeIntersection(): void {
        this.intersection = new IntersectionObserver(
            (entries: Array<IntersectionObserverEntry>) => {
                if (entries[0].intersectionRatio > 0) {
                    this.unobserveIntersection();
                    this.update(this.$store.getters.timeRange); // Use store time range to avoid shifts between charts
                }
            },
            {threshold: 0},
        );

        this.intersection.observe(this.$el);
    }

    private observeResize(): void {
        this.resize = new ResizeObserver(
            debounce(() => {
                if (this.$refs.chart) {
                    requestAnimationFrame(() => {
                        this.draw();
                    });
                }
            }, 200),
        );

        this.resize.observe(this.$refs.chart as HTMLElement);
    }

    private onRefresh(target: Element | null = null): void {
        if (target === null || target === this.$el) {
            this.update();
        }
    }

    private onTimeRange(target: Element | null = null, range: TimeRange | null = null): void {
        if (target === null || target === this.$el) {
            this.update(range);
        }
    }

    private get timezoneUTC(): boolean {
        return this.$store.getters.timezoneUTC;
    }

    private unobserveIntersection(): void {
        if (this.intersection) {
            this.intersection.disconnect();
            this.intersection = null;
        }
    }

    private unobserveResize(): void {
        if (this.resize) {
            this.resize.disconnect();
            this.resize = null;
        }
    }

    private update(range: TimeRange | null = null): void {
        if (this.intersection !== null) {
            return;
        } else if (!this.value?.series || this.value.series.length === 0) {
            this.data = null;
            this.loading = false;
            return;
        }

        const query: SeriesQuery = {
            exprs: this.value.series.map(series => series.expr),
        };

        if (range) {
            Object.assign(query, range);
        } else if (this.range) {
            Object.assign(query, this.range);
        }

        this.loading = true;

        this.$http
            .post("/api/v1/query", query)
            .then(response => response.json())
            .then(
                (response: APIResponse<SeriesResult>) => {
                    this.data = response.data ?? null;
                    this.loading = false;

                    if (this.data) {
                        this.$root.$emit("item-loaded", {from: this.data.from, to: this.data.to});
                        this.$nextTick(() => {
                            this.draw();
                        });
                    }
                },
                this.handleError(() => {
                    this.loading = false;
                }),
            );
    }
}
</script>

<style lang="scss" scoped>
@import "~@facette/boula/dist/style.css";

.v-chart {
    background-color: inherit;
    height: 100%;
    position: relative;

    &.empty {
        align-items: center;
        display: flex;
        justify-content: center;
    }

    .chart-container {
        background-color: inherit;
        height: 100%;
        width: 100%;
    }

    .v-chart-controls {
        bottom: 0;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        left: 0;
        pointer-events: none;
        position: absolute;
        right: 0;
        top: 0;
        visibility: hidden;
        z-index: 1;

        .v-toolbar {
            align-self: flex-end;
            border-bottom-left-radius: 0.2rem;
            height: 2rem;
            min-height: 2rem;
            padding: 0 0 0.5rem;
            pointer-events: auto;

            > .v-button {
                height: 1.5rem;
                line-height: 1.5rem;
                min-height: 1.5rem;
                min-width: 1.5rem;

                &:last-child {
                    margin-right: 0;
                }

                &.icon {
                    width: 1.5rem;

                    .v-icon {
                        font-size: 0.8rem;
                    }
                }
            }

            > .v-divider {
                height: 1.3rem;
            }
        }

        .v-chart-sliders {
            align-items: center;
            display: flex;
            height: calc(50% + 2rem);
            flex-direction: column;
            justify-content: space-between;
            overflow: hidden;

            .v-button {
                pointer-events: auto;
                transition: transform 0.2s var(--timing-function);
                will-change: transform;

                &.active {
                    transform: none !important;
                }

                ::v-deep .v-button-content {
                    background-color: var(--toolbar-background);
                }
            }

            .range {
                align-items: center;
                display: flex;
                justify-content: space-between;
                width: 100%;

                .v-button {
                    font-size: 1.35rem;
                    height: 4rem;
                    line-height: 4rem;
                    width: 2.25rem;

                    &.backward {
                        transform: translateX(-100%);

                        ::v-deep .v-button-content {
                            border-radius: 0 0.2rem 0.2rem 0;
                        }
                    }

                    &.forward {
                        transform: translateX(100%);

                        ::v-deep .v-button-content {
                            border-radius: 0.2rem 0 0 0.2rem;
                        }
                    }
                }
            }

            .summary {
                font-size: 0.75rem;
                height: 1.5rem;
                line-height: 1.5rem;
                transform: translateY(100%);

                ::v-deep .v-button-content {
                    border-radius: 0.2rem 0.2rem 0 0;
                }
            }
        }
    }

    &:hover .v-chart-controls {
        visibility: visible;
    }

    .v-message {
        color: var(--gray);
        font-size: 1rem;
    }

    ::v-deep {
        .selection {
            box-sizing: content-box;
        }

        .tooltip {
            z-index: 700;
        }
    }
}
</style>
