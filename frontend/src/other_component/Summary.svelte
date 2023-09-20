<script lang="ts">
  import type { Summary } from "src/lib/Summary";

  export let summary: Summary;
</script>

<div class="uk-flex uk-flex-center">
  <div class="uk-overflow-auto">
    <table
      class="uk-table uk-table-shrink uk-table-divider uk-table-small uk-table-middle uk-text-nowrap"
    >
      <thead>
        <tr>
          <th />
          <th class="uk-text-center" colspan={summary.tableInfo.shipColspan}
            >艦成績</th
          >
          <th class="uk-text-center" colspan={summary.tableInfo.overallColspan}
            >総合成績</th
          >
        </tr>
        <tr>
          <th />
          {#each summary.values.map((it) => it.label) as label}
            <th class="uk-text-center">{label}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        <tr>
          <td class="uk-text-center">味方チーム平均</td>
          {#each summary.values.map((it) => it.friend) as friend}
            <td class="uk-text-center">{friend}</td>
          {/each}
        </tr>

        <tr>
          <td class="uk-text-center">敵チーム平均</td>
          {#each summary.values.map((it) => it.enemy) as enemy}
            <td class="uk-text-center">{enemy}</td>
          {/each}
        </tr>

        <tr>
          <td class="uk-text-center">差</td>
          {#each summary.values.map((it) => it.diff) as diff}
            <td class="uk-text-center" style="color: {diff.color}"
              >{diff.value}</td
            >
          {/each}
        </tr>
      </tbody>
    </table>
  </div>
</div>
