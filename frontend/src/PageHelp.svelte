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

<div class="mt-3 help">
  <!-- PR color -->
  <div class="accordion" id="PrColor">
    <div class="accordion-item">
      <h2 class="accordion-header" id="headingPrColor">
        <button
          class="accordion-button"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#collapsePrColor"
          aria-expanded="true"
          aria-controls="collapsePrColor"
        >
          プレイヤーごとの配色について
        </button>
      </h2>
      <div
        id="collapsePrColor"
        class="accordion-collapse collapse"
        aria-labelledby="headingPrColor"
        data-bs-parent="#PrColor"
      >
        <div class="accordion-body">
          <p>
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
              href="#"
              on:click={() =>
                BrowserOpenURL("https://asia.wows-numbers.com/personal/rating")}
              >WoWS Stats & Numbers - Personal Rating <i
                class="bi bi-box-arrow-up-right"
              /></a
            >
            に基づいています。
            <br />
            このアプリでは以下のような配色となっています。
          </p>
          <table class="table table-text-color help-table">
            <tbody>
              {#each Object.entries(prs) as [k, v]}
                <tr class={k}>
                  <td>{k}</td>
                  <td>{v}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</div>
