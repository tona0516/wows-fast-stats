<script lang="ts">
  import clone from "clone";
  import { createEventDispatcher } from "svelte";
  import { get } from "svelte/store";
  import {
    ApplyRequiredUserConfig,
    ApplyUserConfig,
    DefaultUserConfig,
    FontSizes,
    OpenDirectory,
    SampleTeams,
    SelectDirectory,
    UserConfig,
  } from "../../wailsjs/go/main/App";
  import { BrowserOpenURL, LogDebug } from "../../wailsjs/runtime/runtime";
  import { storedUserConfig } from "../stores";
  import { Const } from "../Const";
  import type { domain, vo } from "../../wailsjs/go/models";
  import StatisticsTable from "../other_component/StatisticsTable.svelte";
  import {
    Alert,
    Badge,
    Button,
    FormGroup,
    Input,
    Label,
    TabContent,
    TabPane,
  } from "sveltestrap";
  import ConfirmModal from "../other_component/ConfirmModal.svelte";

  const dispatch = createEventDispatcher();

  let inputUserConfig = get(storedUserConfig);
  storedUserConfig.subscribe((it) => {
    inputUserConfig = clone(it);
  });
  let defaultUserConfig: domain.UserConfig;
  let isLoading = false;
  let sampleTeams: domain.Team[] = [];
  let displayKeys: string[] = [];
  let resetDisplaySettingConfirmModal: ConfirmModal;
  let validatedResult: vo.ValidatedResult;

  async function clickApply() {
    isLoading = true;
    try {
      validatedResult = await ApplyRequiredUserConfig(
        inputUserConfig.install_path,
        inputUserConfig.appid
      );

      const errorTexts = Object.values(validatedResult);
      const isValid =
        errorTexts.filter((it) => it == "").length === errorTexts.length;
      if (isValid) {
        const latest = await UserConfig();
        storedUserConfig.set(latest);
        dispatch("UpdateSuccess", { message: "設定を更新しました。" });
      }
    } catch (error) {
      inputUserConfig = get(storedUserConfig);
      dispatch("Failure", { message: error });
    } finally {
      isLoading = false;
    }
  }

  async function silentApply() {
    LogDebug("silentApply");
    // Note: for the following sveltestrap bug
    // https://github.com/bestguy/sveltestrap/issues/461
    await new Promise((resolve) => setTimeout(resolve, 100));

    try {
      await ApplyUserConfig(inputUserConfig);
      const latest = await UserConfig();
      storedUserConfig.set(latest);
    } catch (error) {
      inputUserConfig = get(storedUserConfig);
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

  async function toggleAll(e: any) {
    const isSelectAll: boolean = e.target.checked;

    Object.keys(inputUserConfig.displays.ship).forEach(
      (key) => (inputUserConfig.displays.ship[key] = isSelectAll)
    );
    Object.keys(inputUserConfig.displays.overall).forEach(
      (key) => (inputUserConfig.displays.overall[key] = isSelectAll)
    );

    await silentApply();
  }

  async function clickResetDisplaySetting() {
    inputUserConfig.font_size = defaultUserConfig.font_size;
    inputUserConfig.displays = defaultUserConfig.displays;
    inputUserConfig.custom_color = defaultUserConfig.custom_color;
    inputUserConfig.custom_digit = defaultUserConfig.custom_digit;

    await silentApply();
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

<ConfirmModal
  message="表示設定をリセットしますか？"
  on:onConfirmed={clickResetDisplaySetting}
  bind:this={resetDisplaySettingConfirmModal}
/>

<div class="m-3">
  <TabContent>
    <TabPane class="py-3 center" tabId="required" tab="必須設定" active>
      <!-- install path -->
      <FormGroup class="center">
        <Label
          >World of Warshipsインストールフォルダ <Badge color="danger"
            >必須</Badge
          >
        </Label>
        <div class="d-flex justify-content-center">
          <Input
            type="text"
            class="text-form w-auto"
            style="font-size: {$storedUserConfig.font_size};"
            bind:value={inputUserConfig.install_path}
          />
          <Button
            color="secondary"
            style="font-size: {$storedUserConfig.font_size};"
            on:click={selectDirectory}>選択</Button
          >
        </div>

        {#if validatedResult?.install_path}
          <Alert color="danger" class="m-1">
            {validatedResult?.install_path}
          </Alert>
        {/if}
      </FormGroup>

      <!-- appid -->
      <FormGroup class="center">
        <Label>アプリケーションID <Badge color="danger">必須</Badge></Label>
        <Input
          type="text"
          class="text-form w-auto"
          style="font-size: {$storedUserConfig.font_size};"
          bind:value={inputUserConfig.appid}
        />
        <div>
          <!-- svelte-ignore a11y-invalid-attribute -->
          <a
            class="td-link"
            href="#"
            on:click={() => BrowserOpenURL("https://developers.wargaming.net/")}
            >Developer Room</a
          >で作成したIDを入力してください。
        </div>

        {#if validatedResult?.appid}
          <Alert color="danger" class="m-1">
            {validatedResult?.appid}
          </Alert>
        {/if}
      </FormGroup>

      <!-- apply -->
      <FormGroup class="center">
        <Button
          color="primary"
          style="font-size: {$storedUserConfig.font_size};"
          disabled={isLoading}
          on:click={clickApply}
        >
          {#if isLoading}
            <span
              class="spinner-border spinner-border-sm"
              role="status"
              aria-hidden="true"
            />
            更新中...
          {:else}
            保存
          {/if}
        </Button>
      </FormGroup>
    </TabPane>

    <TabPane class="py-3 center" tabId="display" tab="表示設定">
      <!-- font-size -->
      <FormGroup class="center">
        <Label>文字サイズ</Label>
        <Input
          type="select"
          class="w-auto"
          style="font-size: {$storedUserConfig.font_size};"
          bind:value={inputUserConfig.font_size}
          on:change={silentApply}
        >
          {#await FontSizes() then fontSizes}
            {#each fontSizes as fs}
              <option selected={fs === $storedUserConfig.font_size} value={fs}
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
          on:change={toggleAll}
          checked={Object.values(inputUserConfig.displays.ship).filter(
            (it) => !it
          ).length === 0 &&
            Object.values(inputUserConfig.displays.overall).filter((it) => !it)
              .length === 0}
          label="全選択"
        />

        <table class="table table-sm table-text-color w-auto td-multiple">
          <thead>
            <tr>
              {#each ["項目", "艦成績", "総合成績", "小数点以下の桁数"] as columns}
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
                      bind:checked={inputUserConfig.displays.ship[key]}
                      on:change={silentApply}
                    />
                  {/if}
                </td>

                <td>
                  {#if inputUserConfig.displays.overall[key] !== undefined}
                    <Input
                      type="switch"
                      class="center"
                      bind:checked={inputUserConfig.displays.overall[key]}
                      on:change={silentApply}
                    />
                  {/if}
                </td>

                <td>
                  {#if inputUserConfig.custom_digit[key] !== undefined}
                    <Input
                      type="select"
                      class="p-1 m-1"
                      style="font-size: {$storedUserConfig.font_size};"
                      bind:value={inputUserConfig.custom_digit[key]}
                      on:change={silentApply}
                    >
                      {#each [0, 1, 2] as digit}
                        <option
                          selected={digit === inputUserConfig.custom_digit[key]}
                          value={digit}>{digit}</option
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
              {#each ["スキル", "文字色", "背景色"] as column}
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
                  <Input
                    type="color"
                    class="m-1"
                    bind:value={inputUserConfig.custom_color.skill.text[key]}
                    on:input={silentApply}
                  />
                </td>

                <td>
                  <Input
                    type="color"
                    class="m-1"
                    bind:value={inputUserConfig.custom_color.skill.background[
                      key
                    ]}
                    on:input={silentApply}
                  />
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
              {#each ["Tier", "使用艦", "非使用艦"] as column}
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
                  <Input
                    type="color"
                    class="m-1"
                    bind:value={inputUserConfig.custom_color.tier.own[key]}
                    on:input={silentApply}
                  />
                </td>

                <td>
                  <Input
                    type="color"
                    class="m-1"
                    bind:value={inputUserConfig.custom_color.tier.other[key]}
                    on:input={silentApply}
                  />
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
              {#each ["艦種", "使用艦", "非使用艦"] as column}
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
                  <Input
                    type="color"
                    class="m-1"
                    bind:value={inputUserConfig.custom_color.ship_type.own[key]}
                    on:input={silentApply}
                  />
                </td>

                <td>
                  <Input
                    type="color"
                    class="m-1"
                    bind:value={inputUserConfig.custom_color.ship_type.other[
                      key
                    ]}
                    on:input={silentApply}
                  />
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
            teams={sampleTeams}
            userConfig={inputUserConfig}
            alertPlayers={[]}
          />
        </div>

        <Button
          size="sm"
          color="warning"
          style="font-size: {$storedUserConfig.font_size};"
          on:click={resetDisplaySettingConfirmModal.toggle}>リセット</Button
        >
      </div>
    </TabPane>
    <TabPane class="py-3 center" tabId="team-summary" tab="チームサマリー設定">
      <!-- team average -->
      {#if defaultUserConfig?.team_average}
        {@const teamAvg = defaultUserConfig.team_average}
        <FormGroup class="center">
          <Label>チーム平均に含める最小戦闘数(艦戦績)</Label>
          <Input
            type="range"
            min={teamAvg.min_ship_battles}
            max={100}
            step={1}
            class="text-form w-auto"
            style="font-size: {$storedUserConfig.font_size};"
            bind:value={inputUserConfig.team_average.min_ship_battles}
            on:change={silentApply}
          />
          <div>
            {inputUserConfig.team_average.min_ship_battles}
          </div>
        </FormGroup>

        <FormGroup class="center">
          <Label>チーム平均に含める最小戦闘数(総合成績)</Label>
          <Input
            type="range"
            min={teamAvg.min_overall_battles}
            max={3000}
            step={10}
            class="text-form w-auto"
            style="font-size: {$storedUserConfig.font_size};"
            bind:value={inputUserConfig.team_average.min_overall_battles}
            on:change={silentApply}
          />
          <div>
            {inputUserConfig.team_average.min_overall_battles}
          </div>
        </FormGroup>
      {/if}
    </TabPane>
    <TabPane class="py-3 center" tabId="other" tab="その他設定">
      <FormGroup class="center">
        <ul>
          <li class="mb-3">
            <Input
              type="switch"
              label="自動でスクリーンショットを保存する"
              bind:checked={inputUserConfig.save_screenshot}
              on:change={silentApply}
            />
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
              class="td-link"
              href="#"
              on:click={() => openDirectory("screenshot/")}
              ><i class="bi bi-folder2-open" /> 保存フォルダを開く
            </a>
          </li>
          <li class="mb-3">
            <Input
              type="switch"
              label="【開発用】自動で戦闘情報(tempArenaInfo.json)を保存する"
              bind:checked={inputUserConfig.save_temp_arena_info}
              on:change={silentApply}
            />
            <!-- svelte-ignore a11y-invalid-attribute -->
            <a
              class="td-link"
              href="#"
              on:click={() => openDirectory("temp_arena_info/")}
              ><i class="bi bi-folder2-open" /> 保存フォルダを開く
            </a>
          </li>
          <li class="mb-3">
            <Input
              type="switch"
              label="アプリ改善のためのデータ送信を許可する"
              bind:checked={inputUserConfig.send_report}
              on:change={silentApply}
            />
            <ul>
              <li>アプリバージョン</li>
              <li>エラーログ</li>
              <li>設定値(config/user.json)</li>
              <li>戦闘情報(tempArenaInfo.json)</li>
            </ul>
          </li>
          <li class="mb-3">
            <Input
              type="switch"
              label="新しいバージョンがある場合に通知する"
              bind:checked={inputUserConfig.notify_updatable}
              on:change={silentApply}
            />
          </li>
        </ul>
      </FormGroup>
    </TabPane>
  </TabContent>
</div>
