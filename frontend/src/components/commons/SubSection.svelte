<script lang="ts">
  export let title: string;
</script>

<div class="bg-base-200 shadow-md rounded-md border border-base-300 p-4">
  <span class="font-semibold">{title}</span>
  <slot />
</div>
