/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {expect} from "chai";

import format, {Unit} from "./binary";

describe("formatBinary", () => {
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

    describe("1 KiB", () => {
        const actual: string = format(1, 0, Unit.Kibi);
        expect(actual).to.equal("1 Ki");
    });

    describe("1 Ki from base", () => {
        const actual: string = format(1024, 0);
        expect(actual).to.equal("1 Ki");
    });

    describe("1 Mi", () => {
        const actual: string = format(1, 0, Unit.Mebi);
        expect(actual).to.equal("1 Mi");
    });

    describe("1 Gi", () => {
        const actual: string = format(1, 0, Unit.Gibi);
        expect(actual).to.equal("1 Gi");
    });

    describe("1 Ti", () => {
        const actual: string = format(1, 0, Unit.Tebi);
        expect(actual).to.equal("1 Ti");
    });

    describe("1 Pi", () => {
        const actual: string = format(1, 0, Unit.Pebi);
        expect(actual).to.equal("1 Pi");
    });

    describe("1 Ei", () => {
        const actual: string = format(1, 0, Unit.Exbi);
        expect(actual).to.equal("1 Ei");
    });

    describe("Ki without decimals", () => {
        const actual: string = format(1234, 0);
        expect(actual).to.equal("1 Ki");
    });

    describe("Ki with 1 decimal", () => {
        const actual: string = format(1234, 1);
        expect(actual).to.equal("1.2 Ki");
    });

    describe("Ki with 2 decimals", () => {
        const actual: string = format(1234, 2);
        expect(actual).to.equal("1.21 Ki");
    });

    describe("Ki with 3 decimals", () => {
        const actual: string = format(1234, 3);
        expect(actual).to.equal("1.205 Ki");
    });
});
