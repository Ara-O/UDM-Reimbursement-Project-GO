<template>
    <section class="signup-page">
        <section class="left-section">
            <div class="udmercy-logo-wrapper">
                <img src="../assets/detroit-mercy-logo.png" alt="Detroit mercy logo" class="udmercy-logo" />
            </div>
        </section>
        <section class="right-section">
            <div class="udmercy-logo-wrapper-mobile">
                <img src="../assets/detroit-mercy-logo.png" alt="Detroit mercy logo" class="udmercy-logo-mobile" />
            </div>
            <h3 class="signup-title">Detroit Mercy Reimbursement System</h3>
            <!--  -->
            <slot></slot>
            <!--  -->
            <div class="account-feedback-message">
                <router-link to="/" class="already-have-account">Already have an
                    Account</router-link>
                <h5 class="required-field-note">
                    Note: All required fields (*) must be filled
                </h5>
            </div>
        </section>
    </section>
</template>

<style scoped>
@import url("/styles/registration/registration_layout.css")
</style>