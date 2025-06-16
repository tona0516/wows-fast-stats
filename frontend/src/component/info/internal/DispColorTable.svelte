<script lang="ts">
import type { DispColorTableInfo } from "./DispColorTableInfo";

export let caption: string;
export let tableInfo: DispColorTableInfo;
</script>

<table class="table">
  <caption>{caption}</caption>
  <thead>
    {#each tableInfo.headers as header}
      <th>{header}</th>
    {/each}
  </thead>

  <tbody>
    {#each tableInfo.rows as row}
      <tr>
        {#each row as cell}
          <td style="color: {cell.textColor};"
            >{cell.text}</td
          >
        {/each}
      </tr>
    {/each}
  </tbody>
</table>
