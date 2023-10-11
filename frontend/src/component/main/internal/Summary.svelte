<script lang="ts">
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import type { Summary } from "src/lib/Summary";

  export let summary: Summary;
</script>

<UkTable>
  <thead>
    <tr>
      <th />
      {#each summary.meta.headers as header}
        <th class="uk-text-center" colspan={header.colspan}>{header.title}</th>
      {/each}
    </tr>
    <tr>
      <th />
      {#each summary.meta.columnNames as name}
        <th class="uk-text-center">{name}</th>
      {/each}
    </tr>
  </thead>

  <tbody>
    <tr>
      <td class="uk-text-center">味方チーム平均</td>
      {#each summary.values.friends as friend}
        <td class="uk-text-center">{friend}</td>
      {/each}
    </tr>

    <tr>
      <td class="uk-text-center">敵チーム平均</td>
      {#each summary.values.enemies as enemy}
        <td class="uk-text-center">{enemy}</td>
      {/each}
    </tr>

    <tr>
      <td class="uk-text-center">差</td>
      {#each summary.values.diffs as diff}
        <td class="uk-text-center" style="color: {diff.colorCode}"
          >{diff.diff}</td
        >
      {/each}
    </tr>
  </tbody>
</UkTable>
