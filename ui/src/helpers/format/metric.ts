/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {formatNumber} from ".";

export enum Unit {
    Exa = "E",
    Peta = "P",
    Tera = "T",
    Giga = "G",
    Mega = "M",
    Kilo = "k",
    Base = "",
    Milli = "m",
    Micro = "µ",
    Nano = "n",
    Pico = "p",
    Femto = "f",
    Atto = "a",
}

const units: Record<Unit, number> = {
    [Unit.Exa]: 1e18,
    [Unit.Peta]: 1e15,
    [Unit.Tera]: 1e12,
    [Unit.Giga]: 1e9,
    [Unit.Mega]: 1e6,
    [Unit.Kilo]: 1e3,
    [Unit.Base]: 1,
    [Unit.Milli]: 1e-3,
    [Unit.Micro]: 1e-6,
    [Unit.Nano]: 1e-9,
    [Unit.Pico]: 1e-12,
    [Unit.Femto]: 1e-15,
    [Unit.Atto]: 1e-18,
};

export default function (input: number | null, decimals: number, base: Unit = Unit.Base): string {
    if (input === null) {
        return "null";
    } else if (input === 0) {
        return "0";
    }

    const unit: number | undefined = units[base];
    if (units === undefined) {
        throw Error("unsupported unit base");
    }

    if (base !== Unit.Base) {
        input *= unit;
    }

    for (const key in units) {
        const value: number = input / units[key as Unit];
        if (value >= 1 || key === Unit.Atto) {
            return key === Unit.Base ? formatNumber(value, decimals) : `${formatNumber(value, decimals)} ${key}`;
        }
    }

    return formatNumber(input, decimals);
}
