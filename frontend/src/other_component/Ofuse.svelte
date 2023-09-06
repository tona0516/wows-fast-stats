<script lang="ts">
  import { storedUserConfig } from "src/stores";
  import { Button } from "sveltestrap";
  import { BrowserOpenURL } from "wailsjs/runtime/runtime";
</script>

<div class="center mt-2 mx-2">
  <span> 気に入ってもらえたらカンパお願いします🙏 </span>

  <Button
    size="sm"
    color="primary"
    class="m-1 mb-4"
    on:click={() => BrowserOpenURL("https://ofuse.me/cc0a3d3f")}
    >開発者にカンパする
  </Button>
</div>
