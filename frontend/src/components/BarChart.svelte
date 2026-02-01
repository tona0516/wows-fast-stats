<script lang="ts">
  export let data: {
    label: string;
    friendValue: number;
    enemyValue: number;
  }[] = [];

  function formatValue(value: number): string {
    if (value === 0) return "0";
    if (value >= 1000) {
      return Math.round(value).toLocaleString();
    }
    return value.toFixed(1);
  }

  $: chartData = data.map((item) => {
    const maxValue = Math.max(item.friendValue, item.enemyValue);
    const diff = item.friendValue - item.enemyValue;
    const isPositive = diff > 0;

    const friendSize = maxValue > 0 ? item.friendValue / maxValue : 0;
    const enemySize = maxValue > 0 ? item.enemyValue / maxValue : 0;

    return {
      label: item.label,
      friendValue: item.friendValue,
      enemyValue: item.enemyValue,
      friendSize,
      enemySize,
      diff,
      isPositive,
    };
  });
</script>

<div class="mb-2">
  {#each chartData as item}
    <div class="mb-2 last:mb-0">
      <div class="flex items-center justify-between">
        <span class="font-semibold text-base-content/80 text-sm"
          >{item.label}</span
        >
        <span
          class="font-bold text-sm"
          class:text-success={item.isPositive}
          class:text-error={!item.isPositive}
        >
          {item.isPositive ? "+" : ""}{formatValue(item.diff)}
        </span>
      </div>
      <div class="space-y-1">
        <div class="bar-row">
          <div class="bar-track">
            <div
              class="bar-fill bg-ally"
              style="width: {item.friendSize * 100}%;"
            ></div>
          </div>
          <span class="ml-2 bar-value text-sm"
            >{formatValue(item.friendValue)}</span
          >
        </div>
        <div class="bar-row">
          <div class="bar-track">
            <div
              class="bar-fill bg-enemy"
              style="width: {item.enemySize * 100}%;"
            ></div>
          </div>
          <span class="ml-2 bar-value text-sm"
            >{formatValue(item.enemyValue)}</span
          >
        </div>
      </div>
    </div>
  {/each}
</div>

<style>
  .bar-row {
    display: grid;
    grid-template-columns: 1fr auto;
    align-items: center;
    gap: 0.375rem;
  }

  .bar-track {
    position: relative;
    height: 0.5rem;
    border-radius: 9999px;
    background: rgba(148, 163, 184, 0.25);
    overflow: hidden;
  }

  .bar-fill {
    position: absolute;
    inset: 0;
    width: 0;
    border-radius: 9999px;
  }

  .bar-value {
    line-height: 1;
    color: hsl(var(--bc) / 0.9);
  }
</style>
