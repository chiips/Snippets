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