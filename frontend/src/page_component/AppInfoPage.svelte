<script lang="ts">
  import { Semver } from "../../wailsjs/go/main/App";
  import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
  import iconApp from "../assets/images/appicon.png";
</script>

<div class="mt-3 center">
  <img class="m-1" src={iconApp} alt="" width="128px" height="128px" />

  <div class="m-1 center">
    <div>wows-fast-stats</div>

    {#await Semver() then semver}
      <div>バージョン {semver}</div>
    {/await}
  </div>

  <div class="m-1 center">
    <div>
      <i class="bi bi-question-circle" />
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a
        class="td-link"
        href="#"
        on:click={() =>
          BrowserOpenURL(
            "https://github.com/tona0516/wows-fast-stats/wiki/FAQ"
          )}
      >
        FAQ
      </a>
    </div>

    <div>
      <i class="bi bi-twitter" />
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a
        class="td-link"
        href="#"
        on:click={() => BrowserOpenURL("https://twitter.com/tonango_0516")}
      >
        @tonango_0516
      </a>
    </div>

    <div>
      <i class="bi bi-github" />
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a
        class="td-link"
        href="#"
        on:click={() =>
          BrowserOpenURL("https://github.com/tona0516/wows-fast-stats")}
      >
        wows-fast-stats
      </a>
    </div>
  </div>

  <div class="m-1">Copyright © 2023 tona0516 All Rights Reserved.</div>
</div>
