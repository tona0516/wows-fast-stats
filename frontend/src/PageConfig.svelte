<script lang="ts">
import {
  ApplyUserConfig,
  SelectDirectory,
  OpenDirectory,
} from "../wailsjs/go/main/App.js";
import { createEventDispatcher } from "svelte";
import Const from "./Const.js";
import { BrowserOpenURL } from "../wailsjs/runtime/runtime.js";
import clone from "clone";
import { storedUserConfig } from "./stores.js";
import { get } from "svelte/store";

const dispatch = createEventDispatcher();

let userConfig = get(storedUserConfig);
let inputUserConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => {
  userConfig = it;
  inputUserConfig = clone(it);
});

let isLoading = false;

async function clickApply() {
  isLoading = true;
  try {
    await ApplyUserConfig(inputUserConfig);
    storedUserConfig.set(inputUserConfig);
    dispatch("UpdateSuccess", { message: "設定を更新しました。" });
  } catch (error) {
    dispatch("Failure", { message: error });
  } finally {
    isLoading = false;
  }
}

async function openDirectory(path: string) {
  try {
    await OpenDirectory(path);
  } catch (error) {
    dispatch("Failure", { message: error });
  }
}

async function selectDirectory() {
  try {
    const selected = await SelectDirectory();
    if (!selected) return;
    inputUserConfig.install_path = selected;
  } catch (error) {
    dispatch("Failure", { message: error });
  }
}

function toggleAll(e: any) {
  const isSelectAll: boolean = e.target.checked;

  Object.keys(inputUserConfig.displays.ship).forEach(
    (key) => (inputUserConfig.displays.ship[key] = isSelectAll)
  );
  Object.keys(inputUserConfig.displays.overall).forEach(
    (key) => (inputUserConfig.displays.overall[key] = isSelectAll)
  );
}
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
          style="font-size: {userConfig.font_size};"
          bind:value="{inputUserConfig.install_path}"
        />
        <button
          type="button"
          class="btn btn-secondary"
          style="font-size: {userConfig.font_size};"
          on:click="{selectDirectory}">選択</button
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
        style="font-size: {userConfig.font_size};"
        bind:value="{inputUserConfig.appid}"
      />
      <p>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
          class="td-link"
          href="#"
          on:click="{() => BrowserOpenURL('https://developers.wargaming.net/')}"
          >Developer Room</a
        > で作成したIDを入力してください。
      </p>
    </div>

    <!-- font-size -->
    <div class="mb-3">
      <div class="centerize">
        <label for="font-size" class="form-label">文字サイズ</label>
      </div>
      <select
        class="form-select"
        style="font-size: {userConfig.font_size};"
        bind:value="{inputUserConfig.font_size}"
      >
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
              on:change="{toggleAll}"
              checked="{Object.values(inputUserConfig.displays.ship).filter(
                (it) => !it
              ).length === 0 &&
                Object.values(inputUserConfig.displays.overall).filter(
                  (it) => !it
                ).length === 0}"
            />
            <label class="form-check-label" for="select-all">全選択</label>
          </div>
        </div>
        <div class="col">
          {#each Object.keys(inputUserConfig.displays.ship) as key}
            {@const prefix = "ship"}
            <div class="form-check">
              <input
                class="form-check-input"
                type="checkbox"
                id="{prefix}-{key}"
                bind:checked="{inputUserConfig.displays.ship[key]}"
              />
              <label class="form-check-label" for="{prefix}-{key}"
                >艦:{Const.COLUMN_NAMES[key].full}</label
              >
            </div>
          {/each}
        </div>
        <div class="col">
          {#each Object.keys(inputUserConfig.displays.overall) as key}
            {@const prefix = "overall"}
            <div class="form-check">
              <input
                class="form-check-input"
                type="checkbox"
                id="{prefix}-{key}"
                bind:checked="{inputUserConfig.displays.overall[key]}"
              />
              <label class="form-check-label" for="{prefix}-{key}"
                >総合:{Const.COLUMN_NAMES[key].full}</label
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
            bind:checked="{inputUserConfig.save_screenshot}"
          />
          <label class="form-check-label" for="save-scrrenshot"
            >自動でスクリーンショットを保存する</label
          >
          <br />
          <!-- svelte-ignore a11y-invalid-attribute -->
          <a
            class="td-link"
            href="#"
            on:click="{() => openDirectory('screenshot/')}"
            ><i class="bi bi-folder2-open"></i> 保存フォルダを開く
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
            bind:checked="{inputUserConfig.save_temp_arena_info}"
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
            on:click="{() => openDirectory('temp_arena_info/')}"
            ><i class="bi bi-folder2-open"></i> 保存フォルダを開く</a
          >
        </div>
      </div>
    </div>

    <div class="centerize">
      <!-- apply -->
      <button
        type="button"
        class="btn btn-primary mb-3"
        style="font-size: {userConfig.font_size};"
        disabled="{isLoading}"
        on:click="{clickApply}"
      >
        {#if isLoading}
          <span
            class="spinner-border spinner-border-sm"
            role="status"
            aria-hidden="true"></span>
          更新中...
        {:else}
          適用
        {/if}
      </button>
    </div>
  </form>
</div>
