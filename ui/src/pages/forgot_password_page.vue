<template>
  <section class="login-page">
    <!-- Img is a way you can embed images to your site,
      the ../ before the folder means that we are traversing file systems -->
    <div class="udmercy-logo-wrapper">
      <img src="../assets/detroit-mercy-logo.png" alt="Detroit mercy logo" class="udmercy-logo" />
    </div>
    <br />
    <h3 class="login-title">Detroit Mercy Reimbursement System</h3>
    <br />
    <Form @submit="resetPassword" class="login-form">
      <div class="login-field" style="width: 250px">
        <label for="new-password" style="width: 100%; text-align: center; margin-bottom: 15px;">New Password: </label>
        <span style="width: 100%">
          <Field type="password" class="login-password-input" name="new-password" v-model="password" id="new-password"
            :rules="isNotEmpty" style="width: 100%; margin-left: 0px; text-align: center;" />
          <ErrorMessage as="h3" name="new-password" class="error-field"
            style="text-align: center; width: 100%; max-width: 100%; margin-top: 10px;" />
        </span>
      </div>
      <div class="login-field" style="margin-bottom: 20px; width: 250px">
        <label for="confirm-password" style=" text-align: center; margin-bottom: 15px;width: 100%">Confirm New
          Password:</label>
        <span style="width: 100%">
          <Field v-model="confirmPassword" type="password" class="login-password-input" name="confirm-password"
            style="width: 100%; margin-left: 0px; text-align: center" :rules="isNotEmpty" id="confirm-password " />
          <ErrorMessage as="h3" name="confirm-password" class="error-field"
            style="text-align: center; width: 100%;max-width: 100%;margin-top: 10px" />
        </span>
      </div>
      <span style="display: flex; align-items: center; margin-top: -20px; gap: 10px">
        <router-link to="/signup" style="font-size: 14px">Create an Account</router-link>
        <h3 style="font-weight: 300">|</h3>
        <router-link to="/" style="font-size: 14px">Log In</router-link>
      </span>
      <button class="login-button" style="width: auto">Change Password</button>
    </Form>
  </section>
</template>

<script lang="ts" setup>
import axios from "axios";
import { Form, Field, ErrorMessage } from "vee-validate";
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { isNotEmpty } from "../utils/validators";

let password = ref<string>("bob");
let confirmPassword = ref<string>("bob");
const route = useRoute();
const router = useRouter();
let userToken = ref<string>("");

function resetPassword() {
  if (password.value !== confirmPassword.value) {
    alert("Passwords do not match");
    return
  }

  axios
    .post(
      `${import.meta.env.VITE_API_URL}/api/reset-password`,
      {
        token: userToken.value,
        new_password: password.value,
      }
    )
    .then((res) => {
      console.log(res);
      alert("Password successfully changed");
    })
    .catch((err) => {
      console.log(err)
      if (err?.response?.status === 403) {
        alert("Invalid token, please start the process again")
        router.push("/");
      }
    })
    .finally(() => {
      router.push("/");
    });

}

onMounted(() => {
  if (route.params.userToken === null) {
    router.push("/");
    return
  }

  axios.post(`${import.meta.env.VITE_API_URL}/api/verify-forgot-password-token`, {
    user_token: route.params.userToken
  }).then(() => {
    userToken.value = route.params.userToken as string;
  }).catch((err) => {
    if (err.response.status === 403) {
      alert("Invalid token, please try again")
    } else {
      alert("There was an error, please try again")
    }
    router.push("/")
  })
});
</script>

<style scoped>
@import url("../assets/styles/login-page.css");

.error-field {
  font-weight: 400;
  margin-top: 4px;
}
</style>
