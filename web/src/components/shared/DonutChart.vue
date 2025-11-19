<script lang="ts" setup>
import { computed } from 'vue';
import type { RepositoryTotalBandwidths } from '@/api';
import { useI18n } from 'vue-i18n';
import { useTheme } from 'vuetify';

const props = defineProps<{
    totalBandwidths: RepositoryTotalBandwidths;
}>();

const { t } = useI18n();
const theme = useTheme();

const donutOptions = computed(() => {
    return {
        labels: [t('TOTAL_TX'), t('TOTAL_RX')],
        chart: {
            type: 'donut',
            fontFamily: `inherit`,
            foreColor: '#a1aab2',
            toolbar: {
                show: false
            }
        },
        colors: [theme.current.value.colors.primary, theme.current.value.colors.lightprimary, '#F9F9FD'],
        plotOptions: {
            pie: {
                startAngle: 0,
                endAngle: 360,
                donut: {
                    size: '75%',
                    background: 'transparent'
                }
            }
        },
        stroke: {
            show: false
        },

        dataLabels: {
            enabled: false
        },
        legend: {
            show: false
        },
        tooltip: { theme: 'light', fillSeriesColor: false }
    };
});

const chart = computed(() => [+props.totalBandwidths?.tx.toFixed(6) || 0, +props.totalBandwidths?.rx.toFixed(6) || 0]);
</script>

<template>
    <apexchart :options="donutOptions" :series="chart" height="180" type="donut" />
</template>
