<template>
    <div class="flex-1 flex flex-col">
        <DefaultHeader :items="[]" />
        <main class="flex-1 grid place-items-center">
            <RouterView />
        </main>
    </div>
</template>
