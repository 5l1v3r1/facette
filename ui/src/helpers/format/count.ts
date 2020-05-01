/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {formatNumber} from ".";

export enum Unit {
    Trillion = "t",
    Billion = "B",
    Million = "M",
    Thousand = "K",
    Base = "",
}

const units: Record<Unit, number> = {
    [Unit.Trillion]: 1e12,
    [Unit.Billion]: 1e9,
    [Unit.Million]: 1e6,
    [Unit.Thousand]: 1e3,
    [Unit.Base]: 1,
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
        if (value >= 1 && key !== Unit.Base) {
            return `${formatNumber(value, decimals)} ${key}`;
        }
    }

    return formatNumber(input, decimals);
}
