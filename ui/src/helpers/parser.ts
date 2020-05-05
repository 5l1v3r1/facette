/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

const alpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";

const digits = "0123456789";

export function formatPosition(pos: Position): string {
    return `${pos.line}:${pos.char}`;
}

function isDigit(c: string): boolean {
    return digits.includes(c);
}

function isIdentChar(c: string): boolean {
    return (alpha + digits + "_").includes(c);
}

function isSpace(c: string): boolean {
    return " \t\n".includes(c);
}

export interface Token {
    type: TokenType;
    pos: Position;
    text: string;
}

export enum TokenType {
    INVALID = "invalid",
    EOF = "end of input",

    IDENT = "ident", // ident
    NUMBER = "number", // 123.45
    STRING = "string", // "abc" or 'abc'
    NEWLINE = "new line", // \n
    BADESCAPE = "bad escape", // \b

    EQ = "equal", // =
    NEQ = "not equal", // !=
    EQREGEXP = "equal pattern", // =~
    NEQREGEXP = "not equal pattern", // !~

    LBRACE = "left brace", // {
    RBRACE = "right brace", // }
    LPAREN = "left parenthesis", // (
    RPAREN = "right parenthesis", // )
    COMMA = "comma", // ,
}

export interface Position {
    line: number;
    char: number;
}

export class Parser {
    private last: Position | null = null;

    private lastCh: string | null = null;

    private peeked: Token | null = null;

    private pos: Position = {line: 1, char: 1};

    private text = "";

    public constructor(text: string) {
        this.text = text;
    }

    public expect(type: string): Token {
        const tok = this.next();
        if (tok.type !== type) {
            throw Error(`expected ${type} but got ${tok.type} at ${formatPosition(tok.pos)}`);
        }

        return tok;
    }

    public next(): Token {
        let tok: Token;

        if (this.peeked !== null) {
            tok = this.peeked;
            this.peeked = null;
        } else {
            tok = this.scan();
        }

        return tok;
    }

    public peek(): Token {
        if (this.peeked === null) {
            this.peeked = this.scan();
        }

        return this.peeked;
    }

    private read(): string {
        const c: string = this.text.slice(0, 1);
        this.text = this.text.slice(1);

        this.last = Object.assign({}, this.pos);
        this.lastCh = c;

        if (c === "\n") {
            this.pos.line++;
            this.pos.char = 1;
        } else {
            this.pos.char++;
        }

        return c;
    }

    private run(set: string): string {
        let s = "";

        for (;;) {
            const c = this.read();
            if (!set.includes(c)) {
                break;
            }

            s += c;
        }

        return s;
    }

    private scan(): Token {
        let c = this.read();
        while (isSpace(c)) {
            c = this.read();
        }

        if (c === "") {
            return this.tokenAtPos(TokenType.EOF);
        } else if (c === "-" || isDigit(c)) {
            this.unread();
            return this.scanNumber();
        } else if (isIdentChar(c)) {
            this.unread();
            return this.scanIdent();
        } else if (c === '"' || c === "'") {
            this.unread();
            return this.scanString();
        } else if (c === "=") {
            const next = this.read();
            if (next === "~") {
                return this.tokenAtPos(TokenType.EQREGEXP);
            }

            this.unread();

            return this.tokenAtPos(TokenType.EQ);
        } else if (c === "!") {
            const next = this.read();
            if (next === "=") {
                return this.tokenAtPos(TokenType.NEQ);
            } else if (next === "~") {
                return this.tokenAtPos(TokenType.NEQREGEXP);
            }

            this.unread();
        } else if (c === "{") {
            return this.tokenAtPos(TokenType.LBRACE);
        } else if (c === "}") {
            return this.tokenAtPos(TokenType.RBRACE);
        } else if (c === "(") {
            return this.tokenAtPos(TokenType.LPAREN);
        } else if (c === ")") {
            return this.tokenAtPos(TokenType.RPAREN);
        } else if (c === ",") {
            return this.tokenAtPos(TokenType.COMMA);
        }

        return this.tokenAtPos(TokenType.INVALID, c);
    }

    private scanIdent(): Token {
        let s = "";

        for (;;) {
            const c = this.read();
            if (c === "") {
                break;
            } else if (isIdentChar(c)) {
                s += c;
            } else {
                this.unread();
                break;
            }
        }

        return this.tokenAtPos(TokenType.IDENT, s);
    }

    private scanNumber(): Token {
        let s = "";

        let next = this.read();
        if (next === "-") {
            s += next;
        } else {
            this.unread();
        }

        s += this.run(digits);

        next = this.read();
        if (next === ".") {
            s += next + this.run(digits);
        } else {
            this.unread();
        }

        next = this.read();
        if ("eE".includes(next)) {
            s += next;

            next = this.read();
            if ("+-".includes(next)) {
                s += next;
            } else {
                this.unread();
            }

            s += this.run(digits);
        } else {
            this.unread();
        }

        return this.tokenAtPos(TokenType.NUMBER, s);
    }

    private scanString(): Token {
        const quote = this.read();

        let s = "";

        for (;;) {
            const c = this.read();

            switch (c) {
                case quote:
                    return this.tokenAtPos(TokenType.STRING, s);

                case "":
                    return this.tokenAtPos(TokenType.EOF, s);

                case "\n":
                    return this.tokenAtPos(TokenType.NEWLINE, s);

                case "\\": {
                    const next = this.read();
                    switch (next) {
                        case "\\":
                        case '"':
                        case "'":
                            s += next;
                            break;

                        default:
                            return this.tokenAtPos(TokenType.BADESCAPE, c + next);
                    }

                    break;
                }

                default:
                    s += c;
            }
        }
    }

    private tokenAtPos(type: TokenType, text = ""): Token {
        return {type, pos: Object.assign({}, this.pos), text};
    }

    private unread(): void {
        if (this.last !== null && this.lastCh !== null) {
            this.text = this.lastCh + this.text;
            this.pos = this.last;
            this.last = null;
            this.lastCh = null;
        }
    }
}
