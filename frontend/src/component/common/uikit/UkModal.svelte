<script lang="ts">
// @ts-nocheck
export let id: string;
</script>

<div {id} class="uk-flex-top uk-modal">
  <div class="uk-modal-dialog">
    <button class="uk-modal-close-default" type="button" uk-close></button>

    <div class="uk-modal-body">
      <slot name="body" />
    </div>

    <div class="uk-modal-footer">
      <slot name="footer" />
    </div>
  </div>
</div>
