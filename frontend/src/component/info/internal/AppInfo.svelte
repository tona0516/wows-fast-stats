<script lang="ts">
  import { Semver } from "wailsjs/go/main/App";
  import iconApp from "src/assets/images/appicon.png";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";

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

<div class="uk-padding-small">
  <div class="uk-margin-small">
    <div class="uk-flex uk-flex-center">
      <img src={iconApp} alt="" width="128px" height="128px" />
    </div>

    <div class="uk-flex uk-flex-center">
      <div class="uk-flex uk-text-emphasis">
        wows-fast-stats
        {#await Semver() then semver} {semver} {/await}
      </div>
    </div>
  </div>

  <div class="uk-margin-small">
    {#each links as link}
      <div class="uk-flex uk-flex-center">
        <ExternalLink url={link.url}>
          <UkIcon name={link.icon} />
          <span class="uk-text-middle">{link.text}</span>
        </ExternalLink>
      </div>
    {/each}
  </div>

  <div class="uk-margin-small uk-flex uk-flex-center">
    Copyright © 2023 tona0516 All Rights Reserved.
  </div>
</div>
