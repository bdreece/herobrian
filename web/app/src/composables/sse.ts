import type { EventHookOn } from '@vueuse/core';

export interface UseSseReturn<T> {
    eventSource: Ref<EventSource | undefined>;
    onOpen: EventHookOn<[]>;
    onMessage: EventHookOn<[MessageEvent<T>]>;
    onError: EventHookOn<[]>;
}

export function useSse<T>(stream: string): UseSseReturn<T> {
    const params = new URLSearchParams({ stream });
    const url = new URL(`/event?${params}`, import.meta.env.BASE_URL);
    const eventSource = ref<EventSource | undefined>(undefined);

    const controller = new AbortController();
    const open = createEventHook<[]>();
    const message = createEventHook<[MessageEvent<T>]>();
    const error = createEventHook<[]>();

    onMounted(() => {
        eventSource.value = new EventSource(url);

        eventSource.value.addEventListener('open', () => open.trigger(), {
            signal: controller.signal,
        });

        eventSource.value.addEventListener('message', e => message.trigger(e), {
            signal: controller.signal,
        });

        eventSource.value.addEventListener('error', () => error.trigger(), {
            signal: controller.signal,
        });
    });

    onUnmounted(() => {
        eventSource.value?.close();
        controller.abort();
    });

    return {
        eventSource,
        onOpen: open.on,
        onMessage: message.on,
        onError: error.on,
    };
}
