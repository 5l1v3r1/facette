/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {expect} from "chai";

import format, {Unit} from "./metric";

describe("formatMetric", () => {
    describe("null", () => {
        const actual: string = format(null, 0);
        expect(actual).to.equal("null");
    });

    describe("0", () => {
        const actual: string = format(0, 0);
        expect(actual).to.equal("0");
    });

    describe("1 a", () => {
        const actual: string = format(1, 0, Unit.Atto);
        expect(actual).to.equal("1 a");
    });

    describe("1 f", () => {
        const actual: string = format(1, 0, Unit.Femto);
        expect(actual).to.equal("1 f");
    });

    describe("1 p", () => {
        const actual: string = format(1, 0, Unit.Pico);
        expect(actual).to.equal("1 p");
    });

    describe("1 n", () => {
        const actual: string = format(1, 0, Unit.Nano);
        expect(actual).to.equal("1 n");
    });

    describe("1 µ", () => {
        const actual: string = format(1, 0, Unit.Micro);
        expect(actual).to.equal("1 µ");
    });

    describe("1 m", () => {
        const actual: string = format(1, 0, Unit.Milli);
        expect(actual).to.equal("1 m");
    });

    describe("1 from base", () => {
        const actual: string = format(1, 0);
        expect(actual).to.equal("1");
    });

    describe("1 from kilo", () => {
        const actual: string = format(0.001, 0, Unit.Kilo);
        expect(actual).to.equal("1");
    });

    describe("1 from milli", () => {
        const actual: string = format(1e3, 0, Unit.Milli);
        expect(actual).to.equal("1");
    });

    describe("1 k", () => {
        const actual: string = format(1, 0, Unit.Kilo);
        expect(actual).to.equal("1 k");
    });

    describe("1 M", () => {
        const actual: string = format(1, 0, Unit.Mega);
        expect(actual).to.equal("1 M");
    });

    describe("1 G", () => {
        const actual: string = format(1, 0, Unit.Giga);
        expect(actual).to.equal("1 G");
    });

    describe("1 T", () => {
        const actual: string = format(1, 0, Unit.Tera);
        expect(actual).to.equal("1 T");
    });

    describe("1 P", () => {
        const actual: string = format(1, 0, Unit.Peta);
        expect(actual).to.equal("1 P");
    });

    describe("1 E", () => {
        const actual: string = format(1, 0, Unit.Exa);
        expect(actual).to.equal("1 E");
    });

    describe("k without decimals", () => {
        const actual: string = format(1234, 0);
        expect(actual).to.equal("1 k");
    });

    describe("k with 1 decimal", () => {
        const actual: string = format(1234, 1);
        expect(actual).to.equal("1.2 k");
    });

    describe("k with 2 decimals", () => {
        const actual: string = format(1234, 2);
        expect(actual).to.equal("1.23 k");
    });

    describe("k with 3 decimals", () => {
        const actual: string = format(1234, 3);
        expect(actual).to.equal("1.234 k");
    });

    describe("k with rounded decimals", () => {
        const actual: string = format(1235, 2);
        expect(actual).to.equal("1.24 k");
    });
});
