<script setup lang="ts">
    definePage({
        name: 'login',
        meta: {
            allowAnonymous: true,
            layout: 'auth',
            seo: {
                title: 'herobrian \u2014 Login',
            },
        },
    });

    const formId = useId();
    const router = useRouter();
    const token = useToken();
    const { mutate, error } = useLogin({
        onSuccess({ accessToken }) {
            token.value = accessToken;
            router.push({ name: 'home' });
        },
    });

    function onSubmit(e: SubmitEvent) {
        e.preventDefault();
        const form = e.target as HTMLFormElement;
        const data = new FormData(form);
        /* @ts-expect-error 2345 */
        const params = new URLSearchParams(data);

        mutate(params);
    }
</script>

<template>
    <div class="container mx-auto">
        <div class="card w-96 bg-base-200 shadow-sm mx-auto my-32">
            <div class="card-body">
                <h2 class="card-title mb-3">Login</h2>
                <form
                    :id="formId"
                    @submit="onSubmit"
                >
                    <FormControl
                        label="Username"
                        type="text"
                        name="displayName"
                        placeholder="awesomeshooter12"
                        autocomplete="username"
                        maxlength="63"
                        required
                    />

                    <FormControl
                        label="Password"
                        type="password"
                        name="password"
                        placeholder="********"
                        autocomplete="current-password"
                        required
                    />
                </form>

                <div v-if="error">
                    {{ error.message }}
                </div>

                <div class="card-actions justify-between">
                    <label class="label">
                        <input
                            class="checkbox"
                            type="checkbox"
                            name="rememberMe"
                        />
                        Remember Me?
                    </label>

                    <button
                        class="btn btn-primary"
                        type="submit"
                        :form="formId"
                    >
                        Submit
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>
