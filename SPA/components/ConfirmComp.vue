<template>
<!-- this is our confirm modal for when users try to delete posts or their account -->
  <transition name="modal-fade">
    <div class="modal-backdrop" @click="close">
      <div class="modal" role="dialog" aria-labelledby="modalTitle" aria-describedby="modalDescription">
        <header class="modal-header" id="modalTitle">
          <slot name="header">
            Confirm Deletion
          </slot>
        </header>
        <section class="modal-body" id="modalDescription">
          <slot name="body">
            Are you sure you would like to delete {{ modalprop }}?
          </slot>
        </section>
        <footer class="modal-footer">
          <slot name="footer">
            <button type="button" class="btn-close" @click="close" aria-label="Close modal">
              Cancel
            </button>

            <button type="button" class="btn-green" @click="confirm" aria-label="Confirm modal">
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

<style>
/* sample CSS styles */
.modal-backdrop {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: rgba(0, 0, 0, 0.3);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal {
  background: white;
  overflow-x: auto;
  display: flex;
  flex-direction: column;
}

.modal-header,
.modal-footer {
  padding: 15px;
  display: flex;
}

.modal-header {
  color: green;
  justify-content: space-between;
}

.modal-footer {
  justify-content: flex-end;
}

.modal-body {
  position: relative;
  padding: 20px 10px;
}

.btn-close {
  color: white;
  background: red;
}

.btn-green {
  color: white;
  background: green;
}
</style>
