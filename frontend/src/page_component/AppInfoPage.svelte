<script lang="ts">
import { AppVersion } from "../../wailsjs/go/main/App";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import iconApp from "../assets/images/appicon.png";
</script>

<div class="mt-3 center">
  <img class="m-1" src="{iconApp}" alt="" width="128px" height="128px" />

  <div class="m-1 center">
    <div>wows-fast-stats</div>

    {#await AppVersion() then version}
      <div>バージョン {version.semver}</div>
      <div>リビジョン {version.revision}</div>
    {/await}
  </div>

  <div class="m-1 center">
    <div>
      <i class="bi bi-question-circle"></i>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a
        class="td-link"
        href="#"
        on:click="{() =>
          BrowserOpenURL(
            'https://github.com/tona0516/wows-fast-stats/wiki/FAQ'
          )}"
      >
        FAQ
      </a>
    </div>

    <div>
      <i class="bi bi-twitter"></i>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a
        class="td-link"
        href="#"
        on:click="{() => BrowserOpenURL('https://twitter.com/tonango_0516')}"
      >
        @tonango_0516
      </a>
    </div>

    <div>
      <i class="bi bi-github"></i>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a
        class="td-link"
        href="#"
        on:click="{() =>
          BrowserOpenURL('https://github.com/tona0516/wows-fast-stats')}"
      >
        wows-fast-stats
      </a>
    </div>
  </div>

  <div class="m-1">Copyright © 2023 tona0516 All Rights Reserved.</div>
</div>
