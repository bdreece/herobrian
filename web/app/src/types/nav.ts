import type { RouteLocationRaw } from 'vue-router';

export interface NavItem {
    id: string;
    text: string;
    icon?: string;
    to: RouteLocationRaw;
}
