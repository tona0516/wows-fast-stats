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
        label: "SS",
        data: [player.player_stats.using_ship_type_rate.ss.toFixed(1)],
      },
      {
        label: "DD",
        data: [player.player_stats.using_ship_type_rate.dd.toFixed(1)],
      },
      {
        label: "CL",
        data: [player.player_stats.using_ship_type_rate.cl.toFixed(1)],
      },
      {
        label: "BB",
        data: [player.player_stats.using_ship_type_rate.bb.toFixed(1)],
      },
      {
        label: "CV",
        data: [player.player_stats.using_ship_type_rate.cv.toFixed(1)],
      },
    ],
  };
</script>

{#if config.displays.overall.using_ship_type_rate}
  {#if displayPattern === "full" || displayPattern === "nopr" || displayPattern === "noshipstats"}
    <Cell class="using_ship_type_rate">
      <Bar {data} {options} height={24} />
    </Cell>
  {/if}
{/if}
