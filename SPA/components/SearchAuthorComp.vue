<template>
  <div>
    <h3>
      <router-link :to="{ name: 'author', params: { id: authorID } }">
        {{ authorName }}
      </router-link>
    </h3>
    <p>{{ authorprop.email }}</p>
    <p>{{ created }}</p>

    <div v-if="this.avatar != 'sailboat.jpg'">
      <img v-bind:src="'/api/private/assets/' + authorID + '/' + avatar" />
    </div>
    <div v-else>
      <img src="../assets/sailboat.jpg" />
    </div>
  </div>
</template>

<script>
export default {
  name: "AuthorComp",
  props: {
    authorprop: Object
  },
  data: function() {
    return {
      authorName: this.authorprop.name,
      authorID: this.authorprop.id,
      avatar: this.authorprop.avatar,
      created: new Date(this.authorprop.created).toDateString()
    };
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
