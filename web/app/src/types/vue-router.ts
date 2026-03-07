import type { UseSeoMetaInput } from '@unhead/vue';
import type { MaybeRouteGetter } from '~/composables/route';
import 'vue-router';

declare module 'vue-router' {
    interface RouteMeta {
        breadcrumb?: MaybeRouteGetter<string>;
        seo?: UseSeoMetaInput;
        layout?: string;
        transition?: string;
        allowAnonymous?: boolean;
        authorize?: string[];
    }
}
