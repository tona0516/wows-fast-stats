<script lang="ts">
import clone from "clone";
import { createEventDispatcher } from "svelte";
import { get } from "svelte/store";
import {
  ApplyUserConfig,
  DefaultUserConfig,
  FontSizes,
  OpenDirectory,
  SampleTeams,
  SelectDirectory,
} from "../../wailsjs/go/main/App";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { storedUserConfig } from "../stores";
import { Const } from "../Const";
import ColorPicker from "svelte-awesome-color-picker";
import type { vo } from "../../wailsjs/go/models";
import StatisticsTable from "../other_component/StatisticsTable.svelte";

const dispatch = createEventDispatcher();

let userConfig = get(storedUserConfig);
let inputUserConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => {
  userConfig = it;
  inputUserConfig = clone(it);
});

let isLoading = false;
let defaultUserConfig: vo.UserConfig;
let sampleTeams: vo.Team[] = [];
let displayKeys: string[] = [];

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

async function main() {
  defaultUserConfig = await DefaultUserConfig();
  sampleTeams = await SampleTeams();
  const shipKeys = Object.keys(userConfig.displays.ship);
  const overallKeys = Object.keys(userConfig.displays.overall);
  displayKeys = Array.from(new Set([...shipKeys, ...overallKeys]));
}

main();
</script>

