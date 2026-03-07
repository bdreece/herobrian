<script setup lang="ts">
    const { claims } = useAuth();
</script>

<template>
    <div
        v-if="claims"
        class="dropdown dropdown-end"
    >
        <div
            tabindex="0"
            role="button"
            class="btn btn-ghost btn-circle avatar"
        >
            <div class="w-10 rounded-full">
                <img
                    alt="Avatar"
                    :src="
                        claims.picture_url
                        ?? `https://mc-heads.net/avatar/${claims.preferred_username}`
                    "
                />
            </div>
        </div>

        <ul
            tabindex="-1"
            class="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow"
        >
            <li>
                <RouterLink to="/user/profile">Profile</RouterLink>
            </li>
            <li>
                <RouterLink to="/user/logout">Logout</RouterLink>
            </li>
        </ul>
    </div>
</template>
