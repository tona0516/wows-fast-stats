<script lang="ts">
import { get } from "svelte/store";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { storedUserConfig } from "../stores";

let userConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => (userConfig = it));
</script>

<span class="text-center"> 気に入ってもらえたらカンパお願いします🙏 </span>

<div class="centerize">
  <button
    type="button"
    class="btn btn-sm btn-primary m-1 mb-4"
    style="font-size: {userConfig.font_size};"
    on:click="{() => BrowserOpenURL('https://ofuse.me/cc0a3d3f')}"
    >開発者にカンパする
  </button>
</div>
