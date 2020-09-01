/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

declare type APIError = "notFound" | "unhandled" | null;

declare interface APIResponse<T> {
    data?: T;
    total?: number;
    error?: string;
}

declare interface BulkRequest {
    endpoint: string;
    method: string;
    params?: Record<string, unknown>;
    data?: unknown;
}

declare interface BulkResult {
    status: number;
    headers: Record<string, string>;
    response: APIResponse<unknown>;
}

declare interface ListParams {
    limit?: number;
    offset?: number;
    sort?: string;
    [key: string]: unknown;
}

declare interface Labels {
    entries(name?: boolean): Record<string, string>;
    name(): string | null;
    toString(): string;
}

declare interface Options {
    connectors: Array<string>;
    driver: {
        name: string;
        version: string;
    };
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

declare interface TestResult {
    success: boolean;
}

declare interface TimeRange {
    from?: string;
    to?: string;
}

declare interface Version {
    version?: string;
    branch?: string;
    revision?: string;
    compiler?: string;
    buildDate?: string;
}

// Objects
declare interface ObjectBase {
    id: string;
    name: string;
    createdAt?: Date;
    modifiedAt?: Date;
}

declare interface Chart extends ObjectBase {
    options?: ChartOptions;
    series?: Array<ChartSeries>;
    link?: string;
    template?: boolean;
}

declare interface ChartOptions {
    axes?: ChartAxes;
    legend?: boolean;
    markers?: Array<Marker>;
    title?: string;
    type?: ChartType;
    variables?: Array<TemplateVariable>;
}

declare interface ChartAxes {
    x?: ChartXAxis;
    y?: ChartYAxes;
}

declare interface ChartXAxis {
    show?: boolean;
}

declare interface ChartYAxes {
    center?: boolean;
    left?: ChartYAxis;
    right?: ChartYAxis;
    stack?: StackMode;
}

declare interface ChartYAxis {
    show?: boolean;
    label?: string;
    max?: number;
    min?: number;
    unit?: Unit;
}

declare interface Marker {
    color?: string;
    label?: string;
    value: number;
    axis?: "left" | "right";
}

declare type StackMode = "" | "normal" | "percent";

declare interface Unit {
    type?: UnitType;
    base?: string;
}

declare type UnitType = "" | "binary" | "count" | "duration" | "metric";

declare type ChartType = "area" | "bar" | "line";

declare interface ChartSeries {
    expr: string;
    options?: ChartSeriesOptions;
}

declare interface ChartSeriesOptions {
    color?: string;
    axis?: "left" | "right";
}

declare interface Dashboard extends ObjectBase {
    options?: DashboardOptions;
    layout?: GridLayout;
    items?: Array<DashboardItem>;
    parent?: string;
    link?: string;
    template?: boolean;
    references?: Array<Reference>;
}

declare interface DashboardOptions {
    title?: string;
    variables?: Array<TemplateVariable>;
}

declare interface DashboardItem {
    type: DashboardItemType;
    layout: GridItemLayout;
    options?: Record<string, unknown>;
}

declare type DashboardItemType = "chart" | "text";

declare interface GridLayout {
    columns: number;
    rowHeight: number;
    rows: number;
}

declare interface GridItemLayout {
    x: number;
    y: number;
    w: number;
    h: number;
}

declare interface Provider extends ObjectBase {
    connector: ProviderConnector;
    filters?: Array<FilterRule>;
    pollInterval?: number;
    enabled?: boolean;
    error?: string;
}

declare interface ProviderConnector {
    type: string;
    settings: Record<string, unknown>;
}

declare interface FilterRule {
    action: FilterAction;
    label: string;
    pattern: string;
    into?: string;
    targets?: Record<string, string>;
}

declare type FilterAction = "discard" | "relabel" | "rewrite" | "sieve";

declare interface Reference {
    type: string;
    value: unknown;
}

declare interface TemplateVariable {
    name: string;
    value?: string;
    label?: string;
    filter?: string;
    dynamic: boolean;
}
