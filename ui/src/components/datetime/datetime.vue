<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div class="v-datetime" @autofocus="focus($event.detail.select)">
        <v-input
            icon="calendar-alt"
            ref="input"
            :custom-validity="validity"
            :label="label"
            :placeholder="placeholder"
            :value="value"
            @update:value="$emit('update:value', $event)"
        ></v-input>

        <table class="v-datetime-picker" @mousedown.prevent>
            <thead class="v-datetime-header">
                <tr class="v-datetime-row">
                    <th>
                        <v-icon
                            icon="arrow-left"
                            @click="add(-1 * table.header.delta.count, table.header.delta.unit)"
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
                            @click="add(table.header.delta.count, table.header.delta.unit)"
                        ></v-icon>
                    </th>
                </tr>

                <tr class="v-datetime-row" v-if="pane == 'day'">
                    <th class="v-datetime-cell" :key="day" v-for="day in weekdays">{{ day }}</th>
                </tr>

                <tr class="v-datetime-row placeholder" v-else>
                    <th class="v-datetime-cell" colspan="4"></th>
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
                        {{ cell.value }}
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import {DateTime} from "luxon";
import {ComponentPublicInstance, SetupContext, computed, ref, watch} from "vue";

import {InputComponent} from "types/ui";

import {dateFormatDisplay} from "@/components/chart";
import {parseDate, weekdaysShort} from "@/helpers/date";

