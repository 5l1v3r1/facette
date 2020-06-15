<template>
    <v-modal name="chart-marker">
        <v-form slot-scope="modal">
            <v-label>{{ $t("labels.value") }}</v-label>
            <v-input required v-autofocus v-model.number="marker.value"></v-input>

            <v-label>{{ $tc("labels.labels", 1) }}</v-label>
            <v-input :placeholder="String(marker.value)" v-model="marker.label"></v-input>

            <v-label>{{ $t("labels.color") }}</v-label>
            <v-color class="half" v-model="marker.color"></v-color>

            <v-label>{{ $tc("labels.charts.axes._", 1) }}</v-label>
            <v-select
                class="half"
                :options="axes"
                :placeholder="$t('labels.charts.axes.select')"
                v-model="marker.axis"
            ></v-select>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.close") }}</v-button>
                <v-button primary @click="modal.close(marker)">{{ $t("labels.ok") }} </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import cloneDeep from "lodash/merge";
import merge from "lodash/merge";
import {Component, Vue} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

import {resolveOption} from "@/src/helpers/select";

export interface ModalChartMarkerParams {
    marker: Marker;
}

const axes: Array<SelectOption> = [
    {label: "labels.charts.axes.left", value: "left"},
    {label: "labels.charts.axes.right", value: "right"},
];

const defaultMarker: Marker = {
    value: 0,
    label: "",
    color: "",
    axis: "left",
};

@Component
export default class ModalChartMarkerComponent extends Vue {
    public marker: Marker = cloneDeep(defaultMarker);

    public get axes(): Array<SelectOption> {
        return axes.map(option => resolveOption(this, option));
    }
    public onModalShow(params: ModalChartMarkerParams): void {
        this.marker = merge({}, defaultMarker, params.marker);
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    .v-form > .v-label:first-child {
        margin-top: 0;
    }

    .v-color,
    .v-select {
        min-width: auto;
        width: 100%;
    }
}
</style>
