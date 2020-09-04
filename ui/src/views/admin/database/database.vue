<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-content>
        <v-toolbar clip="content">
            <v-spacer></v-spacer>

            <v-label class="note">{{ i18n.t("labels.database.driver", apiOptions.driver) }}</v-label>
        </v-toolbar>

        <v-form class="third">
            <h1>{{ i18n.t("labels.database._") }}</h1>

            <v-markdown :content="i18n.t('help.database.archive')"></v-markdown>

            <v-message type="warning">{{ i18n.t("help.database.restoreWarning") }}</v-message>

            <v-button icon="download" primary @click="dump">
                {{ i18n.t("labels.database.downloadArchive") }}
            </v-button>

            <v-button icon="upload" @click="$refs.file.click()">
                {{ i18n.t("labels.database.restoreArchive") }}
            </v-button>

            <input accept="application/gzip" ref="file" type="file" @change="restore" />
        </v-form>
    </v-content>
</template>

<script lang="ts">
import {computed, onMounted, ref} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";

import {useUI} from "@/components/ui";
import {prefix as apiPrefix, onFetch} from "@/lib/api";
import {State} from "@/store";
import {ModalConfirmParams} from "@/components/modal/confirm.vue";

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();
        const store = useStore<State>();
        const ui = useUI();

        const file = ref<HTMLInputElement | null>(null);

        const apiOptions = computed(() => store.state.apiOptions);

        const dump = (): void => {
            location.href = `${apiPrefix}/store`;
        };

        const restore = async (): Promise<void> => {
            const curFile = file.value?.files?.[0];
            if (file.value === null || curFile === undefined) {
                throw Error("cannot get file");
            }

            if (curFile.type !== "application/gzip" && curFile.type !== "application/x-gzip") {
                ui.notify(i18n.t("messages.database.invalidFile"), "error");
                file.value.value = "";
                return;
            }

            const ok = await ui.modal<boolean>("confirm", {
                button: {
                    label: i18n.t("labels.database.restore"),
                    danger: true,
                },
                message: i18n.t("messages.database.restore"),
            } as ModalConfirmParams);

            if (!ok) {
                return;
            }

            fetch(`${apiPrefix}/store`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/gzip",
                },
                body: curFile,
            })
                .then(onFetch)
                .then(
                    () => {
                        ui.notify(i18n.t("messages.database.restored"), "success");
                    },
                    async response => {
                        const json = await response.json();
                        ui.notify(i18n.t("messages.database.restoreFailed", [json.error]), "error");
                    },
                );
        };

        onMounted(() => {
            ui.title(i18n.t("labels.database._"));
        });

        return {
            apiOptions,
            dump,
            file,
            i18n,
            restore,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../../mixins";

.v-content {
    @include content;

    .v-button {
        margin-top: 1.5rem;
    }

    input[type="file"] {
        display: none;
    }
}
</style>
