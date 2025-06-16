<script lang="ts">
import { BrowserOpenURL } from "wailsjs/runtime/runtime";

export let url: string;
</script>

<!-- svelte-ignore a11y-invalid-attribute -->
<a href="#" class="underline" on:click={() => BrowserOpenURL(url)}>
  <slot />
</a>
