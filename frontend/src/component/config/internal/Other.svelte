<script lang="ts">
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import { createEventDispatcher } from "svelte";
  import { OpenDirectory } from "wailsjs/go/main/App";
  import { domain } from "wailsjs/go/models";

  export let inputConfig: domain.UserConfig;

  const dispatch = createEventDispatcher();

  const openDirectory = async (path: string) => {
    try {
      await OpenDirectory(path);
    } catch (error) {
      dispatch("Failure", { message: error });
    }
  };
</script>

<div class="uk-padding-small">
  <div class="uk-margin-small-bottom">
    <label
      ><input
        class="uk-checkbox"
        type="checkbox"
        bind:checked={inputConfig.save_screenshot}
        on:change={() => dispatch("Change")}
      /> 自動でスクリーンショットを保存する</label
    >
    <div>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a class="td-link" href="#" on:click={() => openDirectory("screenshot/")}>
        <UkIcon name="folder" />
        <span class="uk-text-middle">保存フォルダを開く</span>
      </a>
    </div>
  </div>

  <div class="uk-margin-small-bottom">
    <label
      ><input
        class="uk-checkbox"
        type="checkbox"
        bind:checked={inputConfig.save_temp_arena_info}
        on:change={() => dispatch("Change")}
      /> 【開発用】自動で戦闘情報(tempArenaInfo.json)を保存する</label
    >
    <div>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a
        class="td-link"
        href="#"
        on:click={() => openDirectory("temp_arena_info/")}
      >
        <UkIcon name="folder" />
        <span class="uk-text-middle">保存フォルダを開く</span>
      </a>
    </div>
  </div>

  <div class="uk-margin-small-bottom">
    <label
      ><input
        class="uk-checkbox"
        type="checkbox"
        bind:checked={inputConfig.send_report}
        on:change={() => dispatch("Change")}
      /> アプリ改善のためのデータ送信を許可する</label
    >
    <div>
      <ul class="uk-list uk-list-disc uk-list-collapse">
        <li>アプリバージョン</li>
        <li>エラーログ</li>
        <li>設定値(config/user.json)</li>
        <li>戦闘情報(tempArenaInfo.json)</li>
      </ul>
    </div>
  </div>

  <div class="uk-margin-small-bottom">
    <label
      ><input
        class="uk-checkbox"
        type="checkbox"
        bind:checked={inputConfig.notify_updatable}
        on:change={() => dispatch("Change")}
      /> 新しいバージョンがある場合に通知する</label
    >
  </div>
</div>
