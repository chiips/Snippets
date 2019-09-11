import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import VueCookie from "vue-cookie";
import Vuelidate from "vuelidate";
import VueTextareaAutosize from "vue-textarea-autosize";
import axios from "axios";

//use our dependencies
Vue.use(VueCookie);
Vue.use(Vuelidate);
Vue.use(VueTextareaAutosize);

//set axios globally
Vue.prototype.$axios = axios;

//on every axios call
Vue.prototype.$axios.interceptors.response.use(
  function(response) {
    //handle access token
    //check if there's a hp-cookie.
    let token = VueCookie.get("token-hp");
    if (token) {
      //decode the payload and extract the user id.
      let id = JSON.parse(window.atob(token.split(".")[1])).id;
      //update the store.
      store.dispatch("updateUser", {
        userID: id
      });
    } else {
      //otherwise reset the store
      store.dispatch("updateUser", {
        userID: ""
      });
    }

    //handle csrf-token.
    //if not nil, set it is as default for all following responses
    let csrfToken = response.headers["x-csrf-token"];
    if (csrfToken != "") {
      Vue.prototype.$axios.defaults.headers.common["x-csrf-token"] = csrfToken;
    }


    //generate random string
    let requestID =
    Math.random()
      .toString(36)
      .substring(2, 15) +
    Math.random()
      .toString(36)
      .substring(2, 15);

    //set request id
    //API will expect custom request id as additional security
    //Logs will be able to trace requests from SPA to API
    Vue.prototype.$axios.defaults.headers.common["X-REQUEST-ID"] = requestID;

    return response;
  },
  //do the same if response is an error
  function(error) {
    let token = VueCookie.get("token-hp");
    if (token) {
      let id = JSON.parse(window.atob(token.split(".")[1])).id;
      store.dispatch("updateUser", {
        userID: id
      });
    } else {
      store.dispatch("updateUser", {
        userID: ""
      });
    }

    let csrfToken = error.response.headers["x-csrf-token"];
    if (csrfToken != "") {
      Vue.prototype.$axios.defaults.headers.common["x-csrf-token"] = csrfToken;
    }

    //generate random string
    let requestID =
    Math.random()
      .toString(36)
      .substring(2, 15) +
    Math.random()
      .toString(36)
      .substring(2, 15);

    //set request id
    //API will expect custom request id as additional security
    //Logs will be able to trace requests from SPA to API
    Vue.prototype.$axios.defaults.headers.common["X-REQUEST-ID"] = requestID;
  }
);

Vue.config.productionTip = false;

//our Vue instance
new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
