<script lang="ts">
  import type { vo } from "wailsjs/go/models.js";
  import {
    ApplyUserConfig,
    UserConfig,
    SelectDirectory,
    Cwd,
  } from "../wailsjs/go/main/App.js";
  import { createEventDispatcher } from "svelte";
  import Const from "./Const.js";
  import Textfield from "@smui/textfield";
  import IconButton from "@smui/icon-button";
  import Select, { Option } from "@smui/select";
  import Checkbox from "@smui/checkbox";
  import Switch from "@smui/switch";
  import Button, { Label } from "@smui/button";
  import FormField from "@smui/form-field";

  const dispatch = createEventDispatcher();

  let inputConfig: vo.UserConfig = Const.DEFAULT_USER_CONFIG;

  let cwd: string;

  const fontSizes = ["x-small", "small", "medium", "large", "x-large"];

  // TODO 翻訳
  const displayLabels: { [key: string]: string } = {
    pr: "PR",
    damage: "Dmg",
    win_rate: "勝率",
    kd_rate: "K/D",
    win_survived_rate: "勝利生存率",
    lose_survived_rate: "敗北生存率",
    exp: "Exp",
    battles: "戦闘数",
    avg_tier: "平均T",
    using_ship_type_rate: "艦種割合",
    using_tier_rate: "ティア割合",
  };

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

  Cwd()
    .then((result) => (cwd = result))
    .catch((error) => "");
</script>

<div class="form-style">
  <div>
    <Textfield
      type="text"
      label="World of Warshipsインストールフォルダ"
      bind:value={inputConfig.install_path}
    />
    <IconButton class="material-icons" on:click={selectDirectory}
      >folder_open</IconButton
    >
  </div>

  <div>
    <Textfield type="text" label="AppID" bind:value={inputConfig.appid} />
  </div>

  <div>
    <Select bind:value={inputConfig.font_size} label="文字サイズ">
      {#each fontSizes as fontSize}
        <Option value={fontSize}>{fontSize}</Option>
      {/each}
    </Select>
  </div>

  <div style="display: flex; align-items: right;">
    <div>
      {#each Object.entries(inputConfig.displays.ship) as [k, v]}
        <FormField>
          <Checkbox bind:checked={inputConfig.displays.ship[k]} />
          <span slot="label">艦:{displayLabels[k]}</span>
        </FormField>
        <br />
      {/each}
    </div>
    <div>
      {#each Object.entries(inputConfig.displays.overall) as [k, v]}
        <FormField>
          <Checkbox bind:checked={inputConfig.displays.overall[k]} />
          <span slot="label">総合:{displayLabels[k]}</span>
        </FormField>
        <br />
      {/each}
    </div>
  </div>

  <div>
    <FormField>
      <Switch bind:checked={inputConfig.save_screenshot} />
      <span slot="label"
        >自動でスクリーンショットを保存する<br />(<i>{cwd}/screenshot</i>)</span
      >
    </FormField>
  </div>

  <div>
    <FormField>
      <Switch bind:checked={inputConfig.save_temp_arena_info} />
      <span slot="label"
        >【開発用】自動で戦闘情報(<i>tempArenaInfo.json</i>)を保存する<br />(<i
          >{cwd}/temp_arena_info</i
        >)</span
      >
    </FormField>
  </div>
</div>

<div class="centerize">
  <Button on:click={clickApply} variant="raised">
    <Label>適用</Label>
  </Button>
</div>
