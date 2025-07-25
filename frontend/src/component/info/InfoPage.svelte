<script lang="ts">
  import iconApp from "src/assets/images/appicon.png";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import { Semver } from "wailsjs/go/main/App";

  const LINKS = [
    {
      icon: "question-circle",
      url: "https://github.com/tona0516/wows-fast-stats/wiki/FAQ",
      text: "FAQ",
    },
    {
      icon: "twitter-x",
      url: "https://x.com/tonango_wows",
      text: "@tonango_0516",
    },
    {
      icon: "github",
      url: "https://github.com/tona0516/wows-fast-stats",
      text: "tona0516/wows-fast-stats",
    },
  ];
</script>

<div class="p-4 flex flex-col items-center">
  <img src={iconApp} alt="" width="128px" height="128px" />
  <div class="pt-1">
    wows-fast-stats {#await Semver() then semver} {semver} {/await}
  </div>
  <div class="pt-2">
    {#each LINKS as link}
      <div class="flex flex-col items-center">
        <ExternalLink url={link.url}>
          <i class="bi bi-{link.icon}">{link.text}</i>
        </ExternalLink>
      </div>
    {/each}
  </div>
  <div class="pt-2">Copyright © 2025 tona0516 All Rights Reserved.</div>
</div>
