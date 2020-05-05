/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {Parser, TokenType} from "./parser";

export const NameLabel = "__name__";

export function labelsToString(labels: Labels): string {
    let s = "";

    const ls = Object.assign({}, labels);

    if (ls[NameLabel]) {
        s += ls[NameLabel];
        delete ls[NameLabel];
    }

    // Remove internal-only hash reference
    delete ls._hash;

    if (ls) {
        s += `{${Object.keys(ls).map((key: string) => `${key}=${JSON.stringify(ls[key])}`)}}`;
    }

    return s;
}

export function parseLabels(text: string): Labels {
    const labels: Labels = {};

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
