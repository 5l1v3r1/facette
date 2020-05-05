/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

declare interface APIResponse<T> {
    data?: T;
    total?: number;
    error?: string;
}

declare interface BulkResult {
    status: number;
    headers: Record<string, string>;
    response: APIResponse<unknown>;
}

declare type Labels = Record<string, string>;

declare interface ListParams {
    limit?: number;
    offset?: number;
    sort?: string;
    [key: string]: unknown;
}

declare interface Point {
    0: number;
    1: number | null;
}

declare interface SeriesQuery {
    from?: string;
    to?: string;
    step?: number | string;
    exprs: Array<string>;
}

declare interface SeriesResult {
    from: string;
    to: string;
    step: number;
    series: Array<Series>;
}

declare interface Series {
    name: string;
    points: Array<Point>;
    summary: SeriesSummary;
}

declare type SeriesSummary = Record<string, number>;

declare interface TimeRange {
    from?: string;
    to?: string;
}
