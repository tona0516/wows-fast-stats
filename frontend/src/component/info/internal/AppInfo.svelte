<script lang="ts">
import iconApp from "src/assets/images/appicon.png";
import ExternalLink from "src/component/common/ExternalLink.svelte";
import { Semver } from "wailsjs/go/main/App";

const links = [
  {
    icon: "question",
    url: "https://github.com/tona0516/wows-fast-stats/wiki/FAQ",
    text: "FAQ",
  },
  {
    icon: "twitter",
    url: "https://twitter.com/tonango_0516",
    text: "@tonango_0516",
  },
  {
    icon: "github",
    url: "https://github.com/tona0516/wows-fast-stats",
    text: "tona0516/wows-fast-stats",
  },
];
</script>

<div>
  <div>
    <div class="flex">
      <img src={iconApp} alt="" width="128px" height="128px" />
    </div>

    <div class="flex">
      <div class="flex">
        wows-fast-stats
        {#await Semver() then semver} {semver} {/await}
      </div>
    </div>
  </div>

  <div>
    {#each links as link}
      <div class="flex">
        <ExternalLink url={link.url}>
          <i class="bi bi-{link.icon}">{link.text}</i>
        </ExternalLink>
      </div>
    {/each}
  </div>

  <div class="flex">Copyright © 2023 tona0516 All Rights Reserved.</div>
</div>
