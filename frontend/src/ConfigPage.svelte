<script lang="ts">
  import type { vo } from "wailsjs/go/models.js";
  import {
    ApplyUserConfig,
    UserConfig,
    SelectDirectory,
  } from "../wailsjs/go/main/App.js";
  import { createEventDispatcher } from "svelte";
  import Const from "./Const.js";

  const dispatch = createEventDispatcher();

  let inputConfig: vo.UserConfig = Const.DEFAULT_USER_CONFIG;

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

  function selectDirectory() {
    SelectDirectory().then((result) => {
      if (!result) return;
      inputConfig.install_path = result;
    });
  }

  function toggleAll(e) {
    const keys = Object.keys(inputConfig.displays);
    const isSelectAll: boolean = e.target.checked;
    keys.forEach((key) => (inputConfig.displays[key] = isSelectAll));
  }
</script>

<div class="mt-3">
  <form>
    <!-- install path -->
    <div class="mb-3 form-style">
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
    <div class="mb-3 form-style">
      <div class="centerize">
        <label for="appid" class="form-label">AppID</label>
      </div>
      <input
        type="text"
        class="form-control"
        id="appid"
        bind:value={inputConfig.appid}
      />
    </div>

    <!-- font-size -->
    <div class="mb-3 form-style">
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
    <div class="mb-3 form-style">
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
              checked={Object.values(inputConfig.displays).length ===
                Object.values(inputConfig.displays).filter((it) => it === true)
                  .length}
            />
            <label class="form-check-label" for="select-all">全選択</label>
          </div>
        </div>
        <div class="col">
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="player-name"
              bind:checked={inputConfig.displays.player_name}
            />
            <label class="form-check-label" for="player-name"
              >プレイヤー名</label
            >
          </div>
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
            <label class="form-check-label" for="shio-damage "
              >艦別:ダメージ</label
            >
          </div>
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="ship-win-rate"
              bind:checked={inputConfig.displays.ship_win_rate}
            />
            <label class="form-check-label" for="ship-win-rate">艦別:勝率</label
            >
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
            <label class="form-check-label" for="ship-battles"
              >艦別:戦闘数</label
            >
          </div>
        </div>
        <div class="col">
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="player-damage"
              bind:checked={inputConfig.displays.player_damage}
            />
            <label class="form-check-label" for="player-damage"
              >総合:ダメージ</label
            >
          </div>
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="player-win-rate"
              bind:checked={inputConfig.displays.player_win_rate}
            />
            <label class="form-check-label" for="player-win-rate"
              >総合:勝率</label
            >
          </div>
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="player-kd-rate"
              bind:checked={inputConfig.displays.player_kd_rate}
            />
            <label class="form-check-label" for="player-kd-rate">総合:K/D</label
            >
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
            <label class="form-check-label" for="player-battles"
              >総合:戦闘数</label
            >
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
          <div class="form-check">
            <input
              class="form-check-input"
              type="checkbox"
              id="player_using_ship_type_rate"
              bind:checked={inputConfig.displays.player_using_ship_type_rate}
            />
            <label class="form-check-label" for="player_using_ship_type_rate"
              >総合:使用艦率(SS|DD|CL|BB|CV)</label
            >
          </div>
        </div>
      </div>

      <!-- save-screenshot -->
      <div class="mb-3 centerize">
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
        </div>
      </div>
    </div>

    <!-- apply -->
    <button type="button" class="btn btn-primary" on:click={clickApply}
      >適用</button
    >
  </form>
</div>
