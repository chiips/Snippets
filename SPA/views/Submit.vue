<template>
  <div>
    <!-- on successful submission, display the submission and a button to submit another -->
    <div v-if="success">
      <p>post submitted!</p>
      <p>{{ title }}</p>
      <p>{{ body }}</p>
      <button @click="another">Submit another post</button>
    </div>
    <!-- otherwise display the submission form -->
    <div v-else>
      <form v-on:submit.prevent="onSubmit">
        <!-- TITLE -->
        <div>
          <label>Title</label>
          <input v-model.trim="$v.title.$model" />
        </div>
        <div v-if="submitted && !$v.title.required">
          Title is required.
        </div>
        <div v-if="submitted && !$v.title.maxLength">
          Title must be less than {{ $v.title.$params.maxLength.max }} characters.
        </div>

        <!-- BODY -->
        <div>
          <label>Body</label>
          <textarea-autosize ref="body" v-model.trim="$v.body.$model"></textarea-autosize>
        </div>
        <div v-if="submitted && !$v.body.required">
          Body is required.
        </div>
        <div v-if="submitted && !$v.body.maxLength">
          Body must be less than {{ $v.title.$params.maxLength.max }} characters.
        </div>

        <div v-if="apiError">
          {{ apiError }}
        </div>
        <button type="submit" :disabled="pending">Submit</button>
      </form>
    </div>
  </div>
</template>

<script>
import { required, maxLength } from "vuelidate/lib/validators";

export default {
  name: "submit",
  data: function() {
    return {
      title: "",
      body: "",
      submitted: false,
      pending: false,
      apiError: "",
      success: false
    };
  },
  validations: {
    //title and body requirements for the form
    title: { required, maxLength: maxLength(50) },
    body: { required, maxLength: maxLength(5000) }
  },
  methods: {
    onSubmit: function() {
      this.submitted = true;

      //return if the form is invalid
      this.$v.$touch();
      if (this.$v.$invalid) {
        return;
      }

      this.pending = true;

      this.title = this.title;
      this.body = this.body;

      this.$axios
        .post("/api/post", {
          title: this.title,
          body: this.body
        })
        .then(response => {
          if (response.status == 200) {
            this.success = true;
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error submitting";
          } else {
            this.apiError = "error submitting";
          }
        })
        .finally(() => {
          this.submitted = false;
          this.pending = false;
        });
    },
    another: function() {
      //if the user clicks the button for another submission, reset the page
      this.success = false;
      this.title = "";
      this.body = "";
    }
  }
};
</script>
