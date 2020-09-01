/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {Parser, TokenType} from "@/lib/parser";

export const NameLabel = "__name__";

export class Labels {
    private value: Record<string, string>;

    constructor(input?: string) {
        this.value = input ? parseLabels(input) : {};
    }

    public entries(name = true): Record<string, string> {
        if (name) {
            return Object.assign({}, this.value);
        }

        const obj = Object.assign({}, this.value);

        return Object.keys(obj).reduce((out: Record<string, string>, key: string) => {
            if (key !== NameLabel) {
                out[key] = obj[key];
            }
            return out;
        }, {});
    }

    public name(): string | null {
        return this.value[NameLabel] ?? null;
    }

    public toString(): string {
        let s = this.name() ?? "";

        const ls = this.entries(false);
        if (ls) {
            s += `{${Object.keys(ls).map(key => `${key}=${JSON.stringify(ls[key])}`)}}`;
        }

        return s;
    }
}

function parseLabels(text: string): Record<string, string> {
    const labels: Record<string, string> = {};

    const parser = new Parser(text);

    let tok = parser.peek();
    const isName = tok.type === TokenType.IDENT;

    if (isName) {
        labels[NameLabel] = tok.text;
        parser.next();
    }

    // Stop if not a left brace and name was at least provided
    tok = parser.peek();
    if (tok.type !== TokenType.LBRACE && isName) {
        return labels;
    }

    tok = parser.expect(TokenType.LBRACE);

    if (parser.peek().type !== TokenType.RBRACE) {
        for (;;) {
            tok = parser.expect(TokenType.IDENT);
            const name = tok.text;

            parser.expect(TokenType.EQ);

            tok = parser.expect(TokenType.STRING);
            labels[name] = tok.text;

            // Allow extraneous comma
            if (parser.peek().type !== TokenType.COMMA) {
                break;
            }

            parser.next();

            if (parser.peek().type === TokenType.RBRACE) {
                break;
            }
        }
    }

    parser.expect(TokenType.RBRACE);

    return labels;
}
