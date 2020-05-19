<template>
    <v-modal name="dashboard-item" ref="modal">
        <v-form slot-scope="modal" v-if="step === 1">
            <v-radio-group name="type" v-autofocus v-model="item.type">
                <v-radio :value="type" :key="type" v-for="(item, type) in types">
                    <v-icon :icon="item.icon"></v-icon>
                    {{ $t(`labels.dashboards.types.${type}`) }}
                </v-radio>
            </v-radio-group>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.cancel") }}</v-button>
                <v-button primary @click="next">{{ $t("labels.continue") }}</v-button>
            </template>
        </v-form>

        <v-form slot-scope="modal" v-else-if="step === 2">
            <v-message type="error" v-if="formFailed">
                {{ $t("messages.dashboards.types.loadFailed", [item.type]) }}
            </v-message>

            <component :is="form" :options="item.options"></component>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.cancel") }}</v-button>
                <v-button primary :disabled="formFailed" @click="modal.close(item)">
                    {{ $t("labels.items.add") }}
                </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import merge from "lodash/merge";
import {VueConstructor} from "vue";
import {Component, Vue, Watch} from "vue-property-decorator";

export interface ModalDashboardItemParams {
    step: number;
    item: DashboardItem;
}

type ItemTypes = Record<string, {icon: string; options: Record<string, unknown>}>;

const types: ItemTypes = {
    chart: {
        icon: "chart-area",
        options: {
            id: "",
        },
    },
};

const defaultItem: DashboardItem = {
    type: "chart",
    layout: {x: 0, y: 0, w: 1, h: 1},
    options: {},
};

@Component
export default class ModalDashboardItemComponent extends Vue {
    public form: VueConstructor | null = null;

    public formFailed = false;

    public item: DashboardItem = cloneDeep(defaultItem);

    public step = 1;

    public types: ItemTypes = types;

    public next(): void {
        this.step++;
    }

    public onModalShow(params: ModalDashboardItemParams): void {
        if (params.step) {
            this.step = params.step;
        }

        this.item = merge({}, defaultItem, params.item);
    }

    @Watch("item.type", {immediate: true})
    public async onType(to: DashboardItemType): Promise<void> {
        merge(this.item.options, types[to].options);

        try {
            Object.assign(this, {
                form: (await import(`../dashboard/${to}-item.vue`)).default,
                formFailed: false,
            });
        } catch {
            this.formFailed = true;
        }
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    .v-radio .v-icon {
        font-size: 1.2rem;
        margin: 0 0.5rem;
        width: 1.25rem;
    }
}
</style>
