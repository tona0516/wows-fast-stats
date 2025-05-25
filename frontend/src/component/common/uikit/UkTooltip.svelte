<script lang="ts">
// @ts-nocheck
export let tooltip: string | undefined;
</script>

{#if tooltip}
  <div uk-tooltip="title: {tooltip}; pos: top">
    <slot />
  </div>
{:else}
  <slot />
{/if}
