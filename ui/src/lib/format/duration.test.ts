/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {expect} from "chai";

import format, {Unit} from "./duration";

describe("formatDuration", () => {
    describe("null", () => {
        const actual: string = format(null, 0);
        expect(actual).to.equal("null");
    });

    describe("0", () => {
        const actual: string = format(0, 0);
        expect(actual).to.equal("0");
    });

    describe("1 ns", () => {
        const actual: string = format(1, 0, Unit.Nanosecond);
        expect(actual).to.equal("1 ns");
    });

    describe("1 ns from microseconds", () => {
        const actual: string = format(0.001, 0, Unit.Microsecond);
        expect(actual).to.equal("1 ns");
    });

    describe("1 µs", () => {
        const actual: string = format(1, 0, Unit.Microsecond);
        expect(actual).to.equal("1 µs");
    });

    describe("1 µs from nanoseconds", () => {
        const actual: string = format(1e3, 0, Unit.Nanosecond);
        expect(actual).to.equal("1 µs");
    });

    describe("1 ms", () => {
        const actual: string = format(1, 0, Unit.Millisecond);
        expect(actual).to.equal("1 ms");
    });

    describe("1 s", () => {
        const actual: string = format(1, 0, Unit.Second);
        expect(actual).to.equal("1 s");
    });

    describe("1 s from base", () => {
        const actual: string = format(1, 0);
        expect(actual).to.equal("1 s");
    });

    describe("1 m", () => {
        const actual: string = format(1, 0, Unit.Minute);
        expect(actual).to.equal("1 m");
    });

    describe("1 h", () => {
        const actual: string = format(1, 0, Unit.Hour);
        expect(actual).to.equal("1 h");
    });

    describe("1 d", () => {
        const actual: string = format(1, 0, Unit.Day);
        expect(actual).to.equal("1 d");
    });

    describe("1 w", () => {
        const actual: string = format(1, 0, Unit.Week);
        expect(actual).to.equal("1 w");
    });

    describe("1 M", () => {
        const actual: string = format(1, 0, Unit.Month);
        expect(actual).to.equal("1 M");
    });

    describe("1 y", () => {
        const actual: string = format(1, 0, Unit.Year);
        expect(actual).to.equal("1 y");
    });

    describe("1.23 y", () => {
        const actual: string = format(38789280, 2);
        expect(actual).to.equal("1.23 y");
    });

    describe("1.23 y from minutes", () => {
        const actual: string = format(646488, 2, Unit.Minute);
        expect(actual).to.equal("1.23 y");
    });

    describe("s without decimals", () => {
        const actual: string = format(1234, 0, Unit.Millisecond);
        expect(actual).to.equal("1 s");
    });

    describe("s with 1 decimal", () => {
        const actual: string = format(1234, 1, Unit.Millisecond);
        expect(actual).to.equal("1.2 s");
    });

    describe("s with 2 decimals", () => {
        const actual: string = format(1230, 2, Unit.Millisecond);
        expect(actual).to.equal("1.23 s");
    });

    describe("s with 3 decimals", () => {
        const actual: string = format(1230, 3, Unit.Millisecond);
        expect(actual).to.equal("1.23 s");
    });

    describe("s with rounded decimals", () => {
        const actual: string = format(1235, 2, Unit.Millisecond);
        expect(actual).to.equal("1.24 s");
    });
});
