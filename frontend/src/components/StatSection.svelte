<script lang="ts">
  export let title: string;
</script>

<div>
  <div class="text-xs font-semibold text-gray-600 mb-2">{title}</div>
  <div class="space-y-2">
    <slot />
  </div>
</div>

<style>
  :global(.text-gray-600) {
    color: #4b5563;
  }
</style>
