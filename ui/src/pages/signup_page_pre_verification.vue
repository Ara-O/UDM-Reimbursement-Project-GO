<template>
  <registration-layout v-if="!basicQuestionsFormIsComplete">
    <section class="signup-form">
      <basic-questions-form :user-signup-data="userSignupData" @form-submitted="register" />
      <h5 class="validating-signup-field" v-if="registrationInProgress">
        {{ progressMessage }}...
      </h5>
    </section>
  </registration-layout>
  <signup-verification-message v-if="basicQuestionsFormIsComplete"
    :user_email="userSignupData.work_email"></signup-verification-message>
</template>

<script lang="ts" setup>
import SignupVerificationMessage from "./signup_verification_message.vue";
import RegistrationLayout from "../layouts/registration_layout.vue"
import BasicQuestionsForm from "../components/signup-flow/BasicQuestionsForm.vue";
import axios from "axios";
import { onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { UserDataPreVerification } from "../types/types";

const router = useRouter();
let basicQuestionsFormIsComplete = ref<boolean>(false);
let registrationInProgress = ref<boolean>(false);
let progressMessage = ref<string>("");
let userSignupData = reactive<UserDataPreVerification>({
  first_name: "oladipea",
  last_name: "bobbingon",
  work_email: "oladipea",
  employment_number: 21324324,
  department: "Computer Science",
  phone_number: 3232223323,
});

function updateProgressMessage(message: string) {
  progressMessage.value = message
  registrationInProgress.value = true
}

async function register() {
  updateProgressMessage("Verifying information...")
  try {
    await axios
      .post(
        `${import.meta.env.VITE_API_URL}/api/verify-user-information`,
        {
          employment_number: userSignupData.employment_number,
          work_email: userSignupData.work_email,
        }
      )

    sendConfirmationEmail()
  } catch (err: any) {
    console.log(err)
    if (err?.response?.status === 409) {
      updateProgressMessage("There is already a similar record in our system. Please check your input and try again.")
    } else {
      updateProgressMessage(err?.response?.data || "There was an error registering, please try again")
    }
  }
}

function sendConfirmationEmail() {
  updateProgressMessage("Registering user...")
  axios
    .post(
      `${import.meta.env.VITE_API_URL}/api/send-confirmation-email`, userSignupData,
    )
    .then(() => {
      basicQuestionsFormIsComplete.value = true;
    })
    .catch((err) => {
      if (err.response.status === 400) {
        updateProgressMessage("Invalid user input");
      } else {
        updateProgressMessage("An error has occured, please try again later");
      }
    });
}

onMounted(() => {
  if (localStorage.getItem("token")?.length ?? 0 > 0) {
    router.push("/dashboard");
  }
});
</script>

<style scoped>
@import url("../assets/styles/signup-page.css");
</style>
