import { defineStore } from 'pinia';
import { OCCTLApi, SystemApi } from '@/api';
import type { ConfigState, ServerState } from '@/types/storeTypes/StoreConfigType';

export const useServerStore = defineStore('server', {
    state: (): ServerState => ({
        Version: '',
        OcctlVersion: ''
    }),
    actions: {
        async getServerInfo() {
            const api = new OCCTLApi();
            await api
                .occtlServerInfoGet()
                .then((res) => {
                    if (res.data) {
                        this.Version = res.data.version || '';
                        this.OcctlVersion = (res.data.occtl_version || '').replace(/\n/g, '<br />');
                    }
                })
                .catch(() => {});
        }
    },
    getters: {
        versionInfo: (state) => state.Version,
        occtlVersionInfo: (state) => state.OcctlVersion
    }
});

export const useConfigStore = defineStore('config', {
    state: (): ConfigState => ({
        setup: false,
        googleCaptchaSiteKey: ''
    }),

    actions: {
        async getConfig() {
            const api = new SystemApi();
            await api.systemInitGet().then((res) => {
                if (res.data) {
                    this.googleCaptchaSiteKey = res.data.google_captcha_site_key || '';
                    this.setup = true;
                }
            });
            return this.setup;
        },
        setConfig(googleCaptchaSiteKey: string | undefined) {
            if (googleCaptchaSiteKey) {
                this.googleCaptchaSiteKey = googleCaptchaSiteKey;
            }
            this.setup = true;
        }
    },
    getters: {
        config(state): ConfigState {
            return {
                setup: state.setup,
                googleCaptchaSiteKey: state.googleCaptchaSiteKey
            };
        }
    }
});
