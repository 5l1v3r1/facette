/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import dayjs from "dayjs";
import dayjsUTC from "dayjs/plugin/utc";
import Vue from "vue";
import {Mixin} from "vue-mixin-decorator";
import {HttpResponse} from "vue-resource/types/vue_resource";
import {Dictionary} from "vue-router/types/router";

dayjs.extend(dayjsUTC);

import {
    BinaryUnit,
    CountUnit,
    DurationUnit,
    MetricUnit,
    formatBinary,
    formatCount,
    formatDuration,
    formatMetric,
    formatNumber,
} from "./helpers/format";

@Mixin
export class CustomMixins extends Vue {
    public erred = false;

    private unwatchGuard: (() => void) | null = null;

    public mounted(): void {
        this.$root.$on("set-error", this.onSetError);
    }

    public beforeDestroy(): void {
        this.$root.$off("set-error", this.onSetError);
    }

    public formatDate(input: Date | string, format = "MMM D YYYY, HH:mm:ss"): string {
        if (this.$store.getters.timezoneUTC) {
            return dayjs(input).utc().format(format) + " UTC";
        }

        return dayjs(input).format(format);
    }

    public formatValue(input: number | null, unit?: Unit, decimals = 2): string {
        if (input === null) {
            return "null";
        }

        if (unit !== undefined) {
            switch (unit.type) {
                case "binary": {
                    return formatBinary(input, decimals, unit.base as BinaryUnit);
                }

                case "count": {
                    return formatCount(input, decimals, unit.base as CountUnit);
                }

                case "duration": {
                    return formatDuration(input, decimals, unit.base as DurationUnit);
                }

                case "metric": {
                    return formatMetric(input, decimals, unit.base as MetricUnit);
                }
            }
        }

        return formatNumber(input, decimals);
    }

    public get guarded(): boolean {
        return this.$store.getters.routeGuarded;
    }

    public guardWatch(...expOrFn: Array<string>): void {
        this.unguard();

        this.unwatchGuard = this.$watch(
            () => expOrFn.map((s: string) => (this as Record<string, unknown>)[s]),
            () => {
                this.guard(true);
                this.guardUnwatch();
            },
            {deep: true},
        );
    }

    public handleError(handler?: () => void): (response: HttpResponse) => void {
        return (response: HttpResponse) => {
            if (response.status >= 500) {
                this.$root.$emit("set-error", true);
            }

            this.$components.notify(
                (response.data?.error
                    ? this.$t("messages.error._", [response.data.error])
                    : this.$t("messages.error.unhandled")) as string,
                "error",
            );

            if (handler) {
                handler();
            }
        };
    }

    public get modifierPressed(): boolean {
        return this.$store.getters.modifierPressed;
    }

    public get params(): Dictionary<string> {
        return this.$route.params;
    }

    public resetError(): void {
        this.$root.$emit("set-error", false);
    }

    public toggleSidebar(): void {
        this.sidebar = !this.sidebar;
    }

    public unguard(): void {
        this.guard(false);
        this.guardUnwatch();
    }

    private guard(state: boolean): void {
        this.$store.commit("routeGuarded", state);
        this.$store.commit("routeGuardFn", state ? this.onBeforeUnload : null);

        window[state ? "addEventListener" : "removeEventListener"]("beforeunload", this.onBeforeUnload);
    }

    private guardUnwatch(): void {
        if (this.unwatchGuard) {
            this.unwatchGuard();
            this.unwatchGuard = null;
        }
    }

    private onBeforeUnload(e: Event): void {
        e.preventDefault();
        e.returnValue = false;
    }

    private onSetError(state: boolean): void {
        this.erred = state;
    }

    private get sidebar(): boolean {
        return this.$store.getters.sidebar;
    }

    private set sidebar(value: boolean) {
        this.$store.commit("sidebar", value);
    }
}