<div class="mt-3">
  <form>
    <!-- install path -->
    <div class="mb-2 form-style">
      <h6 class="text-center">
        World of Warshipsインストールフォルダ <span class="badge bg-danger"
          >必須</span
        >
      </h6>

      <div class="centerize">
        <input
          type="text"
          class="form-control"
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
    <div class="mb-2 form-style">
      <h6 class="text-center">
        アプリケーションID <span class="badge bg-danger">必須</span>
      </h6>

      <input
        type="text"
        class="form-control"
        style="font-size: {userConfig.font_size};"
        bind:value="{inputUserConfig.appid}"
      />
      <div class="centerize">
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
          class="td-link"
          href="#"
          on:click="{() => BrowserOpenURL('https://developers.wargaming.net/')}"
          >Developer Room</a
        > で作成したIDを入力してください。
      </div>
    </div>

    <!-- font-size -->
    <div class="mb-2 form-style">
      <h6 class="text-center">文字サイズ</h6>

      <select
        class="form-select"
        style="font-size: {userConfig.font_size};"
        bind:value="{inputUserConfig.font_size}"
      >
        {#await FontSizes() then fontSizes}
          {#each fontSizes as fs}
            <option selected="{fs === userConfig.font_size}" value="{fs}"
              >{Const.FONT_SIZE[fs]}</option
            >
          {/each}
        {/await}
      </select>
    </div>

    <!-- display values -->
    <div class="mb-2 form-style">
      <h6 class="text-center">表示項目</h6>

      <div class="centerize">
        <div class="form-check form-switch">
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

      <div class="centerize">
        <table class="table table-sm table-text-color w-auto">
          <thead>
            <tr>
              <th>項目</th>
              <th>艦成績</th>
              <th>総合成績</th>
              <th>小数点以下の桁数</th>
            </tr>
          </thead>
          <tbody>
            {#each displayKeys as key}
              <tr>
                <td>
                  <div class="centerize">
                    {Const.COLUMN_NAMES[key].full}
                  </div>
                </td>

                <td>
                  <div class="centerize">
                    {#if inputUserConfig.displays.ship[key] !== undefined}
                      <div class="form-check form-switch my-1">
                        <input
                          class="form-check-input"
                          type="checkbox"
                          bind:checked="{inputUserConfig.displays.ship[key]}"
                        />
                      </div>
                    {/if}
                  </div>
                </td>

                <td>
                  <div class="centerize">
                    {#if inputUserConfig.displays.overall[key] !== undefined}
                      <div class="form-check form-switch my-1">
                        <input
                          class="form-check-input"
                          type="checkbox"
                          bind:checked="{inputUserConfig.displays.overall[key]}"
                        />
                      </div>
                    {/if}
                  </div>
                </td>

                <td>
                  <div class="centerize">
                    {#if inputUserConfig.custom_digit[key] !== undefined}
                      <div class="my-1">
                        <select
                          class="form-select"
                          style="font-size: {userConfig.font_size};"
                          bind:value="{inputUserConfig.custom_digit[key]}"
                        >
                          {#each [0, 1, 2] as digit}
                            <option
                              selected="{digit ===
                                inputUserConfig.custom_digit[key]}"
                              value="{digit}">{digit}</option
                            >
                          {/each}
                        </select>
                      </div>
                    {/if}
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <!-- custom color -->
    <div class="mb-2 form-style">
      <!-- skill color -->
      <h6 class="text-center">スキル別カラー</h6>
      <div class="centerize">
        <table class="table table-sm table-text-color w-auto">
          <thead>
            <tr>
              <th>スキル</th>
              <th>文字色</th>
              <th>背景色</th>
            </tr>
          </thead>
          <tbody>
            {#each Object.keys(inputUserConfig.custom_color.skill.text) as key}
              <tr>
                <td>
                  <div class="centerize">
                    {Const.SKILL_LEVEL_LABELS[key]}
                  </div>
                </td>
                <td>
                  <div class="centerize my-1">
                    <ColorPicker
                      label="{inputUserConfig.custom_color.skill.text[key]}"
                      isAlpha="{false}"
                      canChangeMode="{false}"
                      bind:hex="{inputUserConfig.custom_color.skill.text[key]}"
                    />
                    <div class="mx-2">
                      <button
                        type="button"
                        class="btn btn-sm btn-success"
                        style="font-size: {userConfig.font_size};"
                        on:click="{() => {
                          inputUserConfig.custom_color.skill.text[key] =
                            defaultUserConfig.custom_color.skill.text[key];
                        }}">デフォルト値をセット</button
                      >
                    </div>
                  </div>
                </td>

                <td>
                  <div class="centerize my-1">
                    <ColorPicker
                      label="{inputUserConfig.custom_color.skill.background[
                        key
                      ]}"
                      isAlpha="{false}"
                      canChangeMode="{false}"
                      bind:hex="{inputUserConfig.custom_color.skill.background[
                        key
                      ]}"
                    />
                    <div class="mx-2">
                      <button
                        type="button"
                        class="btn btn-sm btn-success"
                        style="font-size: {userConfig.font_size};"
                        on:click="{() => {
                          inputUserConfig.custom_color.skill.background[key] =
                            defaultUserConfig.custom_color.skill.background[
                              key
                            ];
                        }}">デフォルト値をセット</button
                      >
                    </div>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <!-- tier group color -->
      <h6 class="text-center">Tier別カラー</h6>
      <div class="centerize">
        <table class="table table-sm table-text-color w-auto">
          <thead>
            <tr>
              <th>Tier</th>
              <th>使用艦</th>
              <th>非使用艦</th>
            </tr>
          </thead>
          <tbody>
            {#each Object.keys(inputUserConfig.custom_color.tier.own) as key}
              <tr>
                <td>
                  <div class="centerize">
                    {Const.TIER_GROUP_LABELS[key]}
                  </div>
                </td>
                <td>
                  <div class="centerize my-1">
                    <ColorPicker
                      label="{inputUserConfig.custom_color.tier.own[key]}"
                      isAlpha="{false}"
                      canChangeMode="{false}"
                      bind:hex="{inputUserConfig.custom_color.tier.own[key]}"
                    />
                    <div class="mx-2">
                      <button
                        type="button"
                        class="btn btn-sm btn-success"
                        style="font-size: {userConfig.font_size};"
                        on:click="{() => {
                          inputUserConfig.custom_color.tier.own[key] =
                            defaultUserConfig.custom_color.tier.own[key];
                        }}">デフォルト値をセット</button
                      >
                    </div>
                  </div>
                </td>

                <td>
                  <div class="centerize my-1">
                    <ColorPicker
                      label="{inputUserConfig.custom_color.tier.other[key]}"
                      isAlpha="{false}"
                      canChangeMode="{false}"
                      bind:hex="{inputUserConfig.custom_color.tier.other[key]}"
                    />
                    <div class="mx-2">
                      <button
                        type="button"
                        class="btn btn-sm btn-success"
                        style="font-size: {userConfig.font_size};"
                        on:click="{() => {
                          inputUserConfig.custom_color.tier.other[key] =
                            defaultUserConfig.custom_color.tier.other[key];
                        }}">デフォルト値をセット</button
                      >
                    </div>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <!-- ship type color -->
      <h6 class="text-center">艦別カラー</h6>
      <div class="centerize">
        <table class="table table-sm table-text-color w-auto">
          <thead>
            <tr>
              <th>艦種</th>
              <th>使用艦</th>
              <th>非使用艦</th>
            </tr>
          </thead>
          <tbody>
            {#each Object.keys(inputUserConfig.custom_color.ship_type.own) as key}
              <tr>
                <td>
                  <div class="centerize">
                    {Const.SHIP_TYPE_LABELS[key]}
                  </div>
                </td>
                <td>
                  <div class="centerize my-1">
                    <ColorPicker
                      label="{inputUserConfig.custom_color.ship_type.own[key]}"
                      isAlpha="{false}"
                      canChangeMode="{false}"
                      bind:hex="{inputUserConfig.custom_color.ship_type.own[
                        key
                      ]}"
                    />
                    <div class="mx-2">
                      <button
                        type="button"
                        class="btn btn-sm btn-success"
                        style="font-size: {userConfig.font_size};"
                        on:click="{() => {
                          inputUserConfig.custom_color.ship_type.own[key] =
                            defaultUserConfig.custom_color.ship_type.own[key];
                        }}">デフォルト値をセット</button
                      >
                    </div>
                  </div>
                </td>

                <td>
                  <div class="centerize my-1">
                    <ColorPicker
                      label="{inputUserConfig.custom_color.ship_type.other[
                        key
                      ]}"
                      isAlpha="{false}"
                      canChangeMode="{false}"
                      bind:hex="{inputUserConfig.custom_color.ship_type.other[
                        key
                      ]}"
                    />
                    <div class="mx-2">
                      <button
                        type="button"
                        class="btn btn-sm btn-success"
                        style="font-size: {userConfig.font_size};"
                        on:click="{() => {
                          inputUserConfig.custom_color.ship_type.other[key] =
                            defaultUserConfig.custom_color.ship_type.other[key];
                        }}">デフォルト値をセット</button
                      >
                    </div>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <div class="m-2">
      <h6 class="text-center">プレビュー</h6>

      <div style="font-size: {inputUserConfig.font_size};">
        <StatisticsTable
          teams="{sampleTeams}"
          userConfig="{inputUserConfig}"
          alertPlayers="{[]}"
        />
      </div>
    </div>

    <div class="mb-2 form-style">
      <h6 class="text-center">その他</h6>

      <!-- save-screenshot -->
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

      <!-- save-temp-arena-info -->
      <div class="form-check form-switch">
        <input
          class="form-check-input"
          type="checkbox"
          id="save-temp-arena-info"
          bind:checked="{inputUserConfig.save_temp_arena_info}"
        />
        <label class="form-check-label" for="save-temp-arena-info"
          >【開発用】自動で戦闘情報(<i>tempArenaInfo.json</i>)を保存する</label
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

      <!-- send-report -->
      <div class="form-check form-switch">
        <input
          class="form-check-input"
          type="checkbox"
          id="send-report"
          bind:checked="{inputUserConfig.send_report}"
        />
        <label class="form-check-label" for="send-report"
          >アプリ改善のためのデータ送信を許可する</label
        >
        <ul>
          <li>エラーログ</li>
          <li>設定値(<i>config/user.json</i>)</li>
          <li>戦闘情報(<i>tempArenaInfo.json</i>)</li>
        </ul>
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
