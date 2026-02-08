<script lang="ts">
  export let title: string;
  export let badgeText: string = "";
</script>

<div class="shadow-md rounded-md border border-base-300 p-4">
  <div class="flex items-center mb-4">
    <span class="text-xl font-bold">{title}</span>
    {#if badgeText}
      <span class="ml-2 badge badge-error badge-outline whitespace-nowrap"
        >{badgeText}</span
      >
    {/if}
  </div>

  <slot />
</div>
