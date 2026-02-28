<script setup lang="ts">
    definePage({
        meta: {
            allowAnonymous: true,
            layout: 'auth',
        },
    });

    const formId = useId();

    function onSubmit(e: SubmitEvent) {
        e.preventDefault();
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
