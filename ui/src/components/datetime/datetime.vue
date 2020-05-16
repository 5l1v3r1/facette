<template>
    <div class="v-datetime">
        <v-input
            icon="calendar-alt"
            ref="input"
            :label="label"
            :placeholder="placeholder"
            @input="$emit('input', time.value)"
            v-model="time.value"
        ></v-input>

        <table class="v-datetime-picker" @mousedown.prevent v-if="table">
            <thead class="v-datetime-header">
                <tr class="v-datetime-row">
                    <th>
                        <v-icon
                            icon="arrow-left"
                            @click.native="add(-1 * table.header.delta.count, table.header.delta.unit)"
                        ></v-icon>
                    </th>
                    <th
                        class="v-datetime-title"
                        :colspan="table.header.title.span"
                        @click="table.header.title.nextPane && switchPane(table.header.title.nextPane)"
                    >
                        {{ table.header.title.value }}
                    </th>
                    <th>
                        <v-icon
                            icon="arrow-right"
                            @click.native="add(table.header.delta.count, table.header.delta.unit)"
                        ></v-icon>
                    </th>
                </tr>
                <tr class="v-datetime-row" v-if="pane == 'day'">
                    <th class="v-datetime-cell" :key="day" v-for="day in weekdays">{{ day }}</th>
                </tr>
            </thead>
            <tbody class="v-datetime-body">
                <tr class="v-datetime-row" :key="index" v-for="(row, index) in table.rows">
                    <td
                        class="v-datetime-cell"
                        :aria-disabled="cell.disabled"
                        :class="{
                            current: cell.current,
                            future: cell.future,
                            past: cell.past,
                            today: cell.today,
                        }"
                        :key="subindex"
                        @click="select(cell.time)"
                        v-for="(cell, subindex) in row"
                    >
                        <span>{{ cell.value }}</span>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script lang="ts">
import dayjs from "dayjs";
import dayjsIsSameOrBefore from "dayjs/plugin/isSameOrBefore";
import dayjsIsToday from "dayjs/plugin/isToday";
import dayjsLocaleData from "dayjs/plugin/localeData";
import {Component, Mixins, Prop, Watch} from "vue-property-decorator";

import InputComponent from "@/src/components/vue/input/input.vue";
import {CustomMixins} from "@/src/mixins";

dayjs.extend(dayjsIsSameOrBefore);
dayjs.extend(dayjsIsToday);
dayjs.extend(dayjsLocaleData);

interface Header {
    delta: {
        count: number;
        unit: dayjs.OpUnitType;
    };
    title: {
        nextPane: string | null;
        span: number;
        value: string;
    };
}

interface Cell {
    current: boolean;
    disabled: boolean;
    future?: boolean;
    past?: boolean;
    time: dayjs.Dayjs;
    today?: boolean;
    value: string;
}

interface Table {
    header: Header;
    rows: Array<Array<Cell>>;
}

type Pane = "day" | "hour" | "minute" | "month" | "year";

interface Time {
    min: dayjs.Dayjs | null;
    ref: dayjs.Dayjs | null;
    value: string;
}

function chunkArray(input: Array<Cell>, size: number): Array<Array<Cell>> {
    const output: Array<Array<Cell>> = [];
    while (input.length) {
        output.push(input.splice(0, size));
    }

    return output;
}

@Component
export default class DatetimeComponent extends Mixins<CustomMixins>(CustomMixins) {
    @Prop({default: "YYYY-MM-DD HH:mm:ss", type: String})
    public format!: string;

    @Prop({default: null, type: String})
    public label!: string | null;

    @Prop({default: null, type: String})
    public min!: string | null;

    @Prop({default: null, type: String})
    public placeholder!: string | null;

    @Prop({required: true, type: String})
    public value!: string;

    public pane: Pane = "day";

    public rows: Array<Array<Cell>> = [];

    public time: Time = {
        min: null,
        ref: null,
        value: "",
    };

