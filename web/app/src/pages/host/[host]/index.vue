<script setup lang="ts">
    definePage({
        name: 'host',
        meta: {},
    });

    const hostname = useRouteParams<string>('host');
    const host = useHost(hostname.value);

    const gigabyte = new Intl.NumberFormat(navigator.language, {
        style: 'unit',
        unit: 'gigabyte',
        unitDisplay: 'short',
    });
</script>

<template>
    <div class="stats bg-base-300 shadow">
        <div class="stat">
            <div class="stat-title">Type</div>
            <div class="stat-value">{{ host?.image.type }}</div>
        </div>
        <div class="stat">
            <div class="stat-title">CPU</div>
            <div class="stat-value">
                {{
                    (host?.processor.cores || 0)
                    * (host?.processor.threads || 0)
                }}
                Cores
            </div>
        </div>
        <div class="stat">
            <div class="stat-title">RAM</div>
            <div class="stat-value">
                {{ gigabyte.format((host?.image.memory || 0) / 1024) }}
            </div>
        </div>
        <div class="stat">
            <div class="stat-title">Architecture</div>
            <div class="stat-value">{{ host?.arch }}</div>
        </div>
        <div class="stat">
            <div class="stat-title">Platform</div>
            <div class="stat-value">{{ host?.platform }}</div>
        </div>
        <div class="stat">
            <div class="stat-title">Network Performance</div>
            <div class="stat-value">{{ host?.image.network }}</div>
        </div>
    </div>
</template>
