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
import type { vo } from "../../wailsjs/go/models";
import StatisticsTable from "../other_component/StatisticsTable.svelte";
import { Badge, Button, FormGroup, Input, Label } from "sveltestrap";

const dispatch = createEventDispatcher();

const displayColumns = ["項目", "艦成績", "総合成績", "小数点以下の桁数"];
const skillColorColumns = ["スキル", "文字色", "背景色"];
const tierColorColumns = ["Tier", "使用艦", "非使用艦"];
const shipTypeColorColumns = ["艦種", "使用艦", "非使用艦"];

let inputUserConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => {
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
    const path = await SelectDirectory();
    if (!path) return;
    inputUserConfig.install_path = path;
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
  const shipKeys = Object.keys($storedUserConfig.displays.ship);
  const overallKeys = Object.keys($storedUserConfig.displays.overall);
  displayKeys = Array.from(new Set([...shipKeys, ...overallKeys]));
}

main();
</script>

<div class="m-3 center">
  <!-- install path -->
  <FormGroup class="center">
    <Label
      >World of Warshipsインストールフォルダ <Badge color="danger">必須</Badge>
    </Label>
    <div class="d-flex justify-content-center">
      <Input
        type="text"
        class="text-form"
        style="font-size: {$storedUserConfig.font_size};"
        bind:value="{inputUserConfig.install_path}"
      />
      <Button
        color="secondary"
        style="font-size: {$storedUserConfig.font_size};"
        on:click="{selectDirectory}">選択</Button
      >
    </div>
  </FormGroup>

  <!-- appid -->
  <FormGroup class="center">
    <Label>アプリケーションID <Badge color="danger">必須</Badge></Label>
    <Input
      type="text"
      class="text-form"
      style="font-size: {$storedUserConfig.font_size};"
      bind:value="{inputUserConfig.appid}"
    />
    <div>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a
        class="td-link"
        href="#"
        on:click="{() => BrowserOpenURL('https://developers.wargaming.net/')}"
        >Developer Room</a
      >で作成したIDを入力してください。
    </div>
  </FormGroup>

  <!-- font-size -->
  <FormGroup class="center">
    <Label>文字サイズ</Label>
    <Input
      type="select"
      class="w-auto"
      style="font-size: {$storedUserConfig.font_size};"
      bind:value="{inputUserConfig.font_size}"
    >
      {#await FontSizes() then fontSizes}
        {#each fontSizes as fs}
          <option selected="{fs === $storedUserConfig.font_size}" value="{fs}"
            >{Const.FONT_SIZE[fs]}</option
          >
        {/each}
      {/await}
    </Input>
  </FormGroup>

  <!-- display values -->
  <FormGroup class="center">
    <Label>表示項目</Label>
    <Input
      type="switch"
      style="font-size: {$storedUserConfig.font_size};"
      on:change="{toggleAll}"
      checked="{Object.values(inputUserConfig.displays.ship).filter((it) => !it)
        .length === 0 &&
        Object.values(inputUserConfig.displays.overall).filter((it) => !it)
          .length === 0}"
      label="全選択"
    />

    <table class="table table-sm table-text-color w-auto td-multiple">
      <thead>
        <tr>
          {#each displayColumns as columns}
            <th>{columns}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each displayKeys as key}
          <tr>
            <td>
              {Const.COLUMN_NAMES[key].full}
            </td>

            <td>
              {#if inputUserConfig.displays.ship[key] !== undefined}
                <Input
                  type="switch"
                  class="center"
                  bind:checked="{inputUserConfig.displays.ship[key]}"
                />
              {/if}
            </td>

            <td>
              {#if inputUserConfig.displays.overall[key] !== undefined}
                <Input
                  type="switch"
                  class="center"
                  bind:checked="{inputUserConfig.displays.overall[key]}"
                />
              {/if}
            </td>

            <td>
              {#if inputUserConfig.custom_digit[key] !== undefined}
                <Input
                  type="select"
                  class="p-1 m-1"
                  style="font-size: {$storedUserConfig.font_size};"
                  bind:value="{inputUserConfig.custom_digit[key]}"
                >
                  {#each [0, 1, 2] as digit}
                    <option
                      selected="{digit === inputUserConfig.custom_digit[key]}"
                      value="{digit}">{digit}</option
                    >
                  {/each}
                </Input>
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </FormGroup>

  <!-- skill color -->
  <FormGroup class="center">
    <Label>スキル別カラー</Label>
    <table class="table table-sm table-text-color w-auto td-multiple">
      <thead>
        <tr>
          {#each skillColorColumns as column}
            <th>{column}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each Object.keys(inputUserConfig.custom_color.skill.text) as key}
          <tr>
            <td>
              {Const.SKILL_LEVEL_LABELS[key]}
            </td>
            <td>
              <div class="d-flex justify-content-center">
                <Input
                  type="color"
                  class="m-1"
                  bind:value="{inputUserConfig.custom_color.skill.text[key]}"
                />
                <Button
                  size="sm"
                  color="success"
                  class="m-1"
                  style="font-size: {$storedUserConfig.font_size};"
                  on:click="{() => {
                    inputUserConfig.custom_color.skill.text[key] =
                      defaultUserConfig.custom_color.skill.text[key];
                  }}">デフォルト色に戻す</Button
                >
              </div>
            </td>

            <td>
              <div class="d-flex justify-content-center">
                <Input
                  type="color"
                  class="m-1"
                  bind:value="{inputUserConfig.custom_color.skill.background[
                    key
                  ]}"
                />
                <Button
                  size="sm"
                  color="success"
                  class="m-1"
                  style="font-size: {$storedUserConfig.font_size};"
                  on:click="{() => {
                    inputUserConfig.custom_color.skill.background[key] =
                      defaultUserConfig.custom_color.skill.background[key];
                  }}">デフォルト色に戻す</Button
                >
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </FormGroup>

  <!-- tier group color -->
  <FormGroup class="center">
    <Label>ティア別カラー</Label>
    <table class="table table-sm table-text-color w-auto td-multiple">
      <thead>
        <tr>
          {#each tierColorColumns as column}
            <th>{column}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each Object.keys(inputUserConfig.custom_color.tier.own) as key}
          <tr>
            <td>
              {Const.TIER_GROUP_LABELS[key]}
            </td>
            <td>
              <div class="d-flex justify-content-center">
                <Input
                  type="color"
                  class="m-1"
                  bind:value="{inputUserConfig.custom_color.tier.own[key]}"
                />
                <Button
                  size="sm"
                  color="success"
                  class="m-1"
                  style="font-size: {$storedUserConfig.font_size};"
                  on:click="{() => {
                    inputUserConfig.custom_color.tier.own[key] =
                      defaultUserConfig.custom_color.tier.own[key];
                  }}">デフォルト色に戻す</Button
                >
              </div>
            </td>

            <td>
              <div class="d-flex justify-content-center">
                <Input
                  type="color"
                  class="m-1"
                  bind:value="{inputUserConfig.custom_color.tier.other[key]}"
                />
                <Button
                  size="sm"
                  color="success"
                  class="m-1"
                  style="font-size: {$storedUserConfig.font_size};"
                  on:click="{() => {
                    inputUserConfig.custom_color.tier.other[key] =
                      defaultUserConfig.custom_color.tier.other[key];
                  }}">デフォルト色に戻す</Button
                >
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </FormGroup>

  <!-- ship type color -->
  <FormGroup class="center">
    <Label>艦種別カラー</Label>
    <table class="table table-sm table-text-color w-auto td-multiple">
      <thead>
        <tr>
          {#each shipTypeColorColumns as column}
            <th>{column}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each Object.keys(inputUserConfig.custom_color.ship_type.own) as key}
          <tr>
            <td>
              {Const.SHIP_TYPE_LABELS[key]}
            </td>
            <td>
              <div class="d-flex justify-content-center">
                <Input
                  type="color"
                  class="m-1"
                  bind:value="{inputUserConfig.custom_color.ship_type.own[key]}"
                />
                <Button
                  size="sm"
                  color="success"
                  class="m-1"
                  style="font-size: {$storedUserConfig.font_size};"
                  on:click="{() => {
                    inputUserConfig.custom_color.ship_type.own[key] =
                      defaultUserConfig.custom_color.ship_type.own[key];
                  }}">デフォルト色に戻す</Button
                >
              </div>
            </td>

            <td>
              <div class="d-flex justify-content-center">
                <Input
                  type="color"
                  class="m-1"
                  bind:value="{inputUserConfig.custom_color.ship_type.other[
                    key
                  ]}"
                />
                <Button
                  size="sm"
                  color="success"
                  class="m-1"
                  style="font-size: {$storedUserConfig.font_size};"
                  on:click="{() => {
                    inputUserConfig.custom_color.ship_type.other[key] =
                      defaultUserConfig.custom_color.ship_type.other[key];
                  }}">デフォルト色に戻す</Button
                >
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </FormGroup>

  <div class="center">
    <p>プレビュー</p>

    <div style="font-size: {inputUserConfig.font_size};">
      <StatisticsTable
        teams="{sampleTeams}"
        userConfig="{inputUserConfig}"
        alertPlayers="{[]}"
      />
    </div>
  </div>

  <!-- team average -->
  {#if defaultUserConfig}
    <FormGroup class="center">
      <Label>チーム平均に含める最小戦闘数(艦戦績)</Label>
      <Input
        type="range"
        min="{defaultUserConfig.team_average.min_ship_battles}"
        max="{100}"
        step="{1}"
        class="text-form"
        style="font-size: {$storedUserConfig.font_size};"
        bind:value="{inputUserConfig.team_average.min_ship_battles}"
      />
      <div>
        {inputUserConfig.team_average.min_ship_battles}
      </div>
    </FormGroup>

    <FormGroup class="center">
      <Label>チーム平均に含める最小戦闘数(総合成績)</Label>
      <Input
        type="range"
        min="{defaultUserConfig.team_average.min_overall_battles}"
        max="{10000}"
        step="{100}"
        class="text-form"
        style="font-size: {$storedUserConfig.font_size};"
        bind:value="{inputUserConfig.team_average.min_overall_battles}"
      />
      <div>
        {inputUserConfig.team_average.min_overall_battles}
      </div>
    </FormGroup>
  {/if}

  <FormGroup class="center">
    <Label>その他</Label>
    <ul>
      <li class="my-1">
        <Input
          type="switch"
          label="自動でスクリーンショットを保存する"
          bind:checked="{inputUserConfig.save_screenshot}"
        />
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
          class="td-link"
          href="#"
          on:click="{() => openDirectory('screenshot/')}"
          ><i class="bi bi-folder2-open"></i> 保存フォルダを開く
        </a>
      </li>
      <li class="my-1">
        <Input
          type="switch"
          label="【開発用】自動で戦闘情報(tempArenaInfo.json)を保存する"
          bind:checked="{inputUserConfig.save_temp_arena_info}"
        />
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a
          class="td-link"
          href="#"
          on:click="{() => openDirectory('temp_arena_info/')}"
          ><i class="bi bi-folder2-open"></i> 保存フォルダを開く
        </a>
      </li>
      <li class="my-1">
        <Input
          type="switch"
          label="アプリ改善のためのデータ送信を許可する"
          bind:checked="{inputUserConfig.send_report}"
        />
        <ul>
          <li>エラーログ</li>
          <li>設定値(config/user.json)</li>
          <li>戦闘情報(tempArenaInfo.json)</li>
        </ul>
      </li>
    </ul>
  </FormGroup>

  <!-- apply -->
  <button
    type="button"
    class="btn btn-primary mb-3"
    style="font-size: {$storedUserConfig.font_size};"
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
