<template>
  <div class="profile">
    <!-- pass data into our component as before.
    this time include update functions that take the event which will be emitted up from within the component -->
    <ProfileComp
      v-bind:key="currentUser.id"
      v-bind:profileprop="currentUser"
      v-bind:pending="pending"
      v-on:updateProfileInfo="onUpdateProfileInfo($event)"
      v-on:deleteAccount="onDeleteAccount($event)"
    >
    </ProfileComp>
    <div class="error" v-if="apiError">
      {{ apiError }}
    </div>
    <hr />
    <!-- for these update functions also take the index so we can update the data instantly on successful API calls -->
    <PostsProfileComp
      v-for="(post, index) in posts"
      v-bind:key="post.id"
      v-bind:Postprop="post"
      v-bind:pending="pending"
      v-on:updatePost="onUpdatePost($event, index)"
      v-on:deletePost="onDeletePost($event, index)"
    >
    </PostsProfileComp>
    <div class="error" v-if="apiError">
      {{ apiError }}
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import ProfileComp from "@/components/ProfileComp.vue";
import PostsProfileComp from "@/components/PostsProfileComp.vue";

export default {
  name: "home",
  components: {
    ProfileComp,
    PostsProfileComp
  },
  data: function() {
    return {
      currentUser: {},
      posts: [],
      pending: false, //differentiate pending and loading for delete and scroll
      loading: false,
      apiError: ""
    };
  },
  computed: {
    id() {
      //watch the id of the user in the url
      return this.$route.params.id;
    }
  },
  methods: {
    getProfileInfo: function(id) {
      this.apiError = "";
      this.$axios
        .get(`/api/profile/${id}`)
        .then(response => {
          if (response.status == 200) {
            this.currentUser = response.data;
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error fetching data";
          } else {
            this.apiError = "error fetching data";
          }
        });
    },
    getProfilePosts: function(id) {
      this.apiError = "";
      this.$axios
        .get(`/api/posts/${id}`)
        .then(response => {
          if (response.status == 200) {
            this.posts = response.data;
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error fetching data";
          } else {
            this.apiError = "error fetching data";
          }
        });
    },
    scroll: function(id) {
      window.onscroll = () => {
        let bottomOfWindow = document.documentElement.scrollTop + window.innerHeight === document.documentElement.offsetHeight;

        if (bottomOfWindow && this.loading == false) {
          this.loading = true;
          this.apiError = "";

          let prevDate = this.posts && this.posts[this.posts.length - 1] ? this.posts[this.posts.length - 1].created : "";

          this.$axios
            .get(`/api/posts/${id}?prev=${prevDate}`)
            .then(response => {
              if (response.status == 200) {
                this.posts = this.posts.concat(response.data);
              }
            })
            .catch(err => {
              if (err.response) {
                if (err.response.status == 429) {
                  return;
                }
                this.apiError = err.response.data;
              } else if (err.request) {
                this.apiError = "error fetching data";
              } else {
                this.apiError = "error fetching data";
              }
            })
            .finally(() => {
              this.loading = false;
            });
        }
      };
    },
    onUpdateProfileInfo: function(updatedProfileInfo) {
      this.pending = true;

      //get the updated information
      let name = updatedProfileInfo.name;
      let email = updatedProfileInfo.email;
      let updated = updatedProfileInfo.updated;

      let id = this.currentUser.id;

      //return if no changes made
      if (name === this.currentUser.name && email === this.currentUser.email) {
        this.pending = false;
        return;
      }

      this.apiError = "";

      this.$axios
        .put(`/api/profile/${id}`, {
          name: name,
          email: email
        })
        .then(response => {
          if (response.status == 200) {
            //if successful API response then update our variables locally for quick changes not requiring page reload
            //we also save requiring the API to send back redundant data
            this.currentUser.name = name;
            this.currentUser.email = email;
            this.currentUser.updated = updated;
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error updating";
          } else {
            this.apiError = "error updating";
          }
        })
        .finally(() => {
          this.pending = false;
        });
    },
    onDeleteAccount: function(id) {
      this.pending = true;
      this.apiError = "";

      this.$axios
        .delete(`/api/profile/${id}`)
        .then(response => {
          if (response.status == 200) {
            //if successful API response then log the user out and redirect to home
            this.$store
              .dispatch("updateUser", {
                userID: ""
              })
              .then(() => {
                this.$router.push({ name: "home" });
              })
              .catch(err => {
                if (err.response) {
                  this.apiError = err.response.data;
                } else if (err.request) {
                  this.apiError = "error deleting";
                } else {
                  this.apiError = "error deleting";
                }
              });
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error deleting";
          } else {
            this.apiError = "error deleting";
          }
        })
        .finally(() => {
          this.pending = false;
        });
      return;
    },
    onUpdatePost: function(updatedPost, index) {
      this.pending = true;

      //return if no changes made
      if (updatedPost.title === this.posts[index].title && updatedPost.body === this.posts[index].body) {
        this.pending = false;
        return;
      }

      this.apiError = "";
      this.$axios
        .put(`/api/post`, updatedPost)
        .then(response => {
          if (response.status == 200) {
            //on success update the post locally
            this.$set(this.posts, index, updatedPost);
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error updating";
          } else {
            this.apiError = "error updating";
          }
        })
        .finally(() => {
          this.pending = false;
        });
    },
    onDeletePost: function(PostID, index) {
      this.pending = true;
      this.apiError = "";

      this.$axios
        .delete(`/api/post/${PostID}`)
        .then(response => {
          if (response.status == 200) {
            //on success delete the post locally
            this.$delete(this.posts, index);
          }
        })
        .catch(err => {
          if (err.response) {
            this.apiError = err.response.data;
          } else if (err.request) {
            this.apiError = "error deleting";
          } else {
            this.apiError = "error deleting";
          }
        })
        .finally(() => {
          this.pending = false;
        });
      return;
    }
  },
  created() {
    this.getProfileInfo(this.id);
    this.getProfilePosts(this.id);
  },
  mounted() {
    this.scroll(this.id);
  }
};
</script>