    public add(count: number, unit: dayjs.OpUnitType): void {
        if (this.time.ref !== null) {
            this.time.ref = this.time.ref?.add(count, unit);
        }
    }

    public focus(select: boolean): void {
        (this.$refs.input as InputComponent).focus(select);
    }

    @Watch("value", {immediate: true})
    public onChange(to: string): void {
        const time = this.parseDate(to || undefined, this.format);
        if (!time.isValid()) {
            return;
        }

        Object.assign(this.time, {
            ref: time.startOf("day"),
            value: to,
        });
    }

    @Watch("min", {immediate: true})
    public onMin(to: string, from: string | undefined): void {
        const time = this.parseDate(to, this.format);
        if (!time.isValid()) {
            return;
        }

        // If minimal time is after the reference one, update reference keeping
        // current delta.
        if (this.time.ref !== null && time.isAfter(this.time.ref) && from) {
            this.time.ref.add(time.diff(this.parseDate(from, this.format), "second"), "second");
            this.time.value = this.time.ref.format(this.format);
        }

        this.time.min = time;
    }

    public select(value: dayjs.Dayjs): void {
        let pane = this.pane;
        const time = Object.assign({}, this.time, {ref: value});

        switch (this.pane) {
            case "day":
                pane = "hour";
                break;

            case "hour":
                pane = "minute";
                break;

            case "minute":
                pane = "day";
                time.value = value.format(this.format);
                this.$emit("input", time.value);
                break;

            case "month":
                pane = "day";
                break;

            case "year":
                pane = "month";
                break;
        }

        Object.assign(this, {
            pane,
            time,
        });
    }

    public switchPane(pane: Pane): void {
        this.pane = pane;
    }

    public get table(): Table | null {
        if (!this.time.ref) {
            return null;
        }

        const time = this.parseDate(this.time.value, this.format);

        switch (this.pane) {
            case "day": {
                const month = this.time.ref.month();
                let day = this.time.ref.startOf("month").startOf("week");
                const end = this.time.ref.endOf("month").endOf("week");
                const days: Array<Cell> = [];

                while (day.isBefore(end)) {
                    days.push({
                        current: time && day.isSame(time, "day"),
                        disabled: this.time.min !== null && day.isBefore(this.time.min, "day"),
                        future: day.month() > month,
                        past: day.month() < month,
                        time: day,
                        today: day.isToday(),
                        value: day.format("D"),
                    });

                    day = day.add(1, "day");
                }

                return {
                    header: {
                        delta: {
                            count: 1,
                            unit: "month",
                        },
                        title: {
                            nextPane: "month",
                            span: 5,
                            value: this.time.ref.format("MMM YYYY"),
                        },
                    },
                    rows: chunkArray(days, 7),
                };
            }

            case "hour": {
                let hour = this.time.ref.startOf("day");
                const end = this.time.ref.endOf("day");
                const hours: Array<Cell> = [];

                while (hour.isBefore(end)) {
                    hours.push({
                        current: time && hour.isSame(time, "hour"),
                        disabled: this.time.min !== null && hour.isBefore(this.time.min, "hour"),
                        time: hour,
                        value: hour.format("HH:mm"),
                    });

                    hour = hour.add(1, "hour");
                }

                return {
                    header: {
                        delta: {
                            count: 1,
                            unit: "day",
                        },
                        title: {
                            nextPane: null,
                            span: 2,
                            value: this.time.ref.format("MMM D, YYYY"),
                        },
                    },
                    rows: chunkArray(hours, 4),
                };
            }

            case "minute": {
                let minute = this.time.ref.startOf("hour");
                const end = this.time.ref.endOf("hour");
                const minutes: Array<Cell> = [];

                while (minute.isBefore(end)) {
                    minutes.push({
                        current: time && minute.isSame(time, "minute"),
                        disabled: this.time.min !== null && minute.isSameOrBefore(this.time.min, "minute"),
                        time: minute.clone(),
                        value: minute.format("HH:mm"),
                    });

                    minute = minute.add(5, "minute");
                }

                return {
                    header: {
                        delta: {
                            count: 1,
                            unit: "hour",
                        },
                        title: {
                            nextPane: null,
                            span: 2,
                            value: this.time.ref.format("MMM D, YYYY HH:mm"),
                        },
                    },
                    rows: chunkArray(minutes, 4),
                };
            }

            case "month": {
                let month = this.time.ref.startOf("year");
                const end = this.time.ref.endOf("year");
                const months: Array<Cell> = [];

                while (month.isBefore(end)) {
                    months.push({
                        current: time && month.isSame(time, "month"),
                        disabled: this.time.min !== null && month.isBefore(this.time.min, "month"),
                        time: month.clone(),
                        value: month.format("MMM"),
                    });

                    month = month.add(1, "month");
                }

                return {
                    header: {
                        delta: {
                            count: 1,
                            unit: "year",
                        },
                        title: {
                            nextPane: "year",
                            span: 2,
                            value: this.time.ref.format("YYYY"),
                        },
                    },
                    rows: chunkArray(months, 4),
                };
            }

            case "year": {
                let year = this.parseDate()
                    .year(Math.floor(this.time.ref.year() / 10) * 10 - 1)
                    .startOf("year");
                const start = year.clone();
                const end = year.add(11, "year").endOf("year");
                const years: Array<Cell> = [];

                while (year.isBefore(end)) {
                    years.push({
                        current: time && year.isSame(time, "year"),
                        disabled: this.time.min !== null && year.isBefore(this.time.min, "year"),
                        future: years.length === 11,
                        past: years.length === 0,
                        time: year.clone(),
                        value: year.format("YYYY"),
                    });

                    year = year.add(1, "year");
                }

                return {
                    header: {
                        delta: {
                            count: 12,
                            unit: "year",
                        },
                        title: {
                            nextPane: null,
                            span: 2,
                            value: `${start.add(1, "year").format("YYYY")}-${end.add(-1, "year").format("YYYY")}`,
                        },
                    },
                    rows: chunkArray(years, 4),
                };
            }

            default:
                return null;
        }
    }