interface Header {
    delta: {
        count: number;
        unit: string;
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
    time: DateTime;
    today?: boolean;
    value: string;
}

interface Table {
    header: Header;
    rows: Array<Array<Cell>>;
}

type Pane = "day" | "hour" | "minute" | "month" | "year";

function chunkArray(input: Array<Cell>, size: number): Array<Array<Cell>> {
    const output: Array<Array<Cell>> = [];
    while (input.length) {
        output.push(input.splice(0, size));
    }

    return output;
}

export default {
    props: {
        format: {
            default: dateFormatDisplay,
            type: String,
        },
        label: {
            default: null,
            type: String,
        },
        min: {
            default: null,
            type: String,
        },
        placeholder: {
            default: null,
            type: String,
        },
        value: {
            required: true,
            type: String,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const input = ref<ComponentPublicInstance<InputComponent> | null>(null);
        const minTime = ref<DateTime | null>(null);
        const pane = ref<Pane>("day");
        const refTime = ref<DateTime | null>(null);
        const weekdays = ref<Array<string>>(weekdaysShort());

        const table = computed(
            (): Table => {
                const time = parseDate(props.value || undefined, props.format);

                if (refTime.value === null) {
                    refTime.value = time.isValid ? cloneDeep(time) : parseDate();
                }

                switch (pane.value) {
                    case "day": {
                        let day = refTime.value.startOf("month").startOf("week");

                        const month = refTime.value.month;
                        const end = refTime.value.endOf("month").endOf("week");
                        const days: Array<Cell> = [];

                        while (day < end) {
                            days.push({
                                current: props.value && day.hasSame(time, "day"),
                                disabled: minTime.value !== null && day < minTime.value.startOf("day"),
                                future: day.month > month,
                                past: day.month < month,
                                time: day,
                                today: day.hasSame(parseDate(), "day"),
                                value: day.toFormat("d"),
                            });

                            day = day.plus({days: 1});
                        }

                        return {
                            header: {
                                delta: {count: 1, unit: "month"},
                                title: {nextPane: "month", span: 5, value: refTime.value.toFormat("MMM yyyy")},
                            },
                            rows: chunkArray(days, 7),
                        };
                    }

                    case "hour": {
                        let hour = refTime.value.startOf("day");

                        const end = refTime.value.endOf("day");
                        const hours: Array<Cell> = [];

                        while (hour < end) {
                            hours.push({
                                current: props.value && hour.hasSame(time, "hour"),
                                disabled: minTime.value !== null && hour <= minTime.value.startOf("hour"),
                                time: hour,
                                value: hour.toFormat("HH:mm"),
                            });

                            hour = hour.plus({hours: 1});
                        }

                        return {
                            header: {
                                delta: {count: 1, unit: "day"},
                                title: {nextPane: null, span: 2, value: refTime.value.toFormat("MMM d, yyyy")},
                            },
                            rows: chunkArray(hours, 4),
                        };
                    }

                    case "minute": {
                        let minute = refTime.value.startOf("hour");

                        const end = refTime.value.endOf("hour");
                        const minutes: Array<Cell> = [];

                        while (minute < end) {
                            minutes.push({
                                current: props.value && minute.hasSame(time, "minute"),
                                disabled: minTime.value !== null && minute <= minTime.value.startOf("minute"),
                                time: minute,
                                value: minute.toFormat("HH:mm"),
                            });

                            minute = minute.plus({minutes: 5});
                        }

                        return {
                            header: {
                                delta: {count: 1, unit: "hour"},
                                title: {nextPane: null, span: 2, value: refTime.value.toFormat("MMM d, yyyy HH:mm")},
                            },
                            rows: chunkArray(minutes, 4),
                        };
                    }

                    case "month": {
                        let month = refTime.value.startOf("year");

                        const end = refTime.value.endOf("year");
                        const months: Array<Cell> = [];

                        while (month < end) {
                            months.push({
                                current: props.value && month.hasSame(time, "month"),
                                disabled: minTime.value !== null && month <= minTime.value.startOf("month"),
                                time: month,
                                value: month.toFormat("MMM"),
                            });

                            month = month.plus({months: 1});
                        }

                        return {
                            header: {
                                delta: {count: 1, unit: "year"},
                                title: {nextPane: "year", span: 2, value: refTime.value.toFormat("yyyy")},
                            },
                            rows: chunkArray(months, 4),
                        };
                    }

                    case "year": {
                        let year = DateTime.local(Math.floor(refTime.value.year / 10) * 10 - 1);

                        const start = year;
                        const end = year.plus({years: 11}).endOf("year");
                        const years: Array<Cell> = [];

                        while (year < end) {
                            years.push({
                                current: props.value && year.hasSame(refTime.value, "year"),
                                disabled: minTime.value !== null && year < minTime.value.startOf("year"),
                                future: years.length === 11,
                                past: years.length === 0,
                                time: year,
                                value: year.toFormat("yyyy"),
                            });

                            year = year.plus({years: 1});
                        }

                        return {
                            header: {
                                delta: {count: 12, unit: "year"},
                                title: {
                                    nextPane: null,
                                    span: 2,
                                    value: `${start.plus({years: 1}).toFormat("yyyy")}-${end
                                        .minus({years: 1})
                                        .toFormat("yyyy")}`,
                                },
                            },
                            rows: chunkArray(years, 4),
                        };
                    }
                }
            },
        );

        const add = (count: number, unit: string): void => {
            if (refTime.value !== null) {
                refTime.value = refTime.value.plus({[unit]: count});
            }
        };

        const focus = (select: boolean): void => {
            input.value?.focus(select);
        };

        const select = (value: DateTime): void => {
            refTime.value = value;

            switch (pane.value) {
                case "day":
                    pane.value = "hour";
                    break;

                case "hour":
                    pane.value = "minute";
                    break;

                case "minute":
                    pane.value = "day";
                    ctx.emit("update:value", value.toFormat(props.format));
                    break;

                case "month":
                    pane.value = "day";
                    break;

                case "year":
                    pane.value = "month";
                    break;
            }
        };

        const switchPane = (value: Pane): void => {
            pane.value = value;
        };

        const validity = (value: string): Promise<string> => {
            return Promise.resolve(parseDate(value, props.format).isValid ? "" : "invalid date"); // FIXME: handle i18n
        };

        watch(
            () => props.min,
            (to, from) => {
                let time = parseDate(to, props.format);
                if (!time.isValid) {
                    return;
                }

                // If minimal time is after the reference one, update reference keeping
                // current delta.
                if (refTime.value !== null && time > refTime.value && from) {
                    refTime.value = refTime.value.plus(time.diff(parseDate(from, props.format), "second"));
                    ctx.emit("update:value", refTime.value.toFormat(props.format));
                }

                minTime.value = time;
            },
        );

        return {
            add,
            focus,
            input,
            pane,
            select,
            switchPane,
            table,
            validity,
            weekdays,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-datetime {
    background-color: var(--input-background);
    border-radius: 0.2rem;
    display: inline-block;
    width: 18rem;

    .v-input {
        background-color: transparent;
        border-left: none;
        border-top: 0.25rem solid transparent;
        height: calc(var(--input-height) + 1rem);
        line-height: calc(var(--input-height) + 1rem);
        padding: 0.25rem 1rem 0.5rem;
        width: 100%;
    }

    .v-datetime-picker {
        border-spacing: 0;
        height: 18rem;
        margin: 0 0.5rem 0.5rem;
        table-layout: fixed;
        width: calc(100% - 1rem);
    }

    .v-datetime-title,
    .v-datetime-header .v-icon,
    .v-datetime-body .v-datetime-cell {
        cursor: pointer;
    }

    .v-datetime-header {
        .v-datetime-row:first-child {
            height: 2.25rem;

            th {
                background-color: var(--dark-gray);

                &:first-child {
                    border-radius: 0.2rem 0 0 0.2rem;
                }

                &:last-child {
                    border-radius: 0 0.2rem 0.2rem 0;
                }
            }
        }

        .v-datetime-row:last-child {
            color: var(--light-gray);
            padding: 0.25rem 0;

            &.placeholder {
                height: 0.5rem;
            }
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

    .v-datetime-cell {
        border-radius: 0.2rem;
        position: relative;
        text-align: center;

        &.future,
        &.past {
            color: var(--gray);
        }

        &[aria-disabled="true"] {
            color: var(--dark-gray);
            pointer-events: none;

            &.today {
                color: var(--background);
            }
        }

        &.current {
            background-color: var(--accent) !important;
        }

        &.today {
            background-color: var(--dark-gray);
        }
    }

    .v-datetime-body .v-datetime-cell {
        &:active,
        &:focus,
        &:hover {
            background-color: var(--gray);
            color: var(--color);
        }
    }
}
</style>
