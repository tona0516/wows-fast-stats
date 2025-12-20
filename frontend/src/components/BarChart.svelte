<script lang="ts">
  export let data: {
    label: string;
    friendValue: number;
    enemyValue: number;
  }[] = [];

  const MAX_BAR_WIDTH = 100;

  function getBarWidth(value: number, maxValue: number): number {
    if (maxValue === 0) return 0;
    return (value / maxValue) * MAX_BAR_WIDTH;
  }

  function formatValue(value: number): string {
    if (value === 0) return "0";
    if (value >= 1000) {
      return Math.round(value).toLocaleString();
    }
    return value.toFixed(1);
  }

  $: chartData = data.map((item) => {
    const maxValue = Math.max(item.friendValue, item.enemyValue);
    const friendWidth = getBarWidth(item.friendValue, maxValue);
    const enemyWidth = getBarWidth(item.enemyValue, maxValue);
    const diff = item.friendValue - item.enemyValue;
    const isPositive = diff > 0;

    return {
      label: item.label,
      friendValue: item.friendValue,
      enemyValue: item.enemyValue,
      friendWidth,
      enemyWidth,
      diff,
      isPositive,
    };
  });
</script>

<div class="w-xl grid grid-cols-2 gap-4">
  {#each chartData as item}
    <div class="flex-1 space-y-1">
      <div class="flex items-center justify-between text-sm">
        <span class="font-medium">{item.label}</span>
        <span
          class="text-xs font-mono"
          class:text-success={item.isPositive}
          class:text-error={!item.isPositive}
        >
          {item.isPositive ? "+" : ""}{formatValue(item.diff)}
        </span>
      </div>
      <div class="space-y-1">
        <!-- Friend Team Bar -->
        <div class="flex items-center gap-2">
          <span class="text-xs w-10 text-right text-primary font-medium"
            >味方</span
          >
          <div class="flex-1">
            <div class="h-4 bg-base-300 rounded">
              <div
                class="h-full bg-primary rounded transition-all duration-300"
                style="width: {item.friendWidth}%"
              />
            </div>
          </div>
          <span class="text-xs w-14 text-right font-mono text-primary">
            {formatValue(item.friendValue)}
          </span>
        </div>
        <!-- Enemy Team Bar -->
        <div class="flex items-center gap-2">
          <span class="text-xs w-10 text-right text-secondary font-medium"
            >敵</span
          >
          <div class="flex-1">
            <div class="h-4 bg-base-300 rounded">
              <div
                class="h-full bg-secondary rounded transition-all duration-300"
                style="width: {item.enemyWidth}%"
              />
            </div>
          </div>
          <span class="text-xs w-14 text-right font-mono text-secondary">
            {formatValue(item.enemyValue)}
          </span>
        </div>
      </div>
    </div>
  {/each}
</div>
