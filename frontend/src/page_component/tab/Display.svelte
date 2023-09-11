<script lang="ts">
  import { DispName } from "src/lib/DispName";
  import { sampleTeam } from "src/lib/rating/RatingConst";
  import { displayColumns } from "src/lib/util";
  import ConfirmModal from "src/other_component/ConfirmModal.svelte";
  import StatisticsTable from "src/other_component/StatisticsTable.svelte";
  import { storedUserConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import { Button, FormGroup, Input, Label } from "sveltestrap";
  import { domain } from "wailsjs/go/models";

  export let inputUserConfig: domain.UserConfig;
  export let defaultUserConfig: domain.UserConfig;

  let resetModal: ConfirmModal;

  const dispatch = createEventDispatcher();
  const [commons, shipOnlys, overallOnlys] = displayColumns();

  const onResetConfirmed = () => {
    inputUserConfig.font_size = defaultUserConfig.font_size;
    inputUserConfig.displays = defaultUserConfig.displays;
    inputUserConfig.custom_color = defaultUserConfig.custom_color;
    inputUserConfig.custom_digit = defaultUserConfig.custom_digit;

    dispatch("Change");
  };
</script>

<ConfirmModal
  bind:this={resetModal}
  message="表示設定をリセットしますか？"
  on:onConfirmed={onResetConfirmed}
/>

<FormGroup>
  <Label>文字サイズ</Label>
  <Input
    type="select"
    bind:value={inputUserConfig.font_size}
    on:change={() => dispatch("Change")}
  >
    {#each DispName.FONT_SIZES as fs}
      <option selected={fs.key === $storedUserConfig.font_size} value={fs.key}
        >{fs.value}</option
      >
    {/each}
  </Input>
</FormGroup>

<FormGroup>
  <Label>表示項目</Label>
  <table class="table table-sm table-text-color w-auto td-multiple">
    <thead>
      <tr>
        {#each ["項目", "艦成績", "総合成績", "小数点以下の桁数"] as columns}
          <th>{columns}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each commons as column}
        <tr>
          <td>
            {column.fullDisplayName()}
          </td>

          <td>
            <Input
              type="switch"
              bind:checked={inputUserConfig.displays.ship[column.displayKey()]}
              on:change={() => dispatch("Change")}
            />
          </td>

          <td>
            <Input
              type="switch"
              bind:checked={inputUserConfig.displays.overall[
                column.displayKey()
              ]}
              on:change={() => dispatch("Change")}
            />
          </td>

          <td>
            <FormGroup>
              <Input
                type="select"
                bind:value={inputUserConfig.custom_digit[column.displayKey()]}
                on:change={() => dispatch("Change")}
              >
                {#each [0, 1, 2] as digit}
                  <option
                    selected={digit ===
                      inputUserConfig.custom_digit[column.displayKey()]}
                    value={digit}>{digit}</option
                  >
                {/each}
              </Input>
            </FormGroup>
          </td>
        </tr>
      {/each}

      {#each shipOnlys as column}
        <tr>
          <td>
            {column.fullDisplayName()}
          </td>

          <td>
            <Input
              type="switch"
              bind:checked={inputUserConfig.displays.ship[column.displayKey()]}
              on:change={() => dispatch("Change")}
            />
          </td>

          <td></td>

          <td>
            <FormGroup>
              <Input
                type="select"
                bind:value={inputUserConfig.custom_digit[column.displayKey()]}
                on:change={() => dispatch("Change")}
              >
                {#each [0, 1, 2] as digit}
                  <option
                    selected={digit ===
                      inputUserConfig.custom_digit[column.displayKey()]}
                    value={digit}>{digit}</option
                  >
                {/each}
              </Input>
            </FormGroup>
          </td>
        </tr>
      {/each}

      {#each overallOnlys as column}
        <tr>
          <td>
            {column.fullDisplayName()}
          </td>

          <td></td>

          <td>
            <Input
              type="switch"
              bind:checked={inputUserConfig.displays.overall[
                column.displayKey()
              ]}
              on:change={() => dispatch("Change")}
            />
          </td>

          <td>
            <FormGroup>
              <Input
                type="select"
                bind:value={inputUserConfig.custom_digit[column.displayKey()]}
                on:change={() => dispatch("Change")}
              >
                {#each [0, 1, 2] as digit}
                  <option
                    selected={digit ===
                      inputUserConfig.custom_digit[column.displayKey()]}
                    value={digit}>{digit}</option
                  >
                {/each}
              </Input>
            </FormGroup>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</FormGroup>

<FormGroup>
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
      {#each DispName.SKILL_LEVELS as sl}
        <tr>
          <td>{sl.value}</td>
          <td>
            <Input
              type="color"
              bind:value={inputUserConfig.custom_color.skill.text[sl.key]}
              on:input={() => dispatch("Change")}
            />
          </td>

          <td>
            <Input
              type="color"
              bind:value={inputUserConfig.custom_color.skill.background[sl.key]}
              on:input={() => dispatch("Change")}
            />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</FormGroup>

<FormGroup>
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
      {#each DispName.TIER_GROUPS as tg}
        <tr>
          <td>{tg.value}</td>
          <td>
            <Input
              type="color"
              bind:value={inputUserConfig.custom_color.tier.own[tg.key]}
              on:input={() => dispatch("Change")}
            />
          </td>

          <td>
            <Input
              type="color"
              bind:value={inputUserConfig.custom_color.tier.other[tg.key]}
              on:input={() => dispatch("Change")}
            />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</FormGroup>

<FormGroup>
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
      {#each DispName.SHIP_TYPES as st}
        <tr>
          <td>{st.value}</td>
          <td>
            <Input
              type="color"
              bind:value={inputUserConfig.custom_color.ship_type.own[st.key]}
              on:input={() => dispatch("Change")}
            />
          </td>

          <td>
            <Input
              type="color"
              bind:value={inputUserConfig.custom_color.ship_type.other[st.key]}
              on:input={() => dispatch("Change")}
            />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</FormGroup>

<FormGroup>
  <Label>プレイヤー名の背景色</Label>
  <Input
    type="select"
    bind:value={inputUserConfig.custom_color.player_name}
    on:change={() => dispatch("Change")}
  >
    {#each DispName.PLAYER_NAME_COLORS as pnc}
      <option
        selected={pnc.key === $storedUserConfig.custom_color.player_name}
        value={pnc.key}>{pnc.value}</option
      >
    {/each}
  </Input>
</FormGroup>

<div>
  <p>プレビュー</p>

  <div style="font-size: {inputUserConfig.font_size};">
    <StatisticsTable teams={[sampleTeam()]} userConfig={inputUserConfig} />
  </div>

  <Button
    size="sm"
    color="warning"
    label="リセット"
    on:click={resetModal.toggle}>リセット</Button
  >
</div>
