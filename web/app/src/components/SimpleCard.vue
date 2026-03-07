<script setup lang="ts">
    import type { Size } from '~/types/primitives';

    const { actionsTemplate = 'actions', titleTemplate = 'title' } =
        defineProps<{
            horizontal?: boolean;
            size?: Size;
            title?: string;
            titleTemplate?: string;
            actionsTemplate?: string;
        }>();
</script>

<template>
    <div
        class="card shadow-sm"
        :class="{
            'card-xs': size === 'xs',
            'card-sm': size === 'sm',
            'card-md': size === 'md',
            'card-lg': size === 'lg',
            'card-xl': size === 'xl',
            'card-side': horizontal,
        }"
    >
        <div class="card-body">
            <h2 class="card-title">
                <slot :name="titleTemplate">
                    {{ title }}
                </slot>
            </h2>

            <slot />

            <div
                v-if="$slots[actionsTemplate]"
                class="justify-end card-actions"
            >
                <slot :name="actionsTemplate" />
            </div>
        </div>
    </div>
</template>