    public get weekdays(): Array<string> {
        return this.parseDate().localeData().weekdaysShort();
    }
}
</script>

<style lang="scss" scoped>
.v-datetime {
    display: inline-block;
    width: 18rem;

    .v-input {
        width: 100%;
    }

    .v-datetime-picker {
        border-spacing: 0;
        border: 1px solid var(--dark-gray);
        height: 18rem;
        margin-top: 1rem;
        table-layout: fixed;
        width: 100%;
    }

    .v-datetime-title,
    .v-datetime-header .v-icon,
    .v-datetime-body .v-datetime-cell {
        cursor: pointer;
    }

    .v-datetime-header {
        .v-datetime-row {
            background-color: var(--toolbar-background);
            height: 2rem;
        }

        .v-datetime-cell {
            font-weight: normal;
        }

        .v-icon {
            font-size: 1rem;
        }
    }

    .v-datetime-title {
        white-space: nowrap;
    }

    .v-datetime-row {
        line-height: 2rem;
    }

    .v-datetime-cell {
        padding: 0.25rem;
        position: relative;
        text-align: center;

        &.future span,
        &.past span {
            color: var(--gray);
        }

        &[aria-disabled] {
            pointer-events: none;

            span {
                color: var(--dark-gray);
            }

            &.today span {
                color: var(--background);
            }
        }

        span {
            border-radius: 0.2rem;
            display: block;
        }

        &.current span {
            background-color: var(--blue);
        }

        &.today span {
            background-color: var(--dark-gray);
        }
    }

    .v-datetime-body {
        .v-datetime-cell {
            &:active span,
            &:focus span,
            &:hover span {
                background-color: var(--blue);
                color: var(--white);
            }
        }
    }
}
</style>
