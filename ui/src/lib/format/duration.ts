/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {formatNumber} from ".";

export enum Unit {
    Year = "y",
    Month = "M",
    Week = "w",
    Day = "d",
    Hour = "h",
    Minute = "m",
    Second = "s",
    Millisecond = "ms",
    Microsecond = "µs",
    Nanosecond = "ns",
}

const units: Record<Unit, number> = {
    [Unit.Year]: 31557600, // 86400 * 365.25
    [Unit.Month]: 2630016, // 86400 * 30.44
    [Unit.Week]: 604800,
    [Unit.Day]: 86400,
    [Unit.Hour]: 3600,
    [Unit.Minute]: 60,
    [Unit.Second]: 1,
    [Unit.Millisecond]: 1e-3,
    [Unit.Microsecond]: 1e-6,
    [Unit.Nanosecond]: 1e-9,
};

export default function (input: number | null, decimals: number, base: Unit = Unit.Second): string {
    if (input === null) {
        return "null";
    } else if (input === 0) {
        return "0";
    }

    const unit: number | undefined = units[base];
    if (units === undefined) {
        throw Error("unsupported unit base");
    }

    if (base !== Unit.Second) {
        input *= unit;
    }

    for (const key in units) {
        const value: number = input / units[key as Unit];
        if (value >= 1 || key === Unit.Nanosecond) {
            return `${formatNumber(value, decimals)} ${key}`;
        }
    }

    return formatNumber(input, decimals);
}
