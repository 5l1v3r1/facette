/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {formatNumber} from ".";

export enum Unit {
    Exbi = "Ei",
    Pebi = "Pi",
    Tebi = "Ti",
    Gibi = "Gi",
    Mebi = "Mi",
    Kibi = "Ki",
    Base = "",
}

const units: Record<Unit, number> = {
    [Unit.Exbi]: Math.pow(1024, 6),
    [Unit.Pebi]: Math.pow(1024, 5),
    [Unit.Tebi]: Math.pow(1024, 4),
    [Unit.Gibi]: Math.pow(1024, 3),
    [Unit.Mebi]: Math.pow(1024, 2),
    [Unit.Kibi]: 1024,
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
