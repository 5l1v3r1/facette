/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

declare interface ObjectMeta {
    id: string;
    name: string;
    created?: Date;
    modified?: Date;
}

declare interface Chart extends ObjectMeta {
    options?: ChartOptions;
    series?: Array<ChartSeries>;
    link?: string;
    template?: boolean;
}

declare interface ChartOptions {
    axes?: ChartAxes;
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

declare interface Dashboard extends ObjectMeta {
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

declare type DashboardItemType = "chart";

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

declare interface Provider extends ObjectMeta {
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
