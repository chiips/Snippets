<template>
  <div class="home">
    <SearchAuthorComp
      v-for="user in users"
      v-bind:key="user.id"
      v-bind:authorprop="user"
    >
    </SearchAuthorComp>
    <div class="error" v-if="apiError">
      {{ apiError }}
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import SearchAuthorComp from "@/components/SearchAuthorComp.vue";

export default {
  name: "search",
  components: {
    SearchAuthorComp
  },
  data: function() {
    return {
      users: [],
      pending: false,
      apiError: "",
      q: ""
    };
  },
  methods: {
    searchUsers() {
      this.apiError = "";
      this.q = this.q.trim(); //trim any spaces
      this.$axios
        .get(`/api/search?q=${this.q}`)
        .then(response => {
          this.users = response.data;
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
    scroll: function() {
      window.onscroll = () => {
        let bottomOfWindow = document.documentElement.scrollTop + window.innerHeight === document.documentElement.offsetHeight;

        if (bottomOfWindow && this.pending == false) {
          this.pending = true;
          this.apiError = "";

          let prevDate =
            this.users && this.users[this.users.length - 1] ? this.users[this.users.length - 1].created : "";

          this.q = this.q.trim();
          this.$axios
            .get(`/api/search?prev=${prevDate}&q=${this.q}`)
            .then(response => {
              this.users = this.users.concat(response.data);
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
              this.pending = false;
            });
        }
      };
    }
  },
  created() {
    //when a user visits this view, if the query parameter is present then set it to our local q variable
    if (this.$route.query.q) {
      this.q = this.$route.query.q;
    }
    this.searchUsers();
  },
  beforeRouteUpdate(to, from, next) {
    //when a user updates their search, update our q variable and call the API again
    this.q = to.query.q;
    this.searchUsers();
    next();
  },
  beforeRouteLeave(to, from, next) {
    //when a user leaves the search page, reset our q variable
    this.q = "";
    next();
  },
  mounted() {
    this.scroll();
  }
};
</script>
