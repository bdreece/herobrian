<script setup lang="ts">
    import type { Color, Size } from '~/types/primitives';

    defineProps<{
        color?: Color;
        size?: Size;
    }>();
</script>

<template>
    <span
        class="status"
        :class="{
            'status-neutral': color === 'neutral',
            'status-primary': color === 'primary',
            'status-secondary': color === 'secondary',
            'status-accent': color === 'accent',
            'status-info': color === 'info',
            'status-success': color === 'success',
            'status-warning': color === 'warning',
            'status-error': color === 'error',
            'status-xs': size === 'xs',
            'status-sm': size === 'sm',
            'status-md': size === 'md',
            'status-lg': size === 'lg',
            'status-xl': size === 'xl',
        }"
    />
</template>
