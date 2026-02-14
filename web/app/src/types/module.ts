import type { VueHeadClient } from '@unhead/vue';
import type { Promisable } from 'type-fest';
import type { App } from 'vue';
import type { Router, RouteRecordRaw } from 'vue-router';

declare global {
    namespace Herobrian {
        interface Module {
            install?(ctx: Module.Context): Promisable<void>;
        }
    }

    namespace Herobrian.Module {
        interface Context {
            app: App<Element>;
            router: Router;
            routes: readonly RouteRecordRaw[];
            head: VueHeadClient;
        }
    }
}

export {};
