<script lang="ts">
  import SvelteTable from "svelte-table";
  import { Rows } from "../wailsjs/go/main/App.js";

  let rows = [];
  const columns = [
  {
    key: "name",
    title: "Name",
    value: v => v.name,
  },
  {
    key: "id",
    title: "ID",
    value: v => v.id,
  },
];

function generateRows(): Promise<any> {
    return Rows().then((result) => rows = result);
}

generateRows()

</script>

<main>
  <SvelteTable columns="{columns}" rows="{rows}"></SvelteTable>
</main>

<style>
</style>
