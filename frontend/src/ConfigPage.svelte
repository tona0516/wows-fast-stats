<script lang="ts">
  import {
    ApplyConfig,
    GetConfig,
    SelectDirectory,
  } from "../wailsjs/go/main/App.js";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  let installPath = "";
  let appid = "";
  let fontSize = "";

  GetConfig()
    .then((config) => {
      installPath = config.install_path;
      appid = config.appid;
      fontSize = config.font_size;
    })
    .catch((_) => {
      installPath = "";
      appid = "";
      fontSize = "medium";
    });

  function clickApply() {
    ApplyConfig(installPath, appid, fontSize)
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

  function selectDirectory() {
    SelectDirectory().then((result) => {
      if (!result) return;
      installPath = result;
    });
  }
</script>

<div class="mt-3">
  <form>
    <!-- install path -->
    <div class="mb-3 form-style">
      <label for="install-path" class="form-label"
        >World of Warshipsインストールフォルダ</label
      >
      <div class="horizontal">
        <input
          type="text"
          class="form-control"
          id="install-path"
          bind:value={installPath}
        />
        <button
          type="button"
          class="btn btn-secondary"
          on:click={selectDirectory}>選択</button
        >
      </div>
    </div>

    <!-- appid -->
    <div class="mb-3 form-style">
      <label for="appid" class="form-label">AppID</label>
      <input type="text" class="form-control" id="appid" bind:value={appid} />
    </div>

    <!-- font-size -->
    <div class="mb-3 form-style">
      <label for="font-size" class="form-label">文字サイズ</label>
      <select class="form-select" bind:value={fontSize}>
        <option value="x-small">極小</option>
        <option value="small">小</option>
        <option value="medium">中</option>
        <option value="large">大</option>
        <option value="x-large">極大</option>
      </select>
    </div>

    <!-- apply -->
    <button type="button" class="btn btn-primary" on:click={clickApply}
      >適用</button
    >
  </form>
</div>
