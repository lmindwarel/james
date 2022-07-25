import { App } from 'vue'
import * as constants from '@/constants'

export default {
    install: (app: App) => {
        app.config.globalProperties.$constants = constants
    },
}

declare module '@vue/runtime-core' {
    interface ComponentCustomProperties {
        $constants: typeof constants;
    }
}