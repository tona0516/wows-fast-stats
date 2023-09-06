<script lang="ts">
  import { BrowserOpenURL } from "wailsjs/runtime/runtime";

  export let url: string;
  export let text: string;
</script>

<!-- svelte-ignore a11y-invalid-attribute -->
<a class="td-link" href="#" on:click={() => BrowserOpenURL(url)}>{text}</a>
