<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-message type="placeholder">
        <v-icon icon="dizzy"></v-icon>

        {{ message }}

        <v-button @click="retry">
            {{ i18n.t("labels.retry") }}
        </v-button>
    </v-message>
</template>

<script lang="ts">
import {SetupContext, computed} from "vue";
import {useI18n} from "vue-i18n";

import common from "@/common";

export default {
    props: {
        type: {
            default: null,
            type: String,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const i18n = useI18n();

        const {error, resetError} = common;

        const message = computed((): string => {
            switch (error.value) {
                case "notFound":
                    return i18n.t(`messages.${props.type ?? "error"}.notFound`);
            }

            return i18n.t("messages.error.unhandled");
        });

        const retry = (): void => {
            resetError();
            ctx.emit("retry");
        };

        return {
            i18n,
            message,
            retry,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-message {
    color: var(--gray);
    bottom: 0;
    flex-direction: column;
    font-size: 1rem;
    justify-content: center;
    left: 0;
    position: absolute;
    right: 0;
    top: 0;

    .v-icon {
        color: var(--dark-gray);
        font-size: 4rem;
        margin: 0 0 1rem;
    }

    .v-button {
        font-size: 0.85rem;
        margin-top: 2rem;
    }
}
</style>
