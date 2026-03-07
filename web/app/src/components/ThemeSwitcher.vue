<script setup lang="ts">
    const { toggle } = useTheme();
</script>

<template>
    <label class="swap swap-rotate">
        <input
            type="checkbox"
            @change="toggle()"
        />

        <i class="swap-on iconify solar--sun-line-duotone" />
        <i class="swap-off iconify solar--moon-line-duotone" />
    </label>
</template>
