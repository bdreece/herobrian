<script setup lang="ts"></script>

<template>
    <RouterView v-slot="{ Component, route }">
        <Transition :name="route.meta.transition!">
            <component :is="Component" />
        </Transition>
    </RouterView>
</template>
