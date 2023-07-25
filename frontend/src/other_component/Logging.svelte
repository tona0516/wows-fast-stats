<script lang="ts">
  import { storedLogs } from "../stores";

  export let logs: string[];

  const elementID = "logging";

  storedLogs.subscribe(() => {
    const loggingID = document.getElementById(elementID);
    if (loggingID) {
      // Note: subscribeが更新前にコールされるため更新後にスクロールさせるために少し遅延させる
      setTimeout(() => (loggingID.scrollTop = loggingID.scrollHeight), 0.5);
    }
  });
</script>

<div class="center">
  <div class="logging overflow-scroll bg-dark p-1 mx-1" id={elementID}>
    {#each logs as log}
      {log}
      <br />
    {/each}
  </div>
</div>
