<script lang="ts">
  import { Semver } from "wailsjs/go/main/App";
  import iconApp from "src/assets/images/appicon.png";
  import Logging from "src/other_component/Logging.svelte";
  import { storedLogs } from "src/stores";
  import ExternalLink from "src/other_component/ExternalLink.svelte";
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
      <ExternalLink
        url="https://github.com/tona0516/wows-fast-stats/wiki/FAQ"
        text="FAQ"
      />
    </div>

    <div>
      <i class="bi bi-twitter" />
      <ExternalLink
        url="https://twitter.com/tonango_0516"
        text="@tonango_0516"
      />
    </div>

    <div>
      <i class="bi bi-github" />
      <ExternalLink
        url="https://github.com/tona0516/wows-fast-stats"
        text="tona0516/wows-fast-stats"
      />
    </div>
  </div>

  <Logging logs={$storedLogs} />

  <div class="m-1">Copyright © 2023 tona0516 All Rights Reserved.</div>
</div>
