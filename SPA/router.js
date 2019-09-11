import Vue from "vue";
import Router from "vue-router";
import Home from "./views/Home.vue";
import store from "./store.js";

//user router
Vue.use(Router);

//define router
const router = new Router({
  //use history mode for URL navigation without page reloads (see https://router.vuejs.org/guide/essentials/history-mode.html)
  mode: "history",
  base: process.env.BASE_URL,
  //define our routes
  routes: [
    {
      path: "/",
      name: "home",
      component: Home,
      //note specifically which routes are public, making routes private by default.
      meta: { public: true }
    },
    {
      path: "/about",
      name: "about",
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import("./views/About.vue"),
      meta: { public: true }
    },
    {
      path: "/signup",
      name: "signup",
      component: () => import("./views/Signup.vue"),
      meta: { public: true }
    },
    {
      path: "/login",
      name: "login",
      component: () => import("./views/Login.vue"),
      meta: {
        public: true,
        //prevent accessing the login page if users are already logged in
        onlyWhenLoggedOut: true
      }
    },
    {
      path: "/logout",
      name: "logout",
      component: () => import("./views/Logout.vue")
    },
    {
      path: "/search",
      name: "search",
      component: () => import("./views/Search.vue"),
      meta: { public: true }
    },
    {
      path: "/author/:id",
      name: "author",
      component: () => import("./views/Author.vue"),
      props: true,
      meta: { public: true }
    },
    {
      path: "/profile/:id",
      name: "profile",
      component: () => import("./views/Profile.vue")
    },
    {
      path: "/submit",
      name: "submit",
      component: () => import("./views/Submit.vue")
    },
    {
      path: "*", //should go at end
      name: "catch-all",
      component: () => import("./views/NotFound.vue"),
      meta: { public: true }
    }
  ]
});

router.beforeEach((to, from, next) => {
  const isPublic = to.matched.some(record => record.meta.public);
  const onlyWhenLoggedOut = to.matched.some(record => record.meta.onlyWhenLoggedOut);
  //count user as logged out if either store empty or cookie not present.
  const loggedOut = !store.getters.getUser || !Vue.cookie.get("token-hp");

  //if not a public page and you're logged out
  if (!isPublic && loggedOut) {
    //ensure both cookie and store are clear
    Vue.cookie.delete("token-hp");
    store
      .dispatch("updateUser", { userID: "" })
      .then(() => {
        return next({
          //redirect to the login page.
          //since the navbar in App.vue adapts to user state, this should not be a problem.
          //users will only be redirected like this if they type a private url into the browser while logged out.
          path: "/login",
          query: { redirect: to.fullPath }
        });
      })
      .catch(() => {
        //even if there's an error in dispatching the store update, redirect to login
        return next({
          path: "/login",
          query: { redirect: to.fullPath }
        });
      });
  }

  //if only when logged out page (i.e., /login) and you are already logged in
  if (onlyWhenLoggedOut && !loggedOut) {
    //redirect to home
    //again since the navbar in App.vue adapts to user state, this should not be a problem.
    //users will only be redirected like this if they type the login url into the browser while already logged in.
    return next("/");
  }

  //otherwise go next
  next();
});

export default router;
