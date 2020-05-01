/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

export {default as formatBinary, Unit as BinaryUnit} from "./binary";
export {default as formatCount, Unit as CountUnit} from "./count";
export {default as formatDuration, Unit as DurationUnit} from "./duration";
export {default as formatMetric, Unit as MetricUnit} from "./metric";

export function formatNumber(value: number, decimals: number): string {
    const factor = decimals > 0 ? Math.pow(10, decimals) : 1;
    return (Math.round(value * factor) / factor).toString();
}
