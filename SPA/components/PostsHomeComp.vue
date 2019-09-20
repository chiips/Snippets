<template>
  <div>
    <h3>{{ postprop.title }}</h3>
    <p>{{ postprop.body }}</p>
    <p>{{ created }}</p>
    <p v-if="beenUpdated">{{ updated }}</p>
    <router-link :to="{ name: 'author', params: { id: authorID } }">
      {{ authorName }}
    </router-link>
    <!-- if their avatar has been updated from the default sailboat.jpg -->
    <div v-if="this.avatar != 'sailboat.jpg'">
      <img v-bind:src="'/api/private/assets/' + authorID + '/' + avatar" />
    </div>
    <!-- otherwise display the default -->
    <div v-else>
      <img src="../assets/sailboat.jpg" />
    </div>
  </div>
</template>

<script>
export default {
  name: "PostsHomeComp",
  props: {
    postprop: Object
  },
  data: function() {
    return {
      created: new Date(this.postprop.created).toDateString(),
      updated: new Date(this.postprop.updated).toDateString(),
      authorID: this.postprop && this.postprop.author ? this.postprop.author.id : "",
      authorName: this.postprop && this.postprop.author ? this.postprop.author.name : "",
      avatar: this.postprop && this.postprop.author ? this.postprop.author.avatar : ""
    };
  },
  computed: {
    beenUpdated() {
      if (this.updated === this.created) {
        return false;
      }
      return true;
    }
  }
};
</script>