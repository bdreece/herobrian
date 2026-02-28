<script setup lang="ts">
    defineOptions({
        inheritAttrs: false,
    });

    defineProps<{
        label: string;
    }>();
</script>

<template>
    <label class="floating-label mb-3">
        <span>{{ label }}</span>
        <input
            class="input input-md w-full"
            v-bind="$attrs"
        />
    </label>
</template>
