import type { UseSeoMetaInput } from '@unhead/vue';
import 'vue-router';

declare module 'vue-router' {
    interface RouteMeta {
        breadcrumb?: string;
        seo?: UseSeoMetaInput;
        layout?: string;
        transition?: string;
        allowAnonymous?: boolean;
        authorize?: string[];
    }
}
