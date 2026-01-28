<script lang="ts">
  import "charts.css";

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

<div class="space-y-4">
  {#each chartData as item}
    <div class="space-y-1">
      <div class="flex items-center justify-between text-sm">
        <span class="font-semibold text-base-content/80">{item.label}</span>
        <span
          class="font-bold"
          class:text-success={item.isPositive}
          class:text-error={!item.isPositive}
        >
          {item.isPositive ? "+" : ""}{formatValue(item.diff)}
        </span>
      </div>
      <table class="charts-css column data-end multiple">
        <tbody>
          <tr>
            <th scope="row"></th>
            <td style="--size: {item.friendSize}">
              {#if item.friendSize < 0.5}
                <span class="text-sm data outside">
                  {formatValue(item.friendValue)}
                </span>
              {:else}
                <span class="text-sm text-gray-900">
                  {formatValue(item.friendValue)}
                </span>
              {/if}
            </td>
            <td style="--size: {item.enemySize}">
              {#if item.enemySize < 0.5}
                <span class="text-sm data outside">
                  {formatValue(item.enemyValue)}
                </span>
              {:else}
                <span class="text-sm text-gray-900">
                  {formatValue(item.enemyValue)}
                </span>
              {/if}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  {/each}
</div>

<style>
  :global(.charts-css.column) {
    --color-1: rgb(77, 233, 170);
    --color-2: rgb(255, 76, 42);
  }
</style>
