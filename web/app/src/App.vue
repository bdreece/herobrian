<template>
    <div>
        <Suspense>
            <RouterView v-slot="{ Component, route }">
                <component
                    :is="Component"
                    :key="route"
                />
            </RouterView>
            <template #fallback>
                <LoadingSpinner />
            </template>
        </Suspense>
    </div>
</template>
