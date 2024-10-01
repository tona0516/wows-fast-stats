<script lang="ts">
  // @ts-nocheck
  export let name: string;
  export let onclick: any;
</script>

<!-- svelte-ignore a11y-invalid-attribute -->
<!-- svelte-ignore a11y-missing-content -->
<a href="#" class="uk-icon-button" uk-icon={name} on:click={onclick}></a>
