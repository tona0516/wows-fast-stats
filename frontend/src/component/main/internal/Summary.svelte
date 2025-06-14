<script lang="ts">
import { DispName } from "src/lib/DispName";
import type { Summary } from "src/lib/Summary";

export let summary: Summary;
</script>

<div>
  <div>艦種別平均値</div>

  <div class="grid">
    {#each summary.values.toArray() as item}
      <div>
        {DispName.SHIP_TYPE_FOR_SUMMARY.get(item.key)}

        <table class="table">
          <thead>
            <tr>
              <th />
              {#each summary.meta.headers as header}
                <th colspan={header.colspan}>{header.title}</th>
              {/each}
            </tr>
            <tr>
              <th />
              {#each summary.meta.columnNames as name}
                <th>{name}</th>
              {/each}
            </tr>
          </thead>

          <tbody>
            <tr>
              <td>味方</td>
              {#each item.value.friends as friend}
                <td>{friend}</td>
              {/each}
            </tr>

            <tr>
              <td>敵</td>
              {#each item.value.enemies as enemy}
                <td>{enemy}</td>
              {/each}
            </tr>

            <tr>
              <td>差</td>
              {#each item.value.diffs as diff}
                <td style="color: {diff.colorCode}">{diff.diff}</td>
              {/each}
            </tr>
          </tbody>
        </table>
      </div>
    {/each}
  </div>
</div>
