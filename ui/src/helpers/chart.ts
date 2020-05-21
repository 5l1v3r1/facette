/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import cloneDeep from "lodash/cloneDeep";

import {parseVariables, renderTemplate} from "@/src/helpers/template";

export function cleanupChart(chart: Chart): Chart {
    if (!chart.options?.axes?.y?.left?.unit?.type) {
        delete chart.options?.axes?.y?.left?.unit;
    }

    if (!chart.options?.axes?.y?.right?.unit?.type) {
        delete chart.options?.axes?.y?.right?.unit;
    }

    if (!chart.options?.title) {
        delete chart.options?.title;
    }

    return chart;
}

export function parseChartVariables(chart: Chart): Array<TemplateVariable> {
    let data = "";

    chart.series?.forEach(series => {
        data += `\xff${series.expr}`;
    });

    if (chart.options?.axes?.y?.left?.label) {
        data += `\xff${chart.options.axes.y.left.label}`;
    }

    if (chart.options?.axes?.y?.right?.label) {
        data += `\xff${chart.options.axes.y.right.label}`;
    }

    if (chart.options?.title) {
        data += `\xff${chart.options.title}`;
    }

    return parseVariables(data).map(name => ({name, dynamic: false}));
}

export function renderChart(chart: Chart, data: Record<string, string>): Chart {
    const proxy: Chart = cloneDeep(chart);

    proxy.series?.forEach(series => {
        series.expr = renderTemplate(series.expr, data);
    });

    if (proxy.options?.axes?.y?.left?.label) {
        proxy.options.axes.y.left.label = renderTemplate(proxy.options.axes.y.left.label, data);
    }

    if (proxy.options?.axes?.y?.right?.label) {
        proxy.options.axes.y.right.label = renderTemplate(proxy.options.axes.y.right.label, data);
    }

    if (proxy.options?.title) {
        proxy.options.title = renderTemplate(proxy.options.title, data);
    }

    return proxy;
}
