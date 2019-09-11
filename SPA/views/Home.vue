<template>
  <div class="home">
    <!-- fill our component with our list of posts -->
    <PostsHomeComp
      v-for="post in posts"
      v-bind:key="post.id"
      v-bind:Postprop="post"
    >
    </PostsHomeComp>
    <!-- if there's an error getting the data from the API then display an error -->
    <div class="error" v-if="apiError">
      {{ apiError }}
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import PostsHomeComp from "@/components/PostsHomeComp.vue";

export default {
  name: "home",
  components: {
    PostsHomeComp
  },
  data: function() {
    return {
      posts: [],
      pending: false,
      apiError: ""
    };
  },
  methods: {
    //define our API call
    getPosts: function() {
      this.apiError = "";
      this.$axios
        .get(`/api/posts`)
        .then(response => {
          this.posts = response.data;
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
    //fetch more data from the API when users scroll to the bottom of the screen
    scroll: function() {
      window.onscroll = () => {
        let bottomOfWindow = document.documentElement.scrollTop + window.innerHeight === document.documentElement.offsetHeight;

        //if at the bottom of the screen and the API call is not currently pending
        if (bottomOfWindow && this.pending == false) {
          this.pending = true;
          this.apiError = "";

          //get the created date of the last post displayed
          //posts are displayed in reverse chronological order so our API will take this created date and retrieve earlier posts
          let prevDate = this.posts && this.posts[this.posts.length - 1] ? this.posts[this.posts.length - 1].created : "";

          this.$axios
            .get(`/api/posts?prev=${prevDate}`)
            .then(response => {
              this.posts = this.posts.concat(response.data);
            })
            .catch(err => {
              if (err.response) {
                //if too many requests error then return
                //user will have to reinitiate scroll to try again
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
            //whether successful response or error, always reset pending to false once done
            .finally(() => {
              this.pending = false;
            });
        }
      };
    }
  },
  created() {
    //on creation of this view call the API
    this.getPosts();
  },
  mounted() {
    //mount scroll function to be ready for when the user scrolls
    this.scroll();
  }
};
</script>
