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
    constants?: Array<number>;
    stack?: StackMode;
    title?: string;
    type?: ChartType;
    unit?: Unit;
    variables?: Array<TemplateVariable>;
}

declare type StackMode = "" | "normal" | "percent";

declare type ChartType = "area" | "bar" | "line";

declare interface ChartSeries {
    expr: string;
    options?: ChartSeriesOptions;
}

declare interface ChartSeriesOptions {
    color?: string;
}

declare interface ChartAxes {
    yLeft: ChartAxis;
    yRight: ChartAxis;
}

declare interface ChartAxis {
    center?: boolean;
    label?: string;
    max?: number;
    min?: number;
}

declare interface Unit {
    type?: UnitType;
    base?: string;
}

declare type UnitType = "" | "binary" | "count" | "duration" | "metric";

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

declare interface TemplateVariable {
    name: string;
    value?: string;
    label?: string;
    filter?: string;
    dynamic: boolean;
}
