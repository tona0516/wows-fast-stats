<script lang="ts">
  import { storedDisplayPref } from "@libs/stores";

  export let colspan: number;
  export let text: string;
</script>

<td
  class="{$storedDisplayPref.isBorderVisible
    ? 'border border-gray-500'
    : ''} px-1 py-0.5 text-center"
  {colspan}>{text}</td
>
