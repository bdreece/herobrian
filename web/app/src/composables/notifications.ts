import { useEventBus } from '@vueuse/core';

export interface Notification {
    timestamp: DOMHighResTimeStamp;
    icon: string;
    title: string;
    content: string;
}

export function useNotificationBus() {
    return useEventBus<Notification>('notifications');
}
