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

  const shipTypes: { [key: string]: string } = {
    SS: "bg-primary",
    DD: "bg-info",
    CL: "bg-success",
    BB: "bg-warning",
    CV: "bg-danger",
  };

  const exampleShipTypes: { [key: string]: number } = {
    SS: 10,
    DD: 20,
    CL: 40,
    BB: 20,
    CV: 10,
  };

  const tierRates: { [key: string]: string } = {
    "Tier 1 ~ 4": "bg-primary",
    "Tier 5 ~ 7": "bg-info",
    "Tier 8 ~ ★": "bg-success",
  };

  const exampleTierRate: { [key: string]: number } = {
    "Tier 1 ~ 4": 10,
    "Tier 5 ~ 7": 30,
    "Tier 8 ~ ★": 60,
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

  <!-- ship type rate -->
  <div class="accordion" id="ShipTypeRate">
    <div class="accordion-item">
      <h2 class="accordion-header" id="headingShipTypeRate">
        <button
          class="accordion-button"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#collapseShipTypeRate"
          aria-expanded="true"
          aria-controls="collapseShipTypeRate"
        >
          艦種別プレイ割合について
        </button>
      </h2>
      <div
        id="collapseShipTypeRate"
        class="accordion-collapse collapse"
        aria-labelledby="headingShipTypeRate"
        data-bs-parent="#ShipTypeRate"
      >
        <div class="accordion-body">
          <p>そのプレイヤーの艦種ごとのプレイ率を表しています。</p>
          <table class="table help-table">
            <tbody>
              {#each Object.entries(shipTypes) as [k, v]}
                <tr class={k}>
                  <td>{k}</td>
                  <td>
                    <div class="progress">
                      <div
                        class="progress-bar progress-bar-striped {v}"
                        role="progressbar"
                        style="width: 100%"
                      />
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
          <p>以下のプレイ率の場合、バーの表示は以下のようになります。</p>
          <div class="progress">
            {#each Object.entries(exampleShipTypes) as [k, v]}
              <div
                class="progress-bar progress-bar-striped {shipTypes[k]}"
                role="progressbar"
                style="width: {v}%"
              />
            {/each}
          </div>
          <ul>
            {#each Object.entries(exampleShipTypes) as [k, v]}
              <li>{k}: {v}%</li>
            {/each}
          </ul>
        </div>
      </div>
    </div>
  </div>

  <!-- tier rate -->
  <div class="accordion" id="TierRate">
    <div class="accordion-item">
      <h2 class="accordion-header" id="headingTierRate">
        <button
          class="accordion-button"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#collapseTierRate"
          aria-expanded="true"
          aria-controls="collapseTierRate"
        >
          ティア別プレイ割合について
        </button>
      </h2>
      <div
        id="collapseTierRate"
        class="accordion-collapse collapse"
        aria-labelledby="headingTierRate"
        data-bs-parent="#TierRate"
      >
        <div class="accordion-body">
          <p>そのプレイヤーの低・中・高ティアごとのプレイ率を表しています。</p>
          <table class="table help-table">
            <tbody>
              {#each Object.entries(tierRates) as [k, v]}
                <tr class={k}>
                  <td>{k}</td>
                  <td>
                    <div class="progress">
                      <div
                        class="progress-bar progress-bar-striped {v}"
                        role="progressbar"
                        style="width: 100%"
                      />
                    </div>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
          <p>以下のプレイ率の場合、バーの表示は以下のようになります。</p>
          <div class="progress">
            {#each Object.entries(exampleTierRate) as [k, v]}
              <div
                class="progress-bar progress-bar-striped {tierRates[k]}"
                role="progressbar"
                style="width: {v}%"
              />
            {/each}
          </div>
          <ul>
            {#each Object.entries(exampleTierRate) as [k, v]}
              <li>{k}: {v}%</li>
            {/each}
          </ul>
        </div>
      </div>
    </div>
  </div>
</div>
