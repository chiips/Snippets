<template>
  <div class="profile">
    <form v-on:submit.prevent>
      <!-- if there is a new avatar to preview -->
      <div v-if="avatarPreview">
        <img :src="avatarPreview" />
      </div>
      <!-- otherwise show existing photo -->
      <div v-else>
        <!-- if custom photo -->
        <div v-if="profileprop.avatar != 'sailboat.jpg'">
          <img
            v-bind:src="'/api/private/assets/' + profileprop.id + '/' + avatar"
          />
        </div>
        <!-- otherwise default sailboat -->
        <div v-else>
          <img src="../assets/sailboat.jpg" />
        </div>
      </div>

      <div class="error" v-if="apiError">
        {{ apiError }}
      </div>

      <button @click="launchFilePicker">Select Image</button>

      <p v-if="uploadError">{{ uploadError }}</p>

      <!-- input for new photo -->
      <input
        type="file"
        id="file"
        ref="file"
        accept="image/x-png,image/png,image/jpg,image/jpeg"
        @change="handleFileUpload()"
        style="display:none"
      />
      <!-- if editing, can save -->
      <div v-if="editingPhoto">
        <button v-on:click="submitAvatar()" :disabled="uploadError || pending">
          Save
        </button>
      </div>
    </form>

    <form v-on:submit.prevent>
      <!-- NAME -->
      <div :class="{ 'is-invalid': submitted && $v.name.$error }">
        <label>Name</label>
        <input v-model.trim="$v.name.$model" :disabled="!editing"/>
      </div>
      <div v-if="submitted && !$v.name.required">
        Name is required.
      </div>
      <div v-if="submitted && !$v.name.nameChars">
        Name must be between 1 and 15 characters and contain only standard
        letters, numbers, or underscore ("_").
      </div>

      <!-- EMAIL -->
      <div :class="{ 'is-invalid': $v.email.$error }">
        <label >Email</label>
        <input v-model.trim="$v.email.$model" :disabled="!editing"
        />
      </div>
      <div v-if="submitted && !$v.email.required">
        Email is required.
      </div>
      <div v-if="submitted && !$v.email.email">
        Please enter a valid email.
      </div>

      <p>{{ created }}</p>
      <p v-if="beenUpdated">{{ updated }}</p>

      <div v-if="editing">
        <button v-on:click="updateProfileInfo()" :disabled="pending">
          Save
        </button>
        <button v-on:click="deleteAccount()" :disabled="pending">Delete</button>
      </div>
      <div v-else>
        <button v-on:click.prevent="editing = true">Edit</button>
      </div>
    </form>
    <!--MODAL-->
    <ConfirmComp
      v-show="isModalVisible"
      @confirm="confirmDeletion"
      @close="closeModal"
      v-bind:modalprop="profileprop.name"
    />
  </div>
</template>

<script>
import { required, email, helpers } from "vuelidate/lib/validators";
import ConfirmComp from "@/components/ConfirmComp.vue";

const nameChars = helpers.regex("nameChars", /^[a-zA-Z0-9_]{1,15}$/);

export default {
  name: "ProfileComp",
  props: {
    profileprop: Object,
    pending: Boolean
  },
  components: {
    ConfirmComp
  },
  data: function() {
    return {
      id: this.profileprop.id,
      name: this.profileprop.name,
      email: this.profileprop.email,
      created: new Date(this.profileprop.created).toDateString(),
      updated: new Date(this.profileprop.updated).toDateString(),
      avatar: this.profileprop.avatar,
      file: "",
      submitted: false,
      editing: false,
      editingPhoto: false,
      avatarPreview: null,
      uploadError: "",
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
    name: { required, nameChars },
    email: { required, email }
  },
  methods: {
    //launch file picker for users to select new avatar
    launchFilePicker() {
      this.$refs.file.click();
    },
    //if the user uploads a new avatar then ensure the file meets requirements before previewing it
    handleFileUpload() {
      this.editingPhoto = true;
      this.file = this.$refs.file.files[0];
      if (!this.file) {
        this.avatarPreview = "";
        this.uploadError = "no file chosen";
        return;
      }
      if (!this.file.type.match("image.*")) {
        this.avatarPreview = "";
        this.uploadError = "please select an image";
        return;
      }
      if (this.file.size > 1024 * 1024) {
        this.avatarPreview = "";
        this.uploadError =
          "Your file is too big! Please select an image under 1MB.";
        return;
      }
      this.avatarPreview = URL.createObjectURL(this.file);
      this.uploadError = "";
      return;
    },
    //if submit selected then send the file to the API
    submitAvatar() {
      let formData = new FormData();
      if (!this.file) {
        this.uploadError = "no file chosen";
        return;
      }
      formData.append("avatar", this.file);

      this.$axios
        .put(`/api/profilephoto/${this.id}`, formData, {
          headers: {
            "Content-Type": "multipart/form-data"
          }
        })
        .then(response => {
          if (response.status == 200) {
            //if success then set the avatar to the new file name generated in the API
          this.avatar = response.data.avatar;
          this.file = "";
          this.avatarPreview = "";
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error uploading";
          } else {
            this.apiError = "error uploading";
          }
        })
        .finally(() => {
          this.editingPhoto = false;
        });
    },
    updateProfileInfo() {
      this.submitted = true;

      this.$v.$touch();
      if (this.$v.$invalid) {
        return;
      }

      let now = new Date().getTime();

      let updatedProfileInfo = {
        name: this.name,
        email: this.email,
        body: this.body,
        updated: now
      };
      this.$emit("updateProfileInfo", updatedProfileInfo);
      this.submitted = false;
      this.editing = false;
    },
    deleteAccount() {
      this.submitted = true;
      this.pop()
        .then(res => {
          if (res == true) {
            this.$emit("deleteAccount", this.id);
            this.submitted = false;
          } else {
            this.submitted = false;
            return;
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error selecting";
          } else {
            this.apiError = "error selecting";
          }
        });
    },
    pop() {
      this.isModalVisible = true;
      return new Promise((resolve, reject) => {
        this.resolved = resolve;
      });
    },
    closeModal() {
      this.isModalVisible = false;
      this.resolved(false);
    },
    confirmDeletion() {
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
</style>
