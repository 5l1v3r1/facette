/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {expect} from "chai";

import format, {Unit} from "./count";

describe("formatMetric", () => {
    describe("null", () => {
        const actual: string = format(null, 0);
        expect(actual).to.equal("null");
    });

    describe("0", () => {
        const actual: string = format(0, 0);
        expect(actual).to.equal("0");
    });

    describe("1 from base", () => {
        const actual: string = format(1, 0);
        expect(actual).to.equal("1");
    });

    describe("1 from thousand", () => {
        const actual: string = format(0.001, 0, Unit.Thousand);
        expect(actual).to.equal("1");
    });

    describe("1 K", () => {
        const actual: string = format(1, 0, Unit.Thousand);
        expect(actual).to.equal("1 K");
    });

    describe("1 K from base", () => {
        const actual: string = format(1e3, 0);
        expect(actual).to.equal("1 K");
    });

    describe("1 M", () => {
        const actual: string = format(1, 0, Unit.Million);
        expect(actual).to.equal("1 M");
    });

    describe("1 B", () => {
        const actual: string = format(1, 0, Unit.Billion);
        expect(actual).to.equal("1 B");
    });

    describe("1 t", () => {
        const actual: string = format(1, 0, Unit.Trillion);
        expect(actual).to.equal("1 t");
    });

    describe("K without decimals", () => {
        const actual: string = format(1234, 0);
        expect(actual).to.equal("1 K");
    });

    describe("K with 1 decimal", () => {
        const actual: string = format(1234, 1);
        expect(actual).to.equal("1.2 K");
    });

    describe("K with 2 decimals", () => {
        const actual: string = format(1234, 2);
        expect(actual).to.equal("1.23 K");
    });

    describe("K with 3 decimals", () => {
        const actual: string = format(1234, 3);
        expect(actual).to.equal("1.234 K");
    });

    describe("K with rounded decimals", () => {
        const actual: string = format(1235, 2);
        expect(actual).to.equal("1.24 K");
    });
});
