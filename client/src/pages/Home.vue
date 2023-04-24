<template>
    <div>
        {{ $t('pages.home.name') }}
    </div>
</template>

<script>
    import { defineComponent } from 'vue';
    import { useAuth } from '../stores/auth.js'

    export default defineComponent({
        name: 'Home',
        setup() {
            useAuth()
        }
    })
</script>