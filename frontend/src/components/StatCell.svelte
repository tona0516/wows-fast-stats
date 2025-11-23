<script lang="ts">
  export let label: string;
  export let friendValue: number | string;
  export let enemyValue: number | string;
  export let diffValue: number | string;
  export let diffValueNum: number = 0;
  export let isPercentage: boolean = false;

  function getDiffColorClass(diff: number): string {
    if (diff > 0) return "text-green-600";
    if (diff < 0) return "text-red-600";
    return "text-gray-600";
  }

  const colorClass = getDiffColorClass(diffValueNum);
</script>

<div class="flex items-center gap-3 text-sm">
  <span class="w-12">{label}</span>
  <div class="flex gap-2">
    <div class="text-center w-12">
      <div class="text-xs text-gray-500">味方</div>
      <div>{friendValue}{isPercentage ? "%" : ""}</div>
    </div>
    <div class="text-center w-12">
      <div class="text-xs text-gray-500">敵</div>
      <div>{enemyValue}{isPercentage ? "%" : ""}</div>
    </div>
    <div class="text-center w-12">
      <div class="text-xs text-gray-500">差分</div>
      <div class={colorClass}>
        {diffValue}{isPercentage ? "%" : ""}
      </div>
    </div>
  </div>
</div>

<style>
  :global(.text-green-600) {
    color: #16a34a;
  }

  :global(.text-red-600) {
    color: #dc2626;
  }

  :global(.text-gray-600) {
    color: #4b5563;
  }
</style>
