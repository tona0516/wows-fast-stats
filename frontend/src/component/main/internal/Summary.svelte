<script lang="ts">
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import type { Summary } from "src/lib/Summary";

  export let summary: Summary;

  const COLUMNS = [
    { colspan: summary.tableInfo.shipColspan, text: "艦成績" },
    { colspan: summary.tableInfo.overallColspan, text: "総合成績" },
  ];
</script>

<UkTable>
  <thead>
    <tr>
      <th />
      {#each COLUMNS as column}
        <th class="uk-text-center" colspan={column.colspan}>{column.text}</th>
      {/each}
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
        <td class="uk-text-center" style="color: {diff.color}">{diff.value}</td>
      {/each}
    </tr>
  </tbody>
</UkTable>
