<script lang="ts">
  import MainPage from "src/component/main/MainPage.svelte";
  import ConfigPage from "src/component/config/ConfigPage.svelte";

  import "bootstrap-icons/font/bootstrap-icons.css";
  import "charts.css";

  import { LatestRelease, StartWatching } from "wailsjs/go/main/App";
  import { EventsOn, EventsEmit } from "wailsjs/runtime/runtime";
  import { storedConfig, storedRequiredConfigError } from "src/stores";
  import Modals from "src/component/modal/AlertModals.svelte";
  import { domain } from "wailsjs/go/models";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
  import UkTab from "src/component/common/uikit/UkTab.svelte";
  import { FontSize } from "src/lib/FontSize";
  import { FetchProxy } from "./lib/FetchProxy";
  import { Notifier } from "./lib/Notifier";

  let modals: Modals;
  let mainPage: MainPage | undefined;
  let initialized = false;
  let updatableRelease: domain.GHLatestRelease;

  $: {
    // @ts-ignore
    document.body.style.zoom = FontSize.getZoomRate($storedConfig);
  }

  EventsOn("BATTLE_START", () => mainPage?.fetchBattle());
  EventsOn("BATTLE_ERR", (error: string) => Notifier.failure(error));

  async function main() {
    try {
      const config = await FetchProxy.getConfig();
      const requiredConfigError = await FetchProxy.validateRequiredConfig(
        config.install_path,
        config.appid,
      );
      await FetchProxy.getAlertPlayers();

      initialized = true;

      if (config.notify_updatable) {
        const latestRelease = await LatestRelease();
        if (latestRelease.updatable) {
          updatableRelease = latestRelease;
        }
      }

      if (requiredConfigError.valid) {
        StartWatching();
      }
    } catch (error) {
      Notifier.failure(error);
    }
  }

  window.onload = function () {
    EventsEmit("ONLOAD");
  };

  main();
</script>

<main>
  <Modals bind:this={modals} />

  {#if updatableRelease}
    <div class="uk-flex uk-flex-center uk-background-secondary">
      新しいバージョンがあります:
      <ExternalLink url={updatableRelease.html_url}>
        {updatableRelease.tag_name}
      </ExternalLink>
    </div>
  {/if}

  {#if initialized}
    {@const tabID = "page-tab"}
    <UkTab clazz="uk-margin-remove" id={tabID}>
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#"><UkIcon name="home" /></a>
      </li>
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#">
          <UkIcon name="cog" />
          {#if !$storedRequiredConfigError.valid}
            <span class="uk-text-warning uk-text-small">
              <UkIcon name="warning" />
            </span>
          {/if}
        </a>
      </li>
    </UkTab>
    <ul id={tabID} class="uk-switcher">
      <li>
        <MainPage
          bind:this={mainPage}
          on:EditAlertPlayer={(e) =>
            modals.showEditAlertPlayer(e.detail.target)}
          on:RemoveAlertPlayer={(e) =>
            modals.showRemoveAlertPlayer(e.detail.target)}
        />
      </li>
      <li>
        <ConfigPage
          on:AddAlertPlayer={() => modals.showAddAlertPlayer()}
          on:EditAlertPlayer={(e) =>
            modals.showEditAlertPlayer(e.detail.target)}
          on:RemoveAlertPlayer={(e) =>
            modals.showRemoveAlertPlayer(e.detail.target)}
        />
      </li>
    </ul>
  {:else}
    <div class="uk-overlay-default">
      <div class="uk-position-center">
        <UkSpinner />
      </div>
    </div>
  {/if}
</main>
