<script lang="ts">
// @ts-nocheck
export let id: string;
</script>

<div {id} class="uk-flex-top uk-modal">
  <div class="uk-modal-dialog uk-margin-auto-vertical uk-background-secondary">
    <button class="uk-modal-close-default" type="button" uk-close></button>

    <div class="uk-modal-body">
      <slot name="body" />
    </div>

    <div class="uk-modal-footer uk-text-right uk-background-secondary">
      <slot name="footer" />
    </div>
  </div>
</div>
