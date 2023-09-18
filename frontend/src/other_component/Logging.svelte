<script lang="ts">
  import { storedLogs } from "src/stores";

  export let logs: string[];

  const elementID = "logging";

  storedLogs.subscribe(() => {
    const loggingID = document.getElementById(elementID);
    if (loggingID) {
      // Note: subscribeが更新前にコールされるため更新後にスクロールさせるために少し遅延させる
      setTimeout(() => (loggingID.scrollTop = loggingID.scrollHeight), 100);
    }
  });
</script>

<div class="overflow-scroll bg-dark p-1" id={elementID}>
  {#each logs as log}
    {log}
    <br />
  {/each}
</div>
