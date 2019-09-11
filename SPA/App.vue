<template>
  <div id="app">
    <div id="nav">
      <!-- Visitors can use this form to search users -->
      <form v-on:submit.prevent="search">
        <input v-model.trim="q" placeholder="Search..."></input>
        <button type="submit" v-on:submit.prevent="search" v-on:keyup.enter="search">Search</button>
      </form>
      <router-link to="/">Home</router-link> |
      <router-link to="/about">About</router-link> |

      <!-- if users are logged in then present this menu-->
      <template v-if="getUser && getCookie">
        <router-link to="/submit">Submit</router-link> |
        <router-link :to="{ path: `/profile/${getUser}` }">Profile</router-link>
        | <router-link to="/logout">Logout</router-link> |
      </template>
      <!-- otherwise present this menu -->
      <template v-else>
        <router-link to="/signup">Signup</router-link> |
        <router-link to="/login">Login</router-link>
      </template>
    </div>
    <router-view />
  </div>
</template>

<script>
export default {
  data: function() {
    return {
      //q is our search query
      q: "",
    }
  },
    computed: {
      //Vuex updates the store on log in. Watch for the user information
      getUser() {
        return this.$store.getters.getUser;
      },
      //There should also be a cookie when logged in. Token-hp is our non-HttpOnly cookie for javascript to access.
      getCookie() {
        return this.$cookie.get("token-hp");
      }
    },
  methods: {
    search() {
      //on search, use the router to push to the search view, passing the query paramater.
      this.$router.push({ name: "search", query: { q: this.q} })
    },
  },
   watch: {
    "$route" (to) {
      //if routing to new view and the destination has a query paramater (i.e., a new search) then update q.
      if (to.query.q) {
        this.q=to.query.q
      //otherwise (i.e., a page other than search) reset q to refresh the search bar to blank.
      } else {
        this.q=""
      }
    }
  }
};
</script>

<style>
/*
Default Vuejs CSS generated when creating a new project using Vue CLI.
*/

#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}


</style>
