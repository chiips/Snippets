<template>
<!-- this is our confirm modal for when users try to delete posts or their account -->
  <transition name="modal-fade">
    <div id="modalBackground" @click="close">
      <div id="modal" role="dialog" aria-labelledby="modalTitle" aria-describedby="modalDescription">
        <header id="modalTitle">
          <slot name="header">
            Confirm Deletion
          </slot>
        </header>
        <section id="modalDescription">
          <slot name="body">
            Are you sure you would like to delete {{ modalprop }}?
          </slot>
        </section>
        <footer>
          <slot name="footer">
            <button type="button" @click="close" aria-label="Close modal">
              Cancel
            </button>
            <button type="button" @click="confirm" aria-label="Confirm modal">
              Delete
            </button>
          </slot>
        </footer>
      </div>
    </div>
  </transition>
</template>

<script>
export default {
  name: "modal",
  props: {
    modalprop: String
  },
  methods: {
    //on confirm or close emit the response up to view importing this component
    confirm() {
      this.$emit("confirm");
    },
    close() {
      this.$emit("close");
    }
  }
};
</script>
