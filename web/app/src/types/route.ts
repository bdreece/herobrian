import 'vue-router';

declare module 'vue-router' {
    interface RouteMeta {
        layout?: string;
        transition?: string;
        allowAnonymous?: boolean;
    }
}
