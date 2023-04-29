<script lang="ts">
  import type { vo } from "wailsjs/go/models.js";
  import {
    ApplyUserConfig,
    UserConfig,
    SelectDirectory,
    Cwd,
    OpenDirectory,
  } from "../wailsjs/go/main/App.js";
  import { createEventDispatcher } from "svelte";
  import Const from "./Const.js";
  import { BrowserOpenURL } from "../wailsjs/runtime/runtime.js";

  const dispatch = createEventDispatcher();

  let inputConfig: vo.UserConfig = Const.DEFAULT_USER_CONFIG;

  let cwd: string;

  UserConfig().then((config) => {
    inputConfig = config;
  });

  function clickApply() {
    ApplyUserConfig(inputConfig)
      .then((_) => {
        dispatch("SuccessToast", {
          message: "更新しました。",
        });
      })
      .catch((error) => {
        dispatch("ErrorToast", {
          message: error,
        });
      });
  }

  function openDirectory(path: string) {
    OpenDirectory(path).catch((error) => {
      dispatch("ErrorToast", {
        message: error,
      });
    });
  }

  function selectDirectory() {
    SelectDirectory().then((result) => {
      if (!result) return;
      inputConfig.install_path = result;
    });
  }

  function toggleAll(e) {
    const isSelectAll: boolean = e.target.checked;

    Object.keys(inputConfig.displays.ship).forEach(
      (key) => (inputConfig.displays.ship[key] = isSelectAll)
    );
    Object.keys(inputConfig.displays.overall).forEach(
      (key) => (inputConfig.displays.overall[key] = isSelectAll)
    );
  }

  Cwd()
    .then((result) => (cwd = result))
    .catch((error) => "");
</script>

<div class="mt-3 form-style">
  <form>
    <!-- install path -->
    <div class="mb-3">
      <div class="centerize">
        <label for="install-path" class="form-label"
          >World of Warshipsインストールフォルダ</label
        >
      </div>
      <div class="horizontal">
        <input
          type="text"
          class="form-control"
          id="install-path"
          bind:value={inputConfig.install_path}
        />
        <button
          type="button"
          class="btn btn-secondary"
          on:click={selectDirectory}>選択</button
        >
      </div>
    </div>

    <!-- appid -->
    <div class="mb-3">
      <div class="centerize">
        <label for="appid" class="form-label">AppID</label>
      </div>
      <input
        type="text"
        class="form-control"
        id="appid"
        bind:value={inputConfig.appid}
      />
      <p>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
          class="td-link"
          href="#"
          on:click={() =>
            BrowserOpenURL("https://developers.wargaming.net/applications/")}
          >Developer Room <i class="bi bi-box-arrow-up-right" /></a
        > で作成したIDを入力してください。
      </p>
    </div>

    <!-- font-size -->
    <div class="mb-3">
      <div class="centerize">
        <label for="font-size" class="form-label">文字サイズ</label>
      </div>
      <select class="form-select" bind:value={inputConfig.font_size}>
        <option value="x-small">極小</option>
        <option value="small">小</option>
        <option value="medium">中</option>
        <option value="large">大</option>
        <option value="x-large">極大</option>
      </select>
    </div>

    <!-- display values -->
    <div class="mb-3">
      <div class="centerize">
        <label for="font-column" class="form-lavel">表示項目</label>
      </div>
      <div class="row">
        <div class="col">
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="select-all"
              on:change={toggleAll}
              checked={Object.values(inputConfig.displays.ship).filter(
                (it) => !it
              ).length === 0 &&
                Object.values(inputConfig.displays.overall).filter((it) => !it)
                  .length === 0}
            />
            <label class="form-check-label" for="select-all">全選択</label>
          </div>
        </div>
        <div class="col">
          {#each Object.keys(inputConfig.displays.ship) as key}
            {@const prefix = "ship"}
            <div class="form-check">
              <input
                class="form-check-input"
                type="checkbox"
                id="{prefix}-{key}"
                bind:checked={inputConfig.displays.ship[key]}
              />
              <label class="form-check-label" for="{prefix}-{key}"
                >艦:{Const.COLUMN_NAMES[key].fullName}</label
              >
            </div>
          {/each}
        </div>
        <div class="col">
          {#each Object.keys(inputConfig.displays.overall) as key}
            {@const prefix = "overall"}
            <div class="form-check">
              <input
                class="form-check-input"
                type="checkbox"
                id="{prefix}-{key}"
                bind:checked={inputConfig.displays.overall[key]}
              />
              <label class="form-check-label" for="{prefix}-{key}"
                >総合:{Const.COLUMN_NAMES[key].fullName}</label
              >
            </div>
          {/each}
        </div>
      </div>

      <!-- save-screenshot -->
      <div class="mb-3">
        <div class="form-check form-switch">
          <input
            class="form-check-input"
            type="checkbox"
            id="save-scrrenshot"
            bind:checked={inputConfig.save_screenshot}
          />
          <label class="form-check-label" for="save-scrrenshot"
            >自動でスクリーンショットを保存する</label
          >
          <br />
          <!-- svelte-ignore a11y-invalid-attribute -->
          <a
            class="td-link"
            href="#"
            on:click={() => openDirectory(cwd + "/screenshot")}
            ><i class="bi bi-folder2-open" /> 保存フォルダを開く
          </a>
        </div>
      </div>

      <!-- save-temp-arena-info -->
      <div class="mb-3">
        <div class="form-check form-switch">
          <input
            class="form-check-input"
            type="checkbox"
            id="save-temp-arena-info"
            bind:checked={inputConfig.save_temp_arena_info}
          />
          <label class="form-check-label" for="save-temp-arena-info"
            >【開発用】自動で戦闘情報(<i>tempArenaInfo.json</i
            >)を保存する</label
          >
          <br />
          <!-- svelte-ignore a11y-invalid-attribute -->
          <a
            class="td-link"
            href="#"
            on:click={() => openDirectory(cwd + "/temp_arena_info")}
            ><i class="bi bi-folder2-open" /> 保存フォルダを開く</a
          >
        </div>
      </div>
    </div>

    <div class="centerize">
      <!-- apply -->
      <button type="button" class="btn btn-primary mb-3" on:click={clickApply}
        >適用</button
      >
    </div>
  </form>
</div>
