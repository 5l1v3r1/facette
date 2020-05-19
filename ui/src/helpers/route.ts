/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {Route} from "vue-router";

import {ModalConfirmParams} from "@/src/components/modal/confirm.vue";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function beforeRoute(this: Vue, to: Route, from: Route, next: any): void {
    if (to.path === from.path || !this.$store.getters.routeGuarded) {
        next();
        return;
    }

    this.$components.modal(
        "confirm",
        {
            button: {
                label: this.$t("labels.leavePage"),
                danger: true,
            },
            message: this.$t("messages.confirmLeave"),
        } as ModalConfirmParams,
        (ok: boolean) => {
            if (ok) {
                const fn: ((e: Event) => void) | null = this.$store.getters.routeGuardFn;
                if (fn) {
                    window.removeEventListener("beforeunload", fn);
                }
                next();
            }
        },
    );

    next(false);
}
