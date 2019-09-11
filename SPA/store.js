import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";

//use store
Vue.use(Vuex);

//define store
const store = new Vuex.Store({
  state: {
    userID: ""
  },
  //createPersistedState ensures state persists across page reloads
  plugins: [createPersistedState()],
  getters: {
    getUser: state => {
      return state.userID;
    }
  },
  mutations: {
    updateUser(state, { userID }) {
      state.userID = userID;
    }
  },
  actions: {
    //throughout the SPA call actions for asynchronous operations (see https://vuex.vuejs.org/guide/actions.html).
    //actions call mutations call getters call state.
    updateUser({ commit }, userID) {
      commit("updateUser", userID);
    }
  }
});

export default store;
