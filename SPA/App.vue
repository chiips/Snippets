<template>
  <div id="app" class="has-navbar-fixed-top" >
    <div class="navbar is-fixed-top is-transparent" role="navigation" aria-label="main navigation">

      <div class="navbar-brand">
        <div @click="closeMenu()">
          <router-link to="/" class="navbar-item" ><img src="@/assets/logo.png"></router-link>
        </div>

        <div role="button" v-bind:class="[navbarBurger, isActive ? 'is-active' : '']" aria-label="menu" aria-expanded="false" data-target="navMenu" @click="toggleMenu()">
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </div>

      </div>


      <div id="navMenu" v-bind:class="[navbarMenu, isActive ? 'is-active' : '']" @click="closeMenu()">
        <div class="navbar-start">
          <form v-on:submit.prevent="search" class="navbar-item is-expanded">
            <input v-model.trim="q" class="input" type="text">
            <button type="submit" v-on:submit.prevent="search" v-on:keyup.enter="search" class="button is-primary">Search</button>
          </form>
        </div>

        <div class="navbar-end">
          <router-link to="/" class="navbar-item is-tab">Home</router-link> 
          <router-link to="/about" class="navbar-item is-tab">About</router-link> 

          <template v-if="getUser">
            <router-link to="/submit" class="navbar-item is-tab">Submit</router-link> 
            <router-link :to="{ path: `/profile/${getUser}` }" class="navbar-item is-tab">Profile</router-link> 
            <router-link to="/logout" class="navbar-item is-tab">Logout</router-link> 
          </template>
          <template v-else>
            <div class="navbar-item">
              <div class="buttons">
                <router-link to="/signup" class="button is-link">Sign up</router-link> 
                <router-link to="/login" class="button is-primary">Log in</router-link> 
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>

    <br>

    <router-view class="section"/>
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
      //Vuex updates the store on log in and out. Watch for the user information
      getUser() {
        return this.$store.getters.getUser;
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

<style>
/*
CSS styles applied to the whole SPA
*/
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.container {
  padding-top: 100px;
}

p {
  white-space: pre-wrap;
}

form {
  width: 50%;
  margin: auto;
  text-align: left;

}

</style>