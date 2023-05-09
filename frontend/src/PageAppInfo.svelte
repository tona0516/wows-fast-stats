<script lang="ts">
import { AppVersion } from "../wailsjs/go/main/App";
import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
import iconApp from "./assets/images/appicon.png";
import type { vo } from "wailsjs/go/models";

let version: vo.Version = { semver: "", revision: "" };
AppVersion().then((result) => (version = result));
</script>

<div class="mt-3 appinfo">
  <img class="m-1" src="{iconApp}" alt="" width="128px" height="128px" />

  <p>
    wows-fast-stats
    <br />
    バージョン {version.semver}
    <br />
    リビジョン {version.revision}
  </p>
  <p>
    <i class="bi bi-question-circle"></i>
    <!-- svelte-ignore a11y-invalid-attribute -->
    <a
      class="td-link"
      href="#"
      on:click="{() =>
        BrowserOpenURL('https://github.com/tona0516/wows-fast-stats/wiki/FAQ')}"
    >
      FAQ
    </a>
    <br />
    <i class="bi bi-twitter"></i>
    <!-- svelte-ignore a11y-invalid-attribute -->
    <a
      class="td-link"
      href="#"
      on:click="{() => BrowserOpenURL('https://twitter.com/tonango_0516')}"
    >
      @tonango_0516
    </a>
    <br />
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
  </p>
  <p>Copyright © 2023 tona0516 All Rights Reserved.</p>
</div>
