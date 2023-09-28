<template>
  <registration-layout>
    <section class="signup-form">
      <basic-questions-form :user-signup-data="userSignupData" @form-submitted="register" />
      <h5 class="validating-signup-field" v-if="user_is_registering">
        Registering...
      </h5>
    </section>

    <signup-verification-message v-if="basicQuestionsSectionIsFinished"
      :user_email="userSignupData.work_email"></signup-verification-message>
  </registration-layout>
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
let basicQuestionsSectionIsFinished = ref<boolean>(false);
let user_is_registering = ref<boolean>(false);
let userSignupData = reactive<UserDataPreVerification>({
  first_name: "oladipea",
  last_name: "bobbingon",
  work_email: "bobby",
  employment_number: 21324324,
  department: "Computer Science",
  phone_number: 3232223323,
});

function register() {
  user_is_registering.value = true
  axios
    .post(
      "https://udm-reimbursement-project.onrender.com/api/verifySignupBasicInformation",
      {
        employmentNumber: userSignupData.employment_number,
        workEmail: userSignupData.work_email,
      }
    )
    .then((res) => {
      window.scrollTo({ top: 0, left: 0, behavior: "smooth" });
    })
    .catch((err) => {

    });
}

function sendConfirmationEmail() {
  axios
    .post(
      `${import.meta.env.VITE_API_URL}/api/send-confirmation-email`, userSignupData,
    )
    .then(() => {
      basicQuestionsSectionIsFinished.value = true;
      // verifyingInformation.value = true;
    })
    .catch((err) => {
      console.log(err)
      alert(err?.response?.data || "An error has occured, please try again");
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
