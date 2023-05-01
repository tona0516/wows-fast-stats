<script lang="ts">
import { BrowserOpenURL } from "../wailsjs/runtime/runtime";

const prs: { [key: string]: string } = {
  bad: "0 ~ 750",
  "below-average": "750 ~ 1100",
  average: "1100 ~ 1350",
  good: "1350 ~ 1550",
  "very-good": "1550 ~ 1750",
  great: "1750 ~ 2100",
  unicum: "2100 ~ 2450",
  "super-unicum": "2450 ~ 9999",
};
</script>

<p>
  <!-- svelte-ignore a11y-invalid-attribute -->
  <a
    href="#"
    on:click="{() =>
      BrowserOpenURL('https://asia.wows-numbers.com/personal/rating')}"
    >WoWS Stats & Numbers - Personal Rating <i class="bi bi-box-arrow-up-right"
    ></i></a
  >
  に基づいています。
  <br />
  このアプリでは以下のような配色となっています。
</p>
<table class="table table-text-color help-table">
  <tbody>
    {#each Object.entries(prs) as [k, v]}
      <tr class="{k}">
        <td>{k}</td>
        <td>{v}</td>
      </tr>
    {/each}
  </tbody>
</table>
