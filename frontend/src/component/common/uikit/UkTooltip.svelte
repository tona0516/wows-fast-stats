<script lang="ts">
  // @ts-nocheck
  export let tooltip: string | undefined;
</script>

<div uk-tooltip={tooltip}>
  <slot />
</div>
