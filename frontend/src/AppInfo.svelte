<script lang="ts">
  import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
  import iconApp from "./assets/images/appicon.png";
</script>

<div class="m-2">
  <img class="m-3" src={iconApp} alt="" width="128px" height="128px" />

  <p>
    wows-fast-stats
    <br />
    <!-- TOOD use https://github.com/Masterminds/semver -->
    バージョン 0.0.1
  </p>
  <p>
    Twitter:
    <!-- svelte-ignore a11y-invalid-attribute -->
    <a
      href="#"
      on:click={() => BrowserOpenURL("https://twitter.com/tonango_0516")}
    >
      @tonango_0516
    </a>
    <br />
    GitHub:
    <!-- svelte-ignore a11y-invalid-attribute -->
    <a
      href="#"
      on:click={() =>
        BrowserOpenURL("https://github.com/tona0516/wows-fast-stats")}
    >
      wows-fast-stats
    </a>
  </p>
  <p>Copyright© 2023 tona0516 All Rights Reserved.</p>
</div>
