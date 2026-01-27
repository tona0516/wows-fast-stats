<script lang="ts">
    import { storedDisplayPref } from "@libs/stores";

  export let colspan: number;
  export let text: string;
</script>

<td class="{$storedDisplayPref.showBoarder ? "border border-gray-500" : ""} p-1 text-center" {colspan}>{text}</td>
