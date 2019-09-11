<template>
  <form v-on:submit.prevent>
    <!-- TITLE -->
    <div :class="{ 'is-invalid': submitted && $v.title.$error }">
      <input v-model.trim="$v.title.$model" :disabled="!editing"/>
    </div>
    <div class="error" v-if="submitted && !$v.title.required">
      Title is required.
    </div>
    <div class="error" v-if="submitted && !$v.title.maxLength">
      Title must be less than {{ $v.title.$params.maxLength.max }} characters.
    </div>

    <!-- BODY -->
    <div :class="{ 'is-invalid': submitted && $v.body.$error }">

      <!-- display posts in a textarea that is activated when editing -->
      <textarea-autosize ref="body" v-model.trim="$v.body.$model" :disabled="!editing"
      ></textarea-autosize>

    </div>
    <div v-if="submitted && !$v.body.required">
      Body is required.
    </div>
    <div v-if="submitted && !$v.body.maxLength">
      Body must be less than {{ $v.title.$params.maxLength.max }} characters.
    </div>

    <p>{{ created }}</p>
    <p v-if="beenUpdated">{{ updated }}</p>

    <div v-if="apiError">
      {{ apiError }}
    </div>

    <!-- if editing then display the option to save those edits or delete the post -->
    <div v-if="editing">
      <button v-on:click="updatePost()" :disabled="pending">Save</button>
      <button v-on:click="deletePost()" :disabled="pending">Delete</button>
    </div>
    <!-- otherwise display the option to initiate editing mode -->
    <div v-else>
      <button v-on:click.prevent="editing = true">Edit</button>
    </div>

    <!--MODAL-->
    <ConfirmComp
      v-show="isModalVisible"
      @confirm="confirmDeletion"
      @close="closeModal"
      v-bind:modalprop="postprop.title"
    />
  </form>
</template>

<script>
import { required, maxLength } from "vuelidate/lib/validators";
import ConfirmComp from "@/components/ConfirmComp.vue";

export default {
  name: "PostsProfileComp",
  props: {
    postprop: Object,
    pending: Boolean
  },
  components: {
    ConfirmComp
  },
  data: function() {
    return {
      id: this.postprop.id,
      title: this.postprop.title,
      body: this.postprop.body,
      authorID: this.postprop && this.postprop.author ? this.postprop.author.id : "",
      created: new Date(this.postprop.created).toDateString(),
      updated: new Date(this.postprop.updated).toDateString(),
      submitted: false,
      editing: false,
      isModalVisible: false,
      resolved: null,
      apiError: ""
    };
  },
  computed: {
    beenUpdated() {
      if (this.updated === this.created) {
        return false;
      }
      return true;
    }
  },
  validations: {
    title: { required, maxLength: maxLength(50) },
    body: { required, maxLength: maxLength(5000) }
  },
  methods: {
    updatePost: function() {
      this.submitted = true;

      this.$v.$touch();
      if (this.$v.$invalid) {
        return;
      }

      let updatedPost = {
        id: this.id,
        title: this.title,
        body: this.body,
        author: {
          id: this.authorID
        }
      };
      this.$emit("updatePost", updatedPost); //emit up to parent view
      this.submitted = false;
      this.editing = false;
    },
    deletePost: function() {
      this.submitted = true;
      this.pop()
        .then(res => {
          //if resolved is true then user has clicked confirm deletion
          if (res == true) {
            let postID = this.id;
            this.$emit("deletePost", postID); //emit up to parent view
            this.submitted = false;
          } else {
            this.submitted = false;
            return;
          }
        })
        .catch(err => {
          //if error then resolved is false so user has clicked cancel deletion
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error selecting";
          } else {
            this.apiError = "error selecting";
          }
        })
        .finally(() => {
          this.password = "";
          this.submitted = false;
          this.pending = false;
        });
    },
    pop() {
      //pop presents the modal to confirm or cancel deletion requests
      this.isModalVisible = true;
      //pop returns a promise and sets the local resolved variable to the resolve of this promise
      //this lets us wait for an answer from the modal before proceeding
      return new Promise((resolve, reject) => {
        this.resolved = resolve;
      });
    },
    closeModal() {
      //if the user selects close then close the modal and set this.resolved to false
      this.isModalVisible = false;
      this.resolved(false);
    },
    confirmDeletion() {
      //if the user confirms deletion then close the modal and set this.resolved to true
      this.isModalVisible = false;
      this.resolved(true);
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
textarea {
  box-sizing: border-box;
  width: 70%;
}
</style>
