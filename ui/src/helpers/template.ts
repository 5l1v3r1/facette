/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

interface TemplateNode {
    type: "text" | "variable";
    value: string;
}

function isVariableChar(c: string): boolean {
    return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_".includes(c);
}

class Template {
    private s = "";

    private nodes: Array<TemplateNode> = [];

    public parse(text: string): void {
        this.s = text;
        this.nodes = [];

        while (this.s.length > 0) {
            const next = this.peek(2);
            if (next.length < 2) {
                this.nodes.push({type: "text", value: next});
                break;
            }

            if (next[0] === "$" && (isVariableChar(next[1]) || next[1] === "{")) {
                this.s = this.s.substr(1);
                this.nodes.push(this.readVariable());
            }

            this.nodes.push(this.readText());
        }
    }

    public render(data: Record<string, string>): string {
        return this.nodes
            .map((node: TemplateNode) => {
                switch (node.type) {
                    case "text":
                        return node.value;

                    case "variable":
                        return data[node.value] || "";
                }
            })
            .join("");
    }

    public variables(): Array<string> {
        return this.nodes.reduce((vars: Array<string>, node: TemplateNode) => {
            if (node.type === "variable" && !vars.includes(node.value)) {
                vars.push(node.value);
            }
            return vars;
        }, []);
    }

    private discard(): void {
        this.s = this.s.substr(1);
    }

    private peek(n: number): string {
        return this.s.substr(0, n);
    }

    private read(): string {
        const c: string = this.s.substr(0, 1);
        this.s = this.s.substr(1);
        return c;
    }

    private readText(): TemplateNode {
        let prev = "";
        let value = "";

        for (;;) {
            const c = this.read();
            if (c === "") {
                break;
            }

            if (c === "$" && prev !== "\\") {
                this.unread(c);

                const next = this.peek(2);
                if (isVariableChar(next[1]) || next[1] === "{") {
                    break;
                }

                this.discard();
            }

            value += c;
            prev = c;
        }

        return {type: "text", value};
    }

    private readVariable(): TemplateNode {
        let lbrace = false;
        let rbrace = false;

        const next = this.peek(1);
        if (next === "{") {
            lbrace = true;
            this.discard();
        }

        let value = "";

        for (;;) {
            const c = this.read();
            if (c === "") {
                break;
            }

            if (lbrace && c === "}") {
                rbrace = true;
                break;
            } else if (isVariableChar(c)) {
                value += c;
            } else {
                this.unread(c);
                break;
            }
        }

        if (lbrace && !rbrace) {
            throw Error("unbalanced brace");
        }

        return {type: "variable", value};
    }

    private unread(c: string): void {
        this.s = c + this.s;
    }
}

export function renderTemplate(text: string, data: Record<string, string>): string {
    const tmpl: Template = new Template();
    tmpl.parse(text);
    return tmpl.render(data);
}

export function parseVariables(text: string): Array<string> {
    const tmpl: Template = new Template();
    tmpl.parse(text);
    return tmpl.variables();
}
