<script lang="ts">
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import type { DispColorTableInfo } from "./DispColorTableInfo";

  export let caption: string;
  export let tableInfo: DispColorTableInfo;
</script>

<UkTable>
  <caption>{caption}</caption>
  <thead>
    {#each tableInfo.headers as header}
      <th class="uk-text-center">{header}</th>
    {/each}
  </thead>

  <tbody>
    {#each tableInfo.rows as row}
      <tr>
        {#each row as cell}
          <td
            class="uk-text-center"
            style="color: {cell.textColor}; background-color: {cell.bgColor}"
            >{cell.text}</td
          >
        {/each}
      </tr>
    {/each}
  </tbody>
</UkTable>
