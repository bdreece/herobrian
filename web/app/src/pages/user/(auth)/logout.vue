<script setup lang="ts">
    import axios from 'axios';

    definePage({
        name: 'logout',
        meta: {
            allowAnonymous: true,
            layout: 'auth',
            seo: {
                title: 'herobrian - Logout',
            },
        },
    });

    const router = useRouter();
    const token = useToken();

    await axios.post('/api/user/logout');
    token.value = null;
    router.push({ name: 'login' });
</script>
