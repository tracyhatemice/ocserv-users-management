<script lang="ts" setup>
import {
    HomeApi,
    type HomeCurrentStats,
    type HomeGeneralInfo,
    type HomeGetHomeUser,
    type ModelsDailyTraffic,
    type ModelsIPBanPoints,
    type RepositoryTopBandwidthUsers,
    type RepositoryTotalBandwidths
} from '@/api';
import { onMounted, ref } from 'vue';
import { getAuthorization } from '@/utils/request';

import ServerGeneralInfoOverview from '@/components/dashboard/ServerGeneralInfoOverview.vue';
import ServerCurrentStatsOverview from '@/components/dashboard/ServerCurrentStatsOverview.vue';
import OnlineUsersOverview from '@/components/dashboard/OnlineUsersOverview.vue';
import IpBansPointOverview from '@/components/dashboard/IpBansPointOverview.vue';
import TopBandwidthUsers from '@/components/dashboard/TopBandwidthUsers.vue';
import UsersOverview from '@/components/dashboard/UsersOverview.vue';
import RxTxDonutOverview from '@/components/dashboard/RxTxDonutOverview.vue';
import RxTxChartOverview from '@/components/dashboard/RxTxChartOverview.vue';
import UiParentCard from '@/components/shared/UiParentCard.vue';
import { dummyBanIPs, dummyOnlineUsers, dummyTrafficData } from '@/data/dummy';

const trafficData = ref<ModelsDailyTraffic[]>([]);
const users = ref<HomeGetHomeUser>({});
const ipBanPoints = ref<ModelsIPBanPoints[]>([]);
const topUsers = ref<RepositoryTopBandwidthUsers>({});
const currentStats = ref<HomeCurrentStats>({});
const generalInfo = ref<HomeGeneralInfo>({});
const totalBandwidths = ref<RepositoryTotalBandwidths>({ rx: 0, tx: 0 });

onMounted(() => {
    const api = new HomeApi();
    api.homeGet(getAuthorization()).then((res) => {
        generalInfo.value = res.data?.server_status.general_info || {};
        currentStats.value = res.data?.server_status.current_stats || {};
        users.value = res.data?.users || {};
        trafficData.value = res.data?.statistics || [];
        ipBanPoints.value = res.data?.ip_bans || [];
        topUsers.value = res.data?.top_bandwidth_user || {};
        totalBandwidths.value = res.data?.total_bandwidth || { rx: 0, tx: 0 };

        // // TODO: Remove it
        // // dummy data
        // import { dummyBanIPs, dummyOnlineUsers, dummyTrafficData } from '@/utils/dummy';
        // ipBanPoints.value = dummyBanIPs;
        // users.value = {
        //     total: 4,
        //     online_users_session: dummyOnlineUsers
        // };
        // trafficData.value = dummyTrafficData
        // totalBandwidths.value = {rx: 10.444, tx: 33.44}
    });
});
</script>

<template>
    <v-row>
        <v-col cols="12">
            <UiParentCard>
                <v-row>
                    <!-- Rx Tx overview -->
                    <v-col cols="12" lg="8">
                        <RxTxChartOverview :data="trafficData" />
                    </v-col>

                    <!-- User Overview / Rx Tx overview -->
                    <v-col cols="12" lg="4">
                        <div class="mb-6">
                            <UsersOverview :users="users" />
                        </div>
                        <div>
                            <RxTxDonutOverview :totalBandwidths="totalBandwidths" />
                        </div>
                    </v-col>

                    <!-- Server General Info Overview -->
                    <v-col cols="12" lg="6" sm="6">
                        <ServerGeneralInfoOverview :generalInfo="generalInfo" />
                    </v-col>

                    <!-- Server Stats Overview -->
                    <v-col cols="12" lg="6" sm="6">
                        <ServerCurrentStatsOverview :currentStats="currentStats" />
                    </v-col>

                    <!-- Online Users OverView -->
                    <v-col cols="12" lg="8">
                        <OnlineUsersOverview :sessions="users?.online_users_session || []" />
                    </v-col>

                    <!-- IP Ban Points OverView -->
                    <v-col cols="12" lg="4">
                        <IpBansPointOverview :ipBanPoints="ipBanPoints" />
                    </v-col>

                    <!-- Top Bandwidth Users OverView -->
                    <v-col cols="12">
                        <TopBandwidthUsers :topUsers="topUsers" />
                    </v-col>
                </v-row>
            </UiParentCard>
        </v-col>
    </v-row>
</template>
