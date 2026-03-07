<script setup lang="ts">
    import type { Size } from '~/types/primitives';

    const { size = 'md' } = defineProps<{
        size?: Size;
    }>();
</script>

<template>
    <span
        class="loading loading-spinner"
        :class="{
            'loading-xs': size === 'xs',
            'loading-sm': size === 'sm',
            'loading-md': size === 'md',
            'loading-lg': size === 'lg',
            'loading-xl': size === 'xl',
        }"
    />
</template>
