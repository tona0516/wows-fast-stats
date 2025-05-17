<script lang="ts">
import { storedLogs } from "src/stores";

const elementID = "logging";

storedLogs.subscribe(() => {
  const elem = document.getElementById(elementID);
  if (elem) {
    // Note: subscribeが更新前にコールされるため更新後にスクロールさせるために少し遅延させる
    setTimeout(() => {
      elem.scrollTop = elem.scrollHeight;
    }, 100);
  }
});
</script>

<div class="uk-padding-small">
  <div
    class="overflow-scroll uk-text-emphasis uk-tile-secondary"
    id={elementID}
  >
    {#each $storedLogs as log}
      {log}
      <br />
    {/each}
  </div>
</div>
