<script lang="ts">
  import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
  import iconApp from "./assets/images/appicon.png";
</script>

<div class="mt-3 appinfo">
  <img class="m-1" src={iconApp} alt="" width="128px" height="128px" />

  <p>
    wows-fast-stats
    <br />
    <!-- TOOD use https://github.com/Masterminds/semver -->
    バージョン 0.0.2
  </p>
  <p>
    <i class="bi bi-twitter" />
    <!-- svelte-ignore a11y-invalid-attribute -->
    <a
      class="td-link"
      href="#"
      on:click={() => BrowserOpenURL("https://twitter.com/tonango_0516")}
    >
      @tonango_0516
    </a>
    <br />
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
  </p>
  <p>Copyright © 2023 tona0516 All Rights Reserved.</p>
</div>
