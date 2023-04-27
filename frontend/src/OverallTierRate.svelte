<script lang="ts">
  import type { vo } from "wailsjs/go/models";
  import { Bar } from "svelte-chartjs";
  // Note: need to import for component registration.
  import Chart from "chart.js/auto";
  import { LogDebug } from "../wailsjs/runtime/runtime";
  import { Cell } from "@smui/data-table";
  export let config: vo.UserConfig;
  export let player: vo.Player;
  export let displayPattern: DisplayPattern;

  const options = {
    responsive: true,
    plugins: {
      legend: {
        display: false,
      },
    },
    scales: {
      x: {
        stacked: true,
        display: false,
        max: 100,
      },
      y: {
        stacked: true,
        display: false,
        max: 100,
      },
    },
    indexAxis: "y",
  };

  const data = {
    labels: [""],
    datasets: [
      {
        label: "1~4",
        data: [player.player_stats.using_tier_rate.low.toFixed(1)],
      },
      {
        label: "5~7",
        data: [player.player_stats.using_tier_rate.middle.toFixed(1)],
      },
      {
        label: "9~★",
        data: [player.player_stats.using_tier_rate.high.toFixed(1)],
      },
    ],
  };
</script>

<!-- using tier rate -->
{#if config.displays.overall.using_tier_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <Cell class="using_tier_rate">
      <Bar {data} {options} height={24} />
    </Cell>
  {/if}
{/if}
