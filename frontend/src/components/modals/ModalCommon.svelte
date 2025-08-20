<script lang="ts">
  export let zValue = 50;
  export let close: () => void;
</script>

<dialog class="modal modal-open z-{zValue}">
  <div class="modal-box">
    <button
      class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
      on:click={close}
      ><i class="bi bi-x-lg"></i>
    </button>
    <slot />
  </div>
  <form method="dialog" class="modal-backdrop">
    <button type="button" on:click={close}></button>
  </form>
</dialog>
