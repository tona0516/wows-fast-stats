<script lang="ts">
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { storedUserConfig } from "../stores";
</script>

<div class="center">
  <span> 気に入ってもらえたらカンパお願いします🙏 </span>

  <button
    type="button"
    class="btn btn-sm btn-primary m-1 mb-4"
    style="font-size: {$storedUserConfig.font_size};"
    on:click="{() => BrowserOpenURL('https://ofuse.me/cc0a3d3f')}"
    >開発者にカンパする
  </button>
</div>
