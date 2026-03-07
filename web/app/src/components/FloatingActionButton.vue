<script setup lang="ts">
    export interface Item {
        text?: string;
        tooltip?: string;
        icon?: string;
        onClick?: (e: MouseEvent) => void;
    }

    const { items = [] } = defineProps<{
        items?: Item[];
    }>();
</script>

<template>
    <div class="fab">
        <div
            tabindex="0"
            role="button"
            class="btn btn-lg btn-circle btn-primary"
        >
            <slot />
        </div>

        <div
            v-for="item of items"
            :key="item.text"
            class="tooltip tooltip-left"
            :data-tip="item.tooltip"
        >
            <button
                class="btn btn-lg btn-circle"
                @click="item.onClick ?? (() => {})"
            >
                <i
                    v-if="item.icon"
                    class="iconify"
                    :class="item.icon"
                />
                <span
                    v-else
                    v-text="item.text"
                />
            </button>
        </div>
    </div>
</template>
