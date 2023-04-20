<script lang="ts">
  import type { vo } from "wailsjs/go/models.js";
  import {
    ApplyConfig,
    GetConfig,
    SelectDirectory,
  } from "../wailsjs/go/main/App.js";
  import { createEventDispatcher } from "svelte";
  import Const from "./Const.js";

  const dispatch = createEventDispatcher();

  let inputConfig: vo.UserConfig = Const.DEFAULT_USER_CONFIG;

  GetConfig().then((config) => {
    inputConfig = config;
  });

  function clickApply() {
    ApplyConfig(inputConfig)
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
      inputConfig.install_path = result;
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
    <div class="mb-3 form-style">
      <label for="appid" class="form-label">AppID</label>
      <input
        type="text"
        class="form-control"
        id="appid"
        bind:value={inputConfig.appid}
      />
    </div>

    <!-- font-size -->
    <div class="mb-3 form-style">
      <label for="font-size" class="form-label">文字サイズ</label>
      <select class="form-select" bind:value={inputConfig.font_size}>
        <option value="x-small">極小</option>
        <option value="small">小</option>
        <option value="medium">中</option>
        <option value="large">大</option>
        <option value="x-large">極大</option>
      </select>
    </div>

    <!-- display values -->
    <div class="mb-3 form-style">
      <label for="font-column" class="form-lavel">表示項目</label>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="pr"
          bind:checked={inputConfig.displays.pr}
        />
        <label class="form-check-label" for="pr">PR</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="shio-damage"
          bind:checked={inputConfig.displays.ship_damage}
        />
        <label class="form-check-label" for="shio-damage ">艦別:ダメージ</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="ship-win-rate"
          bind:checked={inputConfig.displays.ship_win_rate}
        />
        <label class="form-check-label" for="ship-win-rate">艦別:勝率</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="ship-kd-rate"
          bind:checked={inputConfig.displays.ship_kd_rate}
        />
        <label class="form-check-label" for="ship-kd-rate">艦別:K/D</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="ship-win-survived-rate"
          bind:checked={inputConfig.displays.ship_win_survived_rate}
        />
        <label class="form-check-label" for="ship-win-survived-rate"
          >艦別:勝利生存率</label
        >
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="ship-lose-survived-rate"
          bind:checked={inputConfig.displays.ship_lose_survived_rate}
        />
        <label class="form-check-label" for="ship-lose-survived-rate"
          >艦別:敗北生存率</label
        >
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="ship-exp"
          bind:checked={inputConfig.displays.ship_exp}
        />
        <label class="form-check-label" for="ship-exp">艦別:経験値</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="ship-battles"
          bind:checked={inputConfig.displays.ship_battles}
        />
        <label class="form-check-label" for="ship-battles">艦別:戦闘数</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="player-damage"
          bind:checked={inputConfig.displays.player_damage}
        />
        <label class="form-check-label" for="player-damage">総合:ダメージ</label
        >
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="player-win-rate"
          bind:checked={inputConfig.displays.player_win_rate}
        />
        <label class="form-check-label" for="player-win-rate">総合:勝率</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="player-kd-rate"
          bind:checked={inputConfig.displays.player_kd_rate}
        />
        <label class="form-check-label" for="player-kd-rate">総合:K/D</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="player-win-survived-rate"
          bind:checked={inputConfig.displays.player_win_survived_rate}
        />
        <label class="form-check-label" for="player-win-survived-rate"
          >総合:勝利生存率</label
        >
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="player-lose-survived-rate"
          bind:checked={inputConfig.displays.player_lose_survived_rate}
        />
        <label class="form-check-label" for="player-lose-survived-rate"
          >総合:敗北生存率</label
        >
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="player-exp"
          bind:checked={inputConfig.displays.player_exp}
        />
        <label class="form-check-label" for="player-exp">総合:経験値</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="player-battles"
          bind:checked={inputConfig.displays.player_battles}
        />
        <label class="form-check-label" for="player-battles">総合:戦闘数</label>
      </div>
      <div class="form-check">
        <input
          class="form-check-input"
          type="checkbox"
          id="avg-tier"
          bind:checked={inputConfig.displays.player_avg_tier}
        />
        <label class="form-check-label" for="avg-tier">総合:平均Tier</label>
      </div>
    </div>

    <!-- apply -->
    <button type="button" class="btn btn-primary" on:click={clickApply}
      >適用</button
    >
  </form>
</div>
