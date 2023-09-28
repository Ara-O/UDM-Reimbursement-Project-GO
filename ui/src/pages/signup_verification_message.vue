<template>
    <article>
        <h3 class="thanks-for-signing-up-msg">Thank you for signing up!</h3>
        <h4 class="email-confirmation-text">
            A confirmation email has being sent to
            {{ user_email.trim() }}@udmercy.edu. You should
            receive an email in 3-5 minutes. Please check your spam/junk folder
            if it does not arrive. You can close out of this page.
        </h4>
    </article>
</template>

<script lang="ts" setup>
defineProps<{
    user_email: string
}>()
</script>