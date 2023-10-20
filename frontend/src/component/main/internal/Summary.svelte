<script lang="ts">
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import { DispName } from "src/lib/DispName";
  import type { Summary } from "src/lib/Summary";

  export let summary: Summary;
</script>

<div class="uk-padding-small">
  <div>
    <div class="uk-text-bold uk-text-center">艦種別平均値</div>
  </div>

  <div class="uk-grid uk-grid-small uk-flex-center">
    {#each summary.values.toArray() as item}
      <div class="uk-margin-remove">
        <div class="uk-text-small uk-text-bold uk-text-center">
          {DispName.SHIP_TYPE_FOR_SUMMARY.get(item.key)}
        </div>

        <UkTable>
          <thead>
            <tr>
              <th />
              {#each summary.meta.headers as header}
                <th class="uk-text-center" colspan={header.colspan}
                  >{header.title}</th
                >
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
              <td class="uk-text-center">味方</td>
              {#each item.value.friends as friend}
                <td class="uk-text-center">{friend}</td>
              {/each}
            </tr>

            <tr>
              <td class="uk-text-center">敵</td>
              {#each item.value.enemies as enemy}
                <td class="uk-text-center">{enemy}</td>
              {/each}
            </tr>

            <tr>
              <td class="uk-text-center">差</td>
              {#each item.value.diffs as diff}
                <td class="uk-text-center" style="color: {diff.colorCode}"
                  >{diff.diff}</td
                >
              {/each}
            </tr>
          </tbody>
        </UkTable>
      </div>
    {/each}
  </div>
</div>
