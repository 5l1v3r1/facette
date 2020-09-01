/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {DateTime, DateTimeOptions} from "luxon";

import store from "@/store";

export function formatDate(input: string, format: string, timezone = true): string {
    if (store.state.timezoneUTC) {
        let value = parseDate(input).toFormat(format);
        if (timezone) {
            value += " UTC";
        }

        return value;
    }

    return parseDate(input).toFormat(format);
}

export function parseDate(input?: string, format?: string): DateTime {
    if (!input) {
        return store.state.timezoneUTC ? DateTime.utc() : DateTime.local();
    }

    const opts: DateTimeOptions | undefined = store.state.timezoneUTC ? {zone: "Etc/UTC"} : undefined;

    if (format) {
        return DateTime.fromFormat(input, format, opts);
    }

    return DateTime.fromISO(input, opts);
}

export function weekdaysShort(): Array<string> {
    const out: Array<string> = [];

    let day = parseDate().startOf("week");
    for (let i = 0; i < 7; i++) {
        out.push(day.weekdayShort);
        day = day.plus({days: 1});
    }

    return out;
}
