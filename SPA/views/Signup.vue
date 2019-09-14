<template>
  <form v-on:submit.prevent="onSubmit">
    <!-- NAME -->
    <div :class="{ 'is-invalid': submitted && $v.name.$error }">
      <label>Name</label>
      <input v-model.trim="$v.name.$model" />
    </div>
    <div v-if="submitted && !$v.name.required">
      Name is required.
    </div>
    <div class="error" v-if="submitted && !$v.name.nameChars">
      Name must be between 1 and 15 characters and contain only standard letters, numbers, or underscore ("_").
    </div>

    <!-- EMAIL -->
    <div :class="{ 'is-invalid': $v.email.$error }">
      <label >Email</label>
      <input v-model.trim="$v.email.$model" />
    </div>
    <div v-if="submitted && !$v.email.required">
      Email is required.
    </div>
    <div v-if="submitted && !$v.email.email">
      Please enter a valid email.
    </div>

    <!-- PASSWORD -->
    <div :class="{ 'is-invalid': $v.password.$error }">
      <label>Password</label>
      <input type="password" v-model.trim="$v.password.$model"/>
    </div>
    <div class="error" v-if="submitted && !$v.password.required">
      Password is required.
    </div>
    <div class="error" v-if="submitted && !$v.password.passwordChars">
      Password must be at least 8 characters and contain at least one lower case letter, one upper case letter, one number, and one special character.
    </div>

    <!-- CONFIRM PASSWORD -->
    <div :class="{ 'is-invalid': $v.confirmPassword.$error }">
      <label class="form__label">Confirm Password</label>
      <input type="password" v-model.trim="$v.confirmPassword.$model"/>
    </div>
    <div class="error" v-if="submitted && !$v.confirmPassword.required">
      Please confirm password.
    </div>
    <div class="error" v-if="submitted && !$v.confirmPassword.sameAsPassword">
      Passwords must match.
    </div>

    <div class="error" v-if="apiError">
      {{ apiError }}
    </div>
    <button class="button" type="submit" :disabled="pending">Sign Up</button>
  </form>
</template>

<script>
//import additional validators
import { required, email, helpers, sameAs } from "vuelidate/lib/validators";

//regular expression for name requirements
const nameChars = helpers.regex("nameChars", /^[a-zA-Z0-9_]{1,15}$/);

//regular expression for password requirements
const passwordChars = helpers.regex("passwordChars", /(?=^.{8,}$)((?=.*\d)(?=.*\W+))(?![.\n])(?=.*[A-Z])(?=.*[a-z]).*$/);

export default {
  name: "signup",
  data: function() {
    return {
      name: "",
      email: "",
      password: "",
      confirmPassword: "",
      submitted: false,
      pending: false,
      apiError: ""
    };
  },
  validations: {
    //set form requirements
    name: { required, nameChars },
    email: { required, email },
    password: { required, passwordChars },
    confirmPassword: { required, sameAsPassword: sameAs("password") }
  },
  methods: {
    onSubmit: function() {
      this.submitted = true;

      this.$v.$touch();
      if (this.$v.$invalid) {
        return;
      }

      this.pending = true;

      this.$axios
        .post("/api/signup", {
          name: this.name,
          email: this.email,
          password: this.confirmPassword
        })
        .then(response => {
          if (response.status == 200) {
            //on success log in the user and redirect to home
            let token = this.$cookie.get("token-hp");
            let id = JSON.parse(window.atob(token.split(".")[1])).id; //split token, get payload (position [1]), b64 decode, parse to JSON, grab id
            this.$store
              .dispatch("updateUser", {
                userID: id
              })
              .then(() => {
                this.$router.push({ name: "home" });
              });
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error signing up";
          } else {
            this.apiError = "error signing up";
          }
        })
        .finally(() => {
          //on success or error reset the password variables and submitted and pending statuses
          this.password = "";
          this.confirmPassword = "";
          this.submitted = false;
          this.pending = false;
        });
    }
  }
};
</script>
